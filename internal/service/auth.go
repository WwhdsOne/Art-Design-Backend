package service

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/authutils"
	"Art-Design-Backend/pkg/constant"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/redisx"
	"Art-Design-Backend/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type AuthService struct {
	UserRepo      *repository.UserRepository         // 用户Repo
	RoleRepo      *repository.RoleRepository         // 角色Repo
	UserRolesRepo *repository.UserRolesRepository    // 用户角色关联Repo
	GormTX        *repository.GormTransactionManager // gorm事务管理
	Redis         *redisx.RedisWrapper               // redis
	Jwt           *jwt.JWT                           // jwt相关
}

// Login 登录
func (s *AuthService) Login(c *gin.Context, u *request.Login) (tokenStr string, err error) {
	// 只验证启用状态的用户
	user, err := s.UserRepo.GetLoginUserByUsername(c, u.Username)
	if err != nil {
		return
	}
	if user.Status != 1 {
		err = fmt.Errorf("用户未启用")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		// todo后续修改为统一鉴权错误
		err = fmt.Errorf("用户名或密码错误")
		return
	}
	claims := jwt.NewBaseClaims(user.ID)
	return s.createToken(claims)
}

func (s *AuthService) RefreshToken(c *gin.Context) (tokenStr string, err error) {
	// 路径参数中获取用户 id
	idStr := c.Param("id")
	// 确保该用户在登录状态
	sessionKey := rediskey.SESSION + idStr
	_, err = s.Redis.Get(sessionKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 用户未登录
			err = errors.New("用户未在登录状态")
			zap.L().Warn("用户未登录", zap.String("user_id", idStr))
			return
		}
		// Redis 出错
		zap.L().Error("Redis 获取 Session 失败", zap.String("key", sessionKey), zap.Error(err))
		return
	}
	// 根据原有的用户 Claims 创建一个新的 token
	claims := authutils.GetClaims(c)
	return s.createToken(claims.BaseClaims)
}

// createToken 创建一个 token 并处理会话
func (s *AuthService) createToken(baseClaims jwt.BaseClaims) (tokenStr string, err error) {
	// 调用 jwt 服务的 CreateToken 方法生成令牌
	tokenStr, err = s.Jwt.CreateToken(baseClaims)
	if err != nil {
		return
	}

	// 将用户 ID 转换为字符串形式
	id := strconv.FormatInt(baseClaims.UserID, 10)

	// 获取 Session
	sessionKey := rediskey.SESSION + id

	// 检查旧会话是否存在
	existToken, err := s.Redis.Get(sessionKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存未命中，正常情况，Debug 级别日志
			// 后续需要创建新的会话
			zap.L().Debug("Redis 获取旧 Session 未命中", zap.String("key", sessionKey), zap.Error(err))
		} else {
			zap.L().Error("Redis 获取 Session 缓存失败", zap.String("key", sessionKey), zap.Error(err))
			return
		}
	}

	// 尝试删除
	delKeys := []string{
		rediskey.LOGIN + existToken,
		sessionKey,
	}
	if err = s.Redis.PipelineDel(delKeys); err != nil {
		zap.L().Error("删除旧Session失败", zap.Strings("keys", delKeys), zap.Error(err))
		return
	}

	// 准备新会话和登录状态键值对
	keyVals := [][2]string{
		{rediskey.LOGIN + tokenStr, id},
		{rediskey.SESSION + id, tokenStr},
	}

	// 把键提取出来用于日志打印
	keys := []string{
		rediskey.LOGIN + tokenStr,
		rediskey.SESSION + id,
	}

	err = s.Redis.PipelineSet(keyVals, s.Jwt.ExpiresTime)
	if err != nil {
		zap.L().Error("设置新Session失败", zap.Strings("keys", keys), zap.Error(err))
		return
	}

	return
}

// LogoutToken 注销 token
func (s *AuthService) LogoutToken(c *gin.Context) (err error) {
	// 从请求头中获取 token
	tokenStr := authutils.GetToken(c)
	// 解析传入的令牌以获取用户信息
	claims, err := s.Jwt.ParseToken(tokenStr)
	if err != nil {
		return
	}

	// 将用户 ID 转换为字符串形式
	id := strconv.FormatInt(claims.BaseClaims.UserID, 10)

	// 准备需要删除的 Redis 键
	delKeys := []string{
		rediskey.SESSION + id,
		rediskey.LOGIN + tokenStr,
	}

	// 使用管道删除 Redis 中的会话和登录状态键
	err = s.Redis.PipelineDel(delKeys)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	return
}

// Register 注册
func (s *AuthService) Register(c *gin.Context, userReq *request.RegisterUser) (err error) {
	var user entity.User
	err = copier.Copy(&user, &userReq)
	// 处理密码（非指针字段）
	password, _ := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	user.Password = string(password)
	// 随机设置头像
	user.Avatar = constant.DefaultAvatar[utils.GenerateRandomNumber(0, len(constant.DefaultAvatar))]
	// 判重
	if err = s.UserRepo.CheckUserDuplicate(&user); err != nil {
		return err
	}
	// 启用事务插入用户表和用户角色关联表
	err = s.GormTX.Transaction(c, func(ctx context.Context) (err error) {
		if err = s.UserRepo.CreateUser(ctx, &user); err != nil {
			return
		}
		// id是新用户的主键ID
		// todo后续可能会换成动态获取
		if err = s.UserRolesRepo.AssignRoleToUser(ctx, user.ID, []int64{42838646763553030}); err != nil {
			return
		}
		return
	})
	if err != nil {
		zap.L().Error("注册失败", zap.Error(err))
		return
	}
	return
}

package service

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/constant"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/loginUtils"
	"Art-Design-Backend/pkg/redisx"
	"Art-Design-Backend/pkg/transaction"
	"Art-Design-Backend/pkg/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type AuthService struct {
	UserRepo    *repository.UserRepository // 用户Repo
	RoleRepo    *repository.RoleRepository // 角色Repo
	GormSession *transaction.GormSession   // gorm事务管理
	Redis       *redisx.RedisWrapper       // redis
	Jwt         *jwt.JWT                   // jwt相关
}

// Login 登录
func (s *AuthService) Login(c *gin.Context, u *request.Login) (tokenStr string, err error) {
	// 只验证启用状态的用户
	user, err := s.UserRepo.GetLoginUserByUsername(c, u.Username)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		// todo后续修改为统一鉴权错误
		err = fmt.Errorf("用户名或密码错误")
		return
	}
	roleIdList, err := s.RoleRepo.GetRoleIDListByUserID(c, user.ID)
	if err != nil {
		return
	}
	if len(roleIdList) == 0 {
		err = fmt.Errorf("用户未分配角色")
		return
	}
	claims := jwt.NewBaseClaims(user.ID, roleIdList)
	return s.createToken(claims)
}

func (s *AuthService) RefreshToken(c *gin.Context) (tokenStr string, err error) {
	// 路径参数中获取用户 id
	idStr := c.Param("id")
	// 确保该用户在登录状态
	val := s.Redis.Get(constant.SESSION + idStr)
	if val == "" {
		// 如果不存在，则返回错误
		// todo后续改为auth专用错误
		err = fmt.Errorf("用户未在登录状态")
		return
	}
	// 根据原有的用户 Claims 创建一个新的 token
	claims := loginUtils.GetClaims(c)
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

	// 检查 Redis 中是否已存在该用户的会话
	existToken := s.Redis.Get(constant.SESSION + id)
	if existToken != "" {
		// 如果存在，则删除现有的会话相关键
		delKeys := []string{constant.LOGIN + existToken, constant.SESSION + id}
		err = s.Redis.PipelineDelete(delKeys)
		if err != nil {
			return
		}
	}

	// 准备新会话和登录状态键值对
	keyVal := [][2]string{
		{constant.LOGIN + tokenStr, id},
		{constant.SESSION + id, tokenStr},
	}

	// 使用管道设置新的键值对到 Redis，并设置过期时间
	err = s.Redis.PipelineSet(keyVal, s.Jwt.ExpiresTime)
	return
}

// LogoutToken 注销 token
func (s *AuthService) LogoutToken(c *gin.Context) (err error) {
	// 从请求头中获取 token
	tokenStr := loginUtils.GetToken(c)
	// 解析传入的令牌以获取用户信息
	claims, err := s.Jwt.ParseToken(tokenStr)
	if err != nil {
		return
	}

	// 将用户 ID 转换为字符串形式
	id := strconv.FormatInt(claims.BaseClaims.UserID, 10)

	// 准备需要删除的 Redis 键
	delKeys := []string{constant.SESSION + id, constant.LOGIN + tokenStr}

	// 使用管道删除 Redis 中的会话和登录状态键
	err = s.Redis.PipelineDelete(delKeys)
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
	user.Avatar = constant.DefaultAvatar[utils.RandomNumber.Intn(len(constant.DefaultAvatar))]
	// 判重
	if err = s.UserRepo.CheckUserDuplicate(&user); err != nil {
		return err
	}
	// 启用事务插入用户表和用户角色关联表
	err = s.GormSession.Transaction(c, func(ctx context.Context) (err error) {
		if err = s.UserRepo.CreateUser(ctx, &user); err != nil {
			return
		}
		// id是新用户的主键ID
		// todo后续可能会换成动态获取
		if err = s.RoleRepo.AssignRoleToUser(ctx, user.ID, []int64{42838646763553030}); err != nil {
			return
		}
		return
	})
	return
}

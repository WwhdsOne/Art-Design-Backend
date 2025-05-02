package service

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/constant"
	"Art-Design-Backend/pkg/errorTypes"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/redisx"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type AuthService struct {
	UserRepo *repository.UserRepository // 用户Repo
	Redis    *redisx.RedisWrapper       // redis
	Jwt      *jwt.JWT                   // jwt相关
}

// Login 登录
func (s *AuthService) Login(ctx *gin.Context, u *request.Login) (tokenStr string, err error) {
	user, err := s.UserRepo.GetUserByUsername(ctx, u.Username)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		// 密码错误
		return "", errorTypes.NewGormError("用户名或密码错误")
	}
	claims := jwt.BaseClaims{
		ID: user.ID,
	}
	// 创建 token
	return s.CreateToken(claims)
}

// CreateToken 创建一个 token 并处理会话
func (s *AuthService) CreateToken(baseClaims jwt.BaseClaims) (tokenStr string, err error) {
	// 调用 jwt 服务的 CreateToken 方法生成令牌
	tokenStr, err = s.Jwt.CreateToken(baseClaims)
	if err != nil {
		return
	}

	// 将用户 ID 转换为字符串形式
	id := strconv.FormatInt(baseClaims.ID, 10)

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
func (s *AuthService) LogoutToken(tokenStr string) (err error) {
	// 解析传入的令牌以获取用户信息
	claims, err := s.Jwt.ParseToken(tokenStr)
	if err != nil {
		return
	}

	// 将用户 ID 转换为字符串形式
	id := strconv.FormatInt(claims.BaseClaims.ID, 10)

	// 准备需要删除的 Redis 键
	delKeys := []string{constant.SESSION + id, constant.LOGIN + tokenStr}

	// 使用管道删除 Redis 中的会话和登录状态键
	err = s.Redis.PipelineDelete(delKeys)
	return
}

// Register 注册
func (s *AuthService) Register(ctx *gin.Context, userReq *request.User) (err error) {
	var user entity.User
	// 属性复制，同时进行邮箱、手机号和密码加密操作
	err = copier.Copy(&user, &userReq)
	if err != nil {
		zap.L().Error("复制属性失败", zap.Error(err))
		return
	}
	// 判重
	if err = s.UserRepo.CheckUserDuplicate(&user); err != nil {
		return
	}
	// 插入
	if err = s.UserRepo.CreateUser(ctx, &user); err != nil {
		// 处理错误
		zap.L().Error("新增用户失败", zap.Error(err))
		return
	}
	return
}

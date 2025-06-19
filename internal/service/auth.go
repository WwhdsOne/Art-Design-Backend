package service

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/internal/repository/db"
	"Art-Design-Backend/pkg/authutils"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo          *repository.UserRepo       // 用户Repo
	AuthRepo          *repository.AuthRepo       // 登录缓存
	RoleRepo          *repository.RoleRepo       // 角色Repo
	GormTX            *db.GormTransactionManager // gorm事务管理
	Jwt               *jwt.JWT                   // jwt相关
	DefaultUserConfig *config.DefaultUserConfig  // 默认用户配置
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
		err = fmt.Errorf("用户名或密码错误")
		return
	}
	claims := jwt.NewBaseClaims(user.ID)
	return s.createToken(claims)
}

// createToken 创建一个 token 并处理会话
func (s *AuthService) createToken(baseClaims jwt.BaseClaims) (tokenStr string, err error) {
	// 1. 创建 JWT Token
	tokenStr, err = s.Jwt.CreateToken(baseClaims)
	if err != nil {
		zap.L().Error("创建 JWT Token 失败", zap.Error(err))
		return
	}
	userID := baseClaims.UserID

	// 2. 获取旧 token（Session）
	oldToken, err := s.AuthRepo.GetTokenByUserID(userID)
	if err != nil && !errors.Is(err, redis.Nil) {
		zap.L().Error("Redis 获取 Session 错误", zap.Error(err))
		return
	}

	// 3. 删除旧 token 映射
	if oldToken != "" {
		err = s.AuthRepo.DeleteOldSession(userID, oldToken)
		if err != nil {
			zap.L().Error("删除旧 Session 失败", zap.Int64("user_id", userID), zap.String("token", oldToken), zap.Error(err))
			return
		}
	}

	// 4. 设置新 token 映射
	err = s.AuthRepo.SetNewSession(userID, tokenStr, s.Jwt.ExpiresTime)
	if err != nil {
		zap.L().Error("设置新 Session 错误", zap.Int64("user_id", userID), zap.String("token", tokenStr), zap.Error(err))
		return
	}

	return
}

// LogoutToken 注销 token
func (s *AuthService) LogoutToken(c *gin.Context) (err error) {
	// 获取 token
	tokenStr := authutils.GetToken(c)

	// 解析 token 获取 claims
	claims, err := s.Jwt.ParseToken(tokenStr)
	if err != nil {
		return
	}
	userID := claims.BaseClaims.UserID

	// 调用 AuthCache 注销 token
	if err = s.AuthRepo.DeleteOldSession(userID, tokenStr); err != nil {
		zap.L().Error("注销 token 失败", zap.Int64("userID", userID), zap.String("token", tokenStr), zap.Error(err))
		return
	}

	return
}

// Register 注册
func (s *AuthService) Register(c *gin.Context, userReq *request.RegisterUser) (err error) {
	var user entity.User
	_ = copier.Copy(&user, &userReq)
	// 处理密码（非指针字段）
	password, _ := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	user.Password = string(password)
	// 随机设置头像
	user.Avatar = s.DefaultUserConfig.Avatars[utils.GenerateRandomNumber(0, len(s.DefaultUserConfig.Avatars))]
	// 判重
	if err = s.UserRepo.CheckUserDuplicate(context.TODO(), &user); err != nil {
		return
	}
	// 启用事务插入用户表和用户角色关联表
	err = s.GormTX.Transaction(c, func(ctx context.Context) (err error) {
		if err = s.UserRepo.CreateUser(ctx, &user); err != nil {
			return
		}
		// id是新用户的主键ID
		// todo后续可能会换成动态获取
		if err = s.RoleRepo.AddRolesToUser(ctx, user.ID, []int64{42838646763553030}); err != nil {
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

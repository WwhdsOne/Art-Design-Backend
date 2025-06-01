package service

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/aliyun"
	"Art-Design-Backend/pkg/authutils"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/redisx"
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"
	"strconv"
)

type UserService struct {
	UserRepo          *repository.UserRepository         // 用户Repo
	RoleRepo          *repository.RoleRepository         // 角色Repo
	UserRolesRepo     *repository.UserRolesRepository    // 用户角色Repo
	GormTX            *repository.GormTransactionManager // gorm事务管理
	OssClient         *aliyun.OssClient                  // 阿里云OSS
	Redis             *redisx.RedisWrapper               // redis
	DefaultUserConfig *config.DefaultUserConfig          // 默认用户配置
}

func (u *UserService) GetUserById(c *gin.Context) (res response.User, err error) {
	var user *entity.User
	id := authutils.GetUserID(c)
	if id == -1 {
		err = fmt.Errorf("当前用户未登录")
		return
	}
	// 获取用户信息
	if user, err = u.UserRepo.GetUserById(c, id); err != nil {
		return
	}
	// 获取用户角色列表
	roleList, err := u.getRoleListByUserID(c, user.ID)
	if len(roleList) == 0 {
		err = fmt.Errorf("当前用户未分配角色")
		return
	}
	user.Roles = roleList
	err = copier.Copy(&res, &user)
	if err != nil {
		zap.L().Error("复制属性失败", zap.Error(err))
		return
	}
	return
}

func (u *UserService) getRoleListByUserID(c context.Context, userID int64) (roleList []*entity.Role, err error) {
	userRoleKey := fmt.Sprintf(rediskey.UserRoleList+"%d", userID)

	// 1. 尝试从 Redis 获取缓存
	cacheStr, err := u.Redis.Get(userRoleKey)
	if err == nil {
		if err = sonic.Unmarshal([]byte(cacheStr), &roleList); err == nil {
			// 缓存命中，返回
			return
		}
		// 缓存解析失败，从数据库重新读取，防止数据污染
		zap.L().Error("用户角色对应关系缓存解析失败", zap.Error(err))
	} else if !errors.Is(err, redis.Nil) {
		zap.L().Error("获取用户角色对应关系缓存失败", zap.Error(err))
		// 若 Redis 出错（但不是未命中），直接返回错误
		return
	} else {
		zap.L().Debug("获取用户角色对应关系缓存未命中", zap.Int64("userID", userID))
	}

	// 2. 从数据库查询ID列表
	roleIDList, err := u.UserRolesRepo.GetRoleIDListByUserID(c, userID)
	if err != nil {
		zap.L().Error("获取用户角色列表失败", zap.Error(err))
		return nil, err
	}

	// 3. 根据ID列表去数据库或缓存查询数据
	for _, roleID := range roleIDList {
		var role *entity.Role
		key := strconv.FormatInt(roleID, 10)
		var roleJson string
		// 3.1 从 Redis 读取
		roleJson, err = u.Redis.Get(rediskey.RoleInfo + key)
		if err == nil {
			_ = sonic.UnmarshalString(roleJson, &role)
			return
		}
		// 3.2 从数据库读取
		role, err = u.RoleRepo.GetRoleByID(c, roleID)
		if err != nil {
			zap.L().Error("获取角色列表失败", zap.Error(err))
			return
		}
		// 3.3 角色信息异步写入 Redis
		go func() {
			roleJson, _ = sonic.MarshalString(role)
			err = u.Redis.Set(rediskey.RoleInfo+key, roleJson, rediskey.RoleInfoTTL)
			if err != nil {
				zap.L().Error("角色缓存写入失败", zap.Error(err))
			}
		}()
		roleList = append(roleList, role)
	}

	// 4. 用户角色对应关系写入 Redis 缓存
	cacheBytes, _ := sonic.MarshalString(roleList)
	if err = u.Redis.Set(userRoleKey, cacheBytes, rediskey.UserRoleListTTL); err != nil {
		zap.L().Error("用户角色对应关系写入缓存失败", zap.Int64("userID", userID), zap.Error(err))
	}

	// 5. 写入 Redis 映射表（每个角色映射该用户缓存key）
	for _, roleID := range roleIDList {
		roleUserDepKey := fmt.Sprintf(rediskey.RoleUserDependencies+"%d", roleID)
		if err = u.Redis.SAdd(roleUserDepKey, userRoleKey); err != nil {
			zap.L().Error("用户角色对应关系写入映射表失败", zap.Error(err))
			// 不 return，非关键失败
		}
	}

	return
}

func (u *UserService) GetUserPage(c *gin.Context, query *query.User) (resp *base.PaginationResp[response.User], err error) {
	users, total, err := u.UserRepo.GetUserPage(c, query)
	if err != nil {
		return
	}

	userResponses := make([]response.User, 0, len(users))
	for _, user := range users {
		var roles []*entity.Role

		// 从 Redis 缓存中获取用户角色
		if roles, err = u.getRoleListByUserID(c, user.ID); err != nil {
			return
		}
		user.Roles = roles

		var userResp response.User
		if err = copier.Copy(&userResp, &user); err != nil {
			zap.L().Error("复制属性失败")
			return
		}
		userResponses = append(userResponses, userResp)
	}

	resp = base.BuildPageResp[response.User](userResponses, total, query.PaginationReq)
	return
}

func (u *UserService) UpdateUserBaseInfo(c *gin.Context, userReq *request.User) (err error) {
	var userDo entity.User
	if err = copier.Copy(&userDo, userReq); err != nil {
		zap.L().Error("复制属性失败")
		return
	}
	if userReq.Email != "" {
		userDo.Email = &userReq.Email
	}
	if userReq.Phone != "" {
		userDo.Phone = &userReq.Phone
	}
	if err = u.UserRepo.CheckUserDuplicate(&userDo); err != nil {
		return
	}
	if err = u.UserRepo.UpdateUser(c, &userDo); err != nil {
		return
	}
	return
}

func (u *UserService) ChangeUserPassword(c *gin.Context, userReq *request.ChangePassword) (err error) {
	var userDo entity.User
	if err = copier.Copy(&userDo, userReq); err != nil {
		zap.L().Error("复制属性失败")
		return
	}
	user, err := u.UserRepo.GetUserById(c, userDo.ID)
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userReq.OldPassword))
	if err != nil {
		err = fmt.Errorf("旧密码错误")
		return
	}
	if userReq.NewPassword != userReq.ConfirmPassword {
		err = fmt.Errorf("两次输入的密码不一致")
		return
	}
	pwd, _ := bcrypt.GenerateFromPassword([]byte(userReq.NewPassword), bcrypt.DefaultCost)
	userDo.Password = string(pwd)
	if err = u.UserRepo.UpdateUser(c, &userDo); err != nil {
		return
	}
	return
}

func (u *UserService) UploadAvatar(c *gin.Context, filename string, src multipart.File) (fileUrl string, err error) {
	url, err := u.OssClient.UploadAvatar(c, filename, src)
	if err != nil {
		zap.L().Error("上传头像失败", zap.Error(err))
		return
	}
	var userDo entity.User
	userDo.ID = authutils.GetUserID(c)
	userDo.Avatar = url
	if err = u.UserRepo.UpdateUser(c, &userDo); err != nil {
		return
	}
	fileUrl = url
	return
}

func (u *UserService) ResetPassword(c *gin.Context, id int64) (err error) {
	var userDo entity.User
	userDo.ID = id
	password, err := bcrypt.GenerateFromPassword([]byte(u.DefaultUserConfig.ResetPassword), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("生成密码失败", zap.Error(err))
		return
	}
	userDo.Password = string(password)
	if err = u.UserRepo.UpdateUser(c, &userDo); err != nil {
		return
	}
	return
}

func (u *UserService) ChangeUserStatus(c *gin.Context, req request.ChangeStatus) (err error) {
	var userDo entity.User
	id := int64(req.ID)
	userDo.ID = id
	userDo.Status = req.Status
	if err = u.UserRepo.UpdateUser(c, &userDo); err != nil {
		return
	}
	go func() {
		// 删除用户登录状态
		id := strconv.FormatInt(id, 10)
		// 获取用户登录状态的 Redis 键
		tokenStr, _ := u.Redis.Get(id)

		// 准备需要删除的 Redis 键
		delKeys := []string{
			rediskey.SESSION + id,
			rediskey.LOGIN + tokenStr,
		}

		// 使用管道删除 Redis 中的会话和登录状态键
		err = u.Redis.PipelineDel(delKeys)
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
	}()
	return
}

//func DeleteUser(ids []int64, deleteBy int64) error {
//	// 开启事务
//	tx := global.DB.Begin()
//
//	if tx.Error != nil {
//		return tx.Error
//	}
//
//	// 更新修改者 ID
//	if err := tx.Model(&entity.User{}).Where("id IN (?)", ids).Update("updated_by", deleteBy).Error; err != nil {
//		tx.Rollback() // 回滚事务
//		zap.L().Error("更新修改者 ID 失败")
//		return err
//	}
//
//	// 删除用户
//	if err := tx.Where("id IN (?)", ids).Delete(&entity.User{}).Error; err != nil {
//		tx.Rollback() // 回滚事务
//		zap.L().Error("删除用户失败")
//		return err
//	}
//
//	// 提交事务
//	if err := tx.Commit().Error; err != nil {
//		tx.Rollback() // 回滚事务
//		zap.L().Error("提交事务失败")
//		return err
//	}
//
//	return nil
//}
//

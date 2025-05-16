package service

import (
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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"
	"time"
)

type UserService struct {
	UserRepo      *repository.UserRepository         // 用户Repo
	RoleRepo      *repository.RoleRepository         // 角色Repo
	UserRolesRepo *repository.UserRolesRepository    // 用户角色Repo
	GormTX        *repository.GormTransactionManager // gorm事务管理
	OssClient     *aliyun.OssClient                  // 阿里云OSS
	Redis         *redisx.RedisWrapper               // redis
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
	roleList, err := u.getUserRoleListFromCache(c, user.ID)
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

// getUserRoleListFromCache 尝试从 Redis 获取用户角色列表，获取不到就查数据库并写入缓存
func (u *UserService) getUserRoleListFromCache(c context.Context, userID int64) (roleList []entity.Role, error error) {
	key := fmt.Sprintf(rediskey.UserRoleList+"%d", userID)
	// 先查 Redis
	cacheStr := u.Redis.Get(key)
	if cacheStr != "" {
		if err := jsoniter.Unmarshal([]byte(cacheStr), &roleList); err == nil {
			return
		}
		// 解析失败也继续走数据库，避免缓存污染
	}
	// 缓存没有命中，从数据库查
	// 获取用户角色ID列表
	roleIDList, err := u.UserRolesRepo.GetRoleIDListByUserID(c, userID)
	if err != nil {
		return
	}
	// 根据上一步获取的用户角色ID获取角色列表
	roleList, err = u.RoleRepo.GetRoleListByRoleIDList(c, roleIDList)
	if err != nil {
		return
	}
	// 存入 Redis（最长缓存 5 分钟）
	cacheBytes, _ := jsoniter.Marshal(roleList)
	err = u.Redis.Set(key, string(cacheBytes), 5*time.Minute)
	if err != nil {
		zap.L().Error("用户角色对应关系写入缓存失败", zap.Int64("userID", userID))
		return
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
		var roles []entity.Role

		// 从 Redis 缓存中获取用户角色
		if roles, err = u.getUserRoleListFromCache(c, user.ID); err != nil {
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

func (u *UserService) UpdateUserPassword(c *gin.Context, userReq *request.ChangePassword) (err error) {
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

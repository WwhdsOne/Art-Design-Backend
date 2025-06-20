package service

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/response"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/internal/repository/db"
	"Art-Design-Backend/pkg/aliyun"
	"Art-Design-Backend/pkg/authutils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"
)

type UserService struct {
	RoleRepo          *repository.RoleRepo       // 角色Repo
	UserRepo          *repository.UserRepo       // 用户Repo
	AuthRepo          *repository.AuthRepo       // 认证Repo
	GormTX            *db.GormTransactionManager // gorm事务管理
	OssClient         *aliyun.OssClient          // 阿里云OSS
	DefaultUserConfig *config.DefaultUserConfig  // 默认用户配置
}

func (u *UserService) GetUserById(c context.Context) (res response.User, err error) {
	var user *entity.User
	id := authutils.GetUserID(c)
	// 获取用户信息
	if user, err = u.UserRepo.GetUserById(c, id); err != nil {
		return
	}
	// 获取用户角色列表
	roleList, err := u.RoleRepo.GetRoleListByUserID(c, user.ID)
	if len(roleList) == 0 {
		err = fmt.Errorf("当前用户未分配角色")
		return
	}
	user.Roles = roleList
	// 不能返回空标签，否则前端无法修改
	if len(user.Tags) == 0 {
		user.Tags = []string{}
	}
	err = copier.Copy(&res, &user)
	if err != nil {
		zap.L().Error("复制属性失败", zap.Error(err))
		return
	}
	return
}

func (u *UserService) GetUserPage(c context.Context, query *query.User) (resp *base.PaginationResp[response.User], err error) {
	users, total, err := u.UserRepo.GetUserPage(c, query)
	if err != nil {
		return
	}

	userResponses := make([]response.User, 0, len(users))
	for _, user := range users {
		var roles []*entity.Role

		// 从 Redis 缓存中获取用户角色
		if roles, err = u.RoleRepo.GetRoleListByUserID(c, user.ID); err != nil {
			return
		}
		user.Roles = roles

		var userResp response.User
		_ = copier.Copy(&userResp, &user)
		userResponses = append(userResponses, userResp)
	}

	resp = base.BuildPageResp[response.User](userResponses, total, query.PaginationReq)
	return
}

func (u *UserService) UpdateUserBaseInfo(c context.Context, userReq *request.User) (err error) {
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
	if err = u.UserRepo.CheckUserDuplicate(c, &userDo); err != nil {
		return
	}
	if err = u.UserRepo.UpdateUser(c, &userDo); err != nil {
		return
	}
	return
}

func (u *UserService) ChangeUserPassword(c context.Context, userReq *request.ChangePassword) (err error) {
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
	go func(id int64) {
		if err := u.AuthRepo.LogoutByUserID(id); err != nil {
			zap.L().Error("登出用户失败", zap.Error(err))
		}
	}(id)
	return
}

func (u *UserService) GetUserRoleBinding(c context.Context, id int64) (res *response.UserRoleBinding, err error) {
	reducedRoleList, err := u.RoleRepo.GetReducedRoleList(c)
	if err != nil {
		zap.L().Error("获取用户角色列表失败", zap.Error(err))
		return
	}
	roleList, err := u.RoleRepo.GetRoleIDListByUserID(c, id)
	if err != nil {
		zap.L().Error("获取用户角色ID列表失败", zap.Error(err))
		return
	}
	simpleRoleList := make([]*response.SimpleRole, 0, len(reducedRoleList))
	for _, role := range reducedRoleList {
		var roleResp response.SimpleRole
		roleResp.ID = role.ID
		roleResp.Name = role.Name
		simpleRoleList = append(simpleRoleList, &roleResp)
	}
	res = &response.UserRoleBinding{
		HasRoleIDs: roleList,
		Roles:      simpleRoleList,
	}
	return
}

func (u *UserService) UpdateUserRoleBinding(c *gin.Context, req *request.UserRoleBinding) (err error) {
	err = u.GormTX.Transaction(c, func(ctx context.Context) (err error) {
		if err = u.RoleRepo.DeleteUserRoleRelationsByUserID(ctx, int64(req.UserId)); err != nil {
			return
		}
		if err = u.RoleRepo.AddRolesToUser(ctx, int64(req.UserId), req.RoleIds); err != nil {
			return
		}
		return
	})
	if err != nil {
		zap.L().Error("更新用户角色绑定失败", zap.Error(err))
		return
	}
	// 缓存清理移出事务
	go func(userID int64, originalRoleIDList []int64) {
		if err := u.UserRepo.InvalidUserRoleCache(c, userID, originalRoleIDList); err != nil {
			zap.L().Error("用户变更角色信息删除缓存失败", zap.Int64("userID", userID), zap.Error(err))
		}
	}(int64(req.UserId), req.OriginalRoleIds)
	return
}

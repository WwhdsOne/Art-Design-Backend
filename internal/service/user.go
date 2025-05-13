package service

import (
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/model/request"
	"Art-Design-Backend/internal/model/resp"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/aliyun"
	"Art-Design-Backend/pkg/loginUtils"
	"Art-Design-Backend/pkg/redisx"
	"Art-Design-Backend/pkg/transaction"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"
)

type UserService struct {
	UserRepo    *repository.UserRepository // 用户Repo
	RoleRepo    *repository.RoleRepository // 角色Repo
	OssClient   *aliyun.OssClient          // 阿里云OSS
	GormSession *transaction.GormSession   // gorm事务管理
	Redis       *redisx.RedisWrapper       // redis
}

func (u *UserService) GetUserById(c *gin.Context) (res resp.User, err error) {
	var user *entity.User
	id := loginUtils.GetUserID(c)
	if id == -1 {
		err = fmt.Errorf("当前用户未登录")
		return
	}
	if user, err = u.UserRepo.GetUserById(c, id); err != nil {
		return
	}
	roleList, err := u.RoleRepo.GetRoleListByUserID(c, user.ID)
	if err != nil {
		return
	}
	if len(roleList) == 0 {
		err = fmt.Errorf("当前用户未分配角色")
		return
	}
	user.Roles = roleList
	err = copier.Copy(&res, &user)
	if err != nil {
		zap.L().Error("复制属性失败")
		return
	}
	return
}

func (u *UserService) GetUserPage(c *gin.Context, userQuery *query.User) (userPageRes *base.PaginationResp[resp.User], err error) {
	userPage, total, err := u.UserRepo.GetUserPage(c, userQuery)
	if err != nil {
		return
	}
	var pageData []resp.User
	roleCache := make(map[int64][]entity.Role) // 使用用户ID作为key，角色列表作为value

	for _, user := range userPage {
		var roleList []entity.Role

		// 先从缓存中查找是否已经查询过该用户的角色
		if cachedRoles, exist := roleCache[user.ID]; exist {
			roleList = cachedRoles
		} else {
			// 如果缓存中没有，则从数据库中查询
			if roleList, err = u.RoleRepo.GetRoleListByUserID(c, user.ID); err != nil {
				return
			}
			// 将查询结果存入缓存
			roleCache[user.ID] = roleList
		}

		user.Roles = roleList

		var pageDataResp resp.User
		if err = copier.Copy(&pageDataResp, &user); err != nil {
			zap.L().Error("复制属性失败")
			return
		}

		pageData = append(pageData, pageDataResp)
	}

	userPageRes = base.BuildPageResp[resp.User](pageData, total, userQuery.PaginationReq)

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
	userDo.ID = loginUtils.GetUserID(c)
	userDo.Avatar = url
	if err = u.UserRepo.UpdateUser(c, &userDo); err != nil {
		return
	}
	fileUrl = url
	return
}

//
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
//func UpdateUser(c *gin.Context, userReq *request.User, id int64) error {
//	//赋值给user
//	user := entity.User{
//		Username: userReq.Username,
//		Nickname: userReq.Nickname,
//	}
//	//判断密码是否为空
//	if user.Password != "" {
//		password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
//		user.Password = string(password)
//	}
//	var count int64
//	//判断重复
//	if err := global.DB.
//		WithContext(c).
//		Model(&entity.User{}).
//		Where("username = ? AND id != ?", user.Username, id).
//		Count(&count).Error; err != nil {
//		zap.L().Error("查询用户名失败", zap.Error(err))
//		return err
//	}
//	if count > 0 {
//		zap.L().Error("用户名已存在")
//		return errors.New("用户名已存在")
//	}
//	if err := global.DB.WithContext(c).Where("id = ?", id).Updates(user).Error; err != nil {
//		// 处理错误
//		zap.L().Error("更新用户失败")
//		return err
//	}
//	return nil
//}
//
////func GetUserById(id int64) (userResp resp.User, err error) {
////	var user entity.User
////	err = global.DB.Where("id = ?", id).First(&user).Error
////	if err != nil {
////		// 处理错误，例如可以返回 nil 或者记录错误日志
////		zap.L().Info("获取用户失败")
////		return
////	}
////	//组装数据
////	userResp = doToResp(user)
////	return
////}
////
////func UserPage(u *query.User) ([]resp.User, int64, error) {
////	var users []entity.User
////	// 获取总记录数
////	var total int64
////	q := global.DB.Model(&entity.User{})
////	// 获取符合条件的总记录数
////	err := q.Count(&total).Error
////	if err != nil {
////		// 处理错误，例如可以返回 nil 或者记录错误日志
////		zap.L().Error("获取用户列表失败")
////		return nil, 0, err
////	}
////	// 执行分页查询
////	global.DB.
////		Scopes(u.PaginationReq.Paginate()). // 组装分页条件
////		Find(&users)
////	//组装数据
////	var userResponses []resp.User
////	for _, user := range users {
////		//组装数据
////		userResponse := doToResp(user)
////		// 查询用户的创建人创建时间
////		var nickName string
////		// 查询用户的昵称
////		err = global.DB.Model(&entity.User{}).
////			Select("nickname").
////			Where("id = ?", userResponse.CreateBy).
////			Scan(&nickName).Error
////		if err != nil {
////			return nil, 0, err
////		}
////		userResponse.CreateUserName = nickName
////		userResponses = append(userResponses, userResponse)
////	}
////
////	return userResponses, total, nil
////}
//

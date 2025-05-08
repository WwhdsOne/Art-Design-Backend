package service

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/resp"
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/loginUtils"
	"Art-Design-Backend/pkg/redisx"
	"Art-Design-Backend/pkg/transaction"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type UserService struct {
	UserRepo    *repository.UserRepository // 用户Repo
	RoleRepo    *repository.RoleRepository // 角色Repo
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

package service

import (
	"Art-Design-Backend/global"
	"Art-Design-Backend/model/entity"
	"Art-Design-Backend/model/query"
	"Art-Design-Backend/model/request"
	"Art-Design-Backend/model/resp"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func doToResp(user entity.User) resp.User {
	return resp.User{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt,
		CreateBy:  user.CreateBy,
	}
}

func Login(ctx *gin.Context, u request.Login) (resp entity.User, err error) {
	global.DB.
		WithContext(ctx).
		Select("id", "password").
		Where("Username = ?", u.Username).
		First(&resp)
	err = bcrypt.CompareHashAndPassword([]byte(resp.Password), []byte(u.Password))
	if err != nil {
		// 密码错误
		return
	}
	return
}

func AddUser(c *gin.Context, userReq *request.User) error {
	//赋值给user
	user := entity.User{
		Username: userReq.Username,
		Nickname: userReq.Nickname,
		Password: userReq.Password,
	}
	//判断密码是否为空
	if user.Password != "" {
		password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(password)
	}
	var count int64
	if err := global.DB.
		WithContext(c).
		Model(&entity.User{}).
		Where("username = ?", user.Username).
		Count(&count).Error; err != nil {

		global.Logger.Error("查询用户名失败", zap.Error(err))
		return err
	}

	if count > 0 {
		global.Logger.Error("用户名已存在")
		return errors.New("用户名已存在")
	}
	if err := global.DB.
		WithContext(c).
		Create(&user).Error; err != nil {
		// 处理错误
		global.Logger.Error("新增用户失败")
		return err
	}
	return nil
}

func DeleteUser(ids []int64, deleteBy int64) error {
	// 开启事务
	tx := global.DB.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// 更新修改者 ID
	if err := tx.Model(&entity.User{}).Where("id IN (?)", ids).Update("updated_by", deleteBy).Error; err != nil {
		tx.Rollback() // 回滚事务
		global.Logger.Error("更新修改者 ID 失败")
		return err
	}

	// 删除用户
	if err := tx.Where("id IN (?)", ids).Delete(&entity.User{}).Error; err != nil {
		tx.Rollback() // 回滚事务
		global.Logger.Error("删除用户失败")
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // 回滚事务
		global.Logger.Error("提交事务失败")
		return err
	}

	return nil
}

func UpdateUser(c *gin.Context, userReq *request.User, id int64) error {
	//赋值给user
	user := entity.User{
		Username: userReq.Username,
		Nickname: userReq.Nickname,
	}
	//判断密码是否为空
	if user.Password != "" {
		password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(password)
	}
	var count int64
	//判断重复
	if err := global.DB.
		WithContext(c).
		Model(&entity.User{}).
		Where("username = ? AND id != ?", user.Username, id).
		Count(&count).Error; err != nil {

		global.Logger.Error("查询用户名失败", zap.Error(err))
		return err
	}
	if count > 0 {
		global.Logger.Error("用户名已存在")
		return errors.New("用户名已存在")
	}
	if err := global.DB.WithContext(c).Where("id = ?", id).Updates(user).Error; err != nil {
		// 处理错误
		global.Logger.Error("更新用户失败")
		return err
	}
	return nil
}

func GetUserById(id int64) (userResp resp.User, err error) {
	var user entity.User
	err = global.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		// 处理错误，例如可以返回 nil 或者记录错误日志
		global.Logger.Info("获取用户失败")
		return
	}
	//组装数据
	userResp = doToResp(user)
	return
}

func UserPage(u *query.User) ([]resp.User, int64, error) {
	var users []entity.User
	// 获取总记录数
	var total int64
	q := global.DB.Model(&entity.User{})
	// 获取符合条件的总记录数
	err := q.Count(&total).Error
	if err != nil {
		// 处理错误，例如可以返回 nil 或者记录错误日志
		global.Logger.Error("获取用户列表失败")
		return nil, 0, err
	}
	// 执行分页查询
	global.DB.
		Scopes(u.PaginationQ.Paginate()). // 组装分页条件
		Find(&users)
	//组装数据
	var userResponses []resp.User
	for _, user := range users {
		//组装数据
		userResponse := doToResp(user)
		// 查询用户的创建人创建时间
		var nickName string
		// 查询用户的昵称
		err = global.DB.Model(&entity.User{}).
			Select("nickname").
			Where("id = ?", userResponse.CreateBy).
			Scan(&nickName).Error
		if err != nil {
			return nil, 0, err
		}
		userResponse.CreateUserName = nickName
		userResponses = append(userResponses, userResponse)
	}

	return userResponses, total, nil
}

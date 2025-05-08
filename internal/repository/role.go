package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/constant"
	"Art-Design-Backend/pkg/errorTypes"
	"Art-Design-Backend/pkg/transaction"
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type RoleRepository struct {
	db *gorm.DB // 用户表数据库连接
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) CheckRoleDuplicate(c context.Context, role *entity.Role) (err error) {
	var result struct {
		NameExists bool
		CodeExists bool
	}

	// 检查当前记录是否有ID，如果有，则在查询中排除它
	excludeID := ""
	if role.ID != 0 {
		excludeID = fmt.Sprintf("AND id != %d", role.ID)
	}

	// 构建动态查询条件
	var query strings.Builder
	args := make([]interface{}, 0)
	conditions := make([]string, 0)

	// 只检查非空字段
	if role.Name != "" {
		conditions = append(conditions, "EXISTS(SELECT 1 FROM role WHERE name = ? "+excludeID+") AS name_exists")
		args = append(args, role.Name)
	}

	if role.Code != "" {
		conditions = append(conditions, "EXISTS(SELECT 1 FROM role WHERE code = ? "+excludeID+") AS code_exists")
		args = append(args, role.Code)
	}

	// 如果没有需要检查的字段，直接返回
	if len(conditions) == 0 {
		return nil
	}

	// 构建完整查询
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(conditions, ","))

	// 执行查询
	if err = r.db.WithContext(c).Raw(query.String(), args...).Scan(&result).Error; err != nil {
		return err
	}

	// 检查结果
	switch {
	case result.NameExists:
		return errorTypes.NewGormError("角色名称已存在")
	case result.CodeExists:
		return errorTypes.NewGormError("角色编码已存在")
	}

	return
}
func (r *RoleRepository) CreateRole(c context.Context, role *entity.Role) (err error) {
	if err = r.db.WithContext(c).Create(role).Error; err != nil {
		zap.L().Error("创建角色失败", zap.Error(err))
		return errorTypes.NewGormError("创建角色失败")
	}
	return
}

func (r *RoleRepository) GetRoleListByUserID(c context.Context, userID int64) (roleList []entity.Role, err error) {
	// 使用JOIN查询关联角色
	if err = r.db.WithContext(c).
		Table(constant.RoleTableName).
		Select("role.*").
		Joins("JOIN user_roles ON user_roles.role_id = role.id").
		Where("user_roles.user_id = ?", userID).
		Where("status = 1").
		Find(&roleList).Error; err != nil {
		zap.L().Error("查询用户角色列表失败", zap.Error(err))
		err = errorTypes.NewGormError("查询用户角色列表失败")
		return
	}
	return
}

func (r *RoleRepository) GetRoleIDListByUserID(c context.Context, userID int64) (roleIDList []int64, err error) {
	if err = r.db.WithContext(c).
		Table(constant.RoleTableName).
		Select("id").
		Joins("JOIN user_roles ON user_roles.role_id = role.id").
		Where("user_roles.user_id = ?", userID).
		Where("status = 1").
		Find(&roleIDList).Error; err != nil {
		zap.L().Error("查询角色ID列表失败")
		err = errorTypes.NewGormError("查询角色ID列表失败")
		return
	}

	return
}

func (r *RoleRepository) AssignRoleToUser(c context.Context, userID int64, roleIDList []int64) (err error) {
	// 删除原有关联
	if err = transaction.DB(c, r.db).
		Table("user_roles").
		Where("user_id = ?", userID).
		Delete(nil).Error; err != nil {
		zap.L().Error("删除原有关联失败")
		err = errorTypes.NewGormError("删除原有关联失败")
		return
	}
	if len(roleIDList) > 0 {
		userRoleList := make([]entity.UserRoles, len(roleIDList))
		for _, roleID := range roleIDList {
			userRoleList = append(userRoleList, entity.UserRoles{
				UserID: userID,
				RoleID: roleID,
			})
		}
		// 创建新的关联
		if err = transaction.DB(c, r.db).
			Table("user_roles").
			Create(userRoleList).Error; err != nil {
			zap.L().Error("创建新的关联失败")
			err = errorTypes.NewGormError("创建新的关联失败")
		}
	}
	return
}

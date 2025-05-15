package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/pkg/constant"
	"Art-Design-Backend/pkg/errors"
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
	var queryCondition strings.Builder
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
	queryCondition.WriteString("SELECT ")
	queryCondition.WriteString(strings.Join(conditions, ","))

	// 执行查询
	if err = r.db.WithContext(c).Raw(queryCondition.String(), args...).Scan(&result).Error; err != nil {
		return err
	}

	// 检查结果
	switch {
	case result.NameExists:
		return errors.NewDBError("角色名称已存在")
	case result.CodeExists:
		return errors.NewDBError("角色编码已存在")
	}

	return
}
func (r *RoleRepository) CreateRole(c context.Context, role *entity.Role) (err error) {
	if err = DB(c, r.db).Create(role).Error; err != nil {
		zap.L().Error("创建角色失败", zap.Error(err))
		return errors.NewDBError("创建角色失败")
	}
	return
}

func (r *RoleRepository) GetRoleListByUserID(c context.Context, userID int64) (roleList []entity.Role, err error) {
	// 使用JOIN查询关联角色
	if err = DB(c, r.db).
		Table(constant.RoleTableName).
		Select("role.*").
		Joins("JOIN user_roles ON user_roles.role_id = role.id").
		Where("user_roles.user_id = ?", userID).
		Where("status = 1").
		Find(&roleList).Error; err != nil {
		zap.L().Error("查询用户角色列表失败", zap.Error(err))
		err = errors.NewDBError("查询用户角色列表失败")
		return
	}
	return
}

func (r *RoleRepository) GetRoleIDListByUserID(c context.Context, userID int64) (roleIDList []int64, err error) {
	if err = DB(c, r.db).
		Table(constant.RoleTableName).
		Select("id").
		Joins("JOIN user_roles ON user_roles.role_id = role.id").
		Where("user_roles.user_id = ?", userID).
		Where("status = 1").
		Find(&roleIDList).Error; err != nil {
		zap.L().Error("查询角色ID列表失败")
		err = errors.NewDBError("查询角色ID列表失败")
		return
	}

	return
}

func (r *RoleRepository) AssignRoleToUser(c context.Context, userID int64, roleIDList []int64) (err error) {
	// 删除原有关联
	if err = DB(c, r.db).
		Table("user_roles").
		Where("user_id = ?", userID).
		Delete(nil).Error; err != nil {
		zap.L().Error("删除原有关联失败")
		err = errors.NewDBError("删除原有关联失败")
		return
	}
	if len(roleIDList) > 0 {
		userRoleList := make([]entity.UserRoles, 0, len(roleIDList))
		for _, roleID := range roleIDList {
			userRoleList = append(userRoleList, entity.UserRoles{
				UserID: userID,
				RoleID: roleID,
			})
		}
		// 创建新的关联
		if err = DB(c, r.db).
			Table("user_roles").
			Create(userRoleList).Error; err != nil {
			zap.L().Error("创建新的关联失败")
			err = errors.NewDBError("创建新的关联失败")
		}
	}
	return
}

func (r *RoleRepository) GetRolePage(c context.Context, role *query.Role) (rolePage []*entity.Role, total int64, err error) {
	db := DB(c, r.db)

	// 构建通用查询条件
	queryConditions := db.Table(constant.RoleTableName)

	if role.Name != "" {
		queryConditions = queryConditions.Where("name LIKE ?", "%"+role.Name+"%")
	}

	// 查询总数
	if err = queryConditions.Count(&total).Error; err != nil {
		zap.L().Error("获取角色分页失败")
		err = errors.NewDBError("获取角色分页失败")
		return
	}

	// 查询分页数据
	if err = queryConditions.Scopes(role.Paginate()).Find(&rolePage).Error; err != nil {
		zap.L().Error("获取角色分页数据失败")
		err = errors.NewDBError("获取角色分页数据失败")
		return
	}
	return
}

func (r *RoleRepository) UpdateRole(c context.Context, role *entity.Role) (err error) {
	if err = DB(c, r.db).Updates(role).Error; err != nil {
		zap.L().Error("更新角色失败")
		err = errors.NewDBError("更新角色失败")
		return
	}
	return
}

func (r *RoleRepository) DeleteRoleByID(c context.Context, roleID int64) (err error) {
	if err = DB(c, r.db).Delete(&entity.Role{}, roleID).Error; err != nil {
		zap.L().Error("删除角色失败")
		err = errors.NewDBError("删除角色失败")
		return
	}
	return
}

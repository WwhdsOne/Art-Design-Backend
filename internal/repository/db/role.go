package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/pkg/errors"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type RoleDB struct {
	db *gorm.DB // 用户表数据库连接
}

func NewRoleDB(db *gorm.DB) *RoleDB {
	return &RoleDB{
		db: db,
	}
}

func (r *RoleDB) CheckRoleDuplicate(c context.Context, role *entity.Role) (err error) {
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
		conditions = append(conditions, "EXISTS(SELECT 1 FROM \"role\" WHERE \"name\" = ? "+excludeID+") AS name_exists")
		args = append(args, role.Name)
	}

	if role.Code != "" {
		conditions = append(conditions, "EXISTS(SELECT 1 FROM \"role\" WHERE \"code\" = ? "+excludeID+") AS code_exists")
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
	if err = DB(c, r.db).Raw(queryCondition.String(), args...).Scan(&result).Error; err != nil {
		return err
	}

	// 检查结果
	switch {
	case result.NameExists:
		return errors.NewDBError("角色名称已存在")
	case result.CodeExists:
		return errors.NewDBError("角色编码已存在")
	}

	return nil
}
func (r *RoleDB) CreateRole(c context.Context, role *entity.Role) (err error) {
	if err = DB(c, r.db).Create(role).Error; err != nil {
		return errors.WrapDBError(err, "创建角色失败")
	}
	return
}

// GetEnableRoleByID 查询有效角色
func (r *RoleDB) GetEnableRoleByID(c context.Context, roleID int64) (role *entity.Role, err error) {
	if err = DB(c, r.db).
		Where("id = ?", roleID).
		Where("status = 1").
		First(&role).Error; err != nil {
		err = errors.NewDBError("查询角色失败")
		return
	}
	return
}

func (r *RoleDB) GetRolePage(c context.Context, role *query.Role) (rolePage []*entity.Role, total int64, err error) {
	db := DB(c, r.db)

	// 构建通用查询条件
	queryConditions := db.Model(&entity.Role{})

	if role.Name != "" {
		queryConditions = queryConditions.Where("name LIKE ?", "%"+role.Name+"%")
	}

	// 查询总数
	if err = queryConditions.Count(&total).Error; err != nil {
		err = errors.WrapDBError(err, "获取角色分页失败")
		return
	}

	// 查询分页数据
	if err = queryConditions.Scopes(role.Paginate()).Find(&rolePage).Error; err != nil {
		err = errors.WrapDBError(err, "获取角色分页数据失败")
		return
	}
	return
}

func (r *RoleDB) UpdateRole(c context.Context, role *entity.Role) (err error) {
	if err = DB(c, r.db).Updates(role).Error; err != nil {
		err = errors.WrapDBError(err, "更新角色失败")
		return
	}
	return
}

func (r *RoleDB) DeleteRoleByID(c context.Context, roleID int64) (err error) {
	if err = DB(c, r.db).Where("id = ?", roleID).Delete(nil).Error; err != nil {
		err = errors.WrapDBError(err, "删除角色失败")
		return
	}
	return
}

func (r *RoleDB) FilterValidRoleIDs(ctx context.Context, roleIDs []int64) (validIDList []int64, err error) {
	err = DB(ctx, r.db).
		Model(&entity.Role{}).
		Where("id IN ? AND status = ?", roleIDs, 1). // 1 = 启用状态
		Pluck("id", &validIDList).Error
	if err != nil {
		err = errors.WrapDBError(err, "查询有效角色失败")
		return
	}
	return
}

func (r *RoleDB) GetReducedRoleList(ctx context.Context) (roleList []*entity.Role, err error) {
	if err = DB(ctx, r.db).
		Select("id", "name").
		Where("status = 1").
		Find(&roleList).Error; err != nil {
		err = errors.NewDBError("获取精简角色列表失败")
		return
	}
	return
}

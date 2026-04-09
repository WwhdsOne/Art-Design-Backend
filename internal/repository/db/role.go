package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/internal/repository/dupcheck"
	"Art-Design-Backend/pkg/errors"
	"context"

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

func (r *RoleDB) CheckRoleDuplicate(c context.Context, role *entity.Role) error {
	return dupcheck.Check(DB(c, r.db), "\"role\"", role.ID, []dupcheck.Field{
		{Column: "\"name\"", ErrMsg: "角色名称已存在", Value: role.Name},
		{Column: "\"code\"", ErrMsg: "角色编码已存在", Value: role.Code},
	})
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

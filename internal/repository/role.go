package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/errors"
	"Art-Design-Backend/pkg/redisx"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type RoleRepository struct {
	db    *gorm.DB             // 用户表数据库连接
	redis *redisx.RedisWrapper // redis
}

func NewRoleRepository(db *gorm.DB, redis *redisx.RedisWrapper) *RoleRepository {
	return &RoleRepository{
		db:    db,
		redis: redis,
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

func (r *RoleRepository) GetRoleByID(c context.Context, roleID int64) (role *entity.Role, err error) {
	key := strconv.FormatInt(roleID, 10)
	// 从 Redis 读取
	roleJson, err := r.redis.Get(rediskey.RoleInfo + key)
	if err == nil {
		_ = sonic.UnmarshalString(roleJson, &role)
		return
	}
	if err = DB(c, r.db).
		Where("id = ?", roleID).
		Where("status = 1").
		First(&role).Error; err != nil {
		err = errors.NewDBError("查询角色失败")
		return
	}
	// 异步写入 Redis
	go func() {
		roleJson, _ = sonic.MarshalString(role)
		err = r.redis.Set(rediskey.RoleInfo+key, roleJson, rediskey.RoleInfoTTL)
	}()
	return
}

func (r *RoleRepository) GetRoleListByRoleIDList(c context.Context, roleIDList []int64) (roleList []entity.Role, err error) {
	for _, roleID := range roleIDList {
		var role *entity.Role
		role, err = r.GetRoleByID(c, roleID)
		if err != nil {
			return
		}
		roleList = append(roleList, *role)
	}
	return
}

func (r *RoleRepository) GetRolePage(c context.Context, role *query.Role) (rolePage []*entity.Role, total int64, err error) {
	db := DB(c, r.db)

	// 构建通用查询条件
	queryConditions := db.Model(&entity.Role{})

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

// InvalidRoleInfoCache 删除角色信息缓存
// 同时也删除映射表缓存
func (r *RoleRepository) InvalidRoleInfoCache(roleID int64) (err error) {
	// 删除角色信息缓存
	key := rediskey.RoleInfo + strconv.FormatInt(roleID, 10)
	err = r.redis.Del(key)
	if err != nil {
		err = errors.NewCacheError("删除角色信息缓存失败")
		return
	}
	// 根据角色用户映射表删除用户角色信息缓存
	key = fmt.Sprintf(rediskey.RoleUserDependencies+"%d", roleID)
	err = r.redis.DelBySetMembers(key)
	if err != nil {
		err = errors.NewCacheError("删除角色菜单映射表缓存失败")
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
	if err = DB(c, r.db).Where("id = ?", roleID).Delete(nil).Error; err != nil {
		zap.L().Error("删除角色失败")
		err = errors.NewDBError("删除角色失败")
		return
	}
	return
}

func (r *RoleRepository) FilterValidRoleIDs(ctx context.Context, roleIDs []int64) (validIDList []int64, err error) {
	err = DB(ctx, r.db).
		Model(&entity.Role{}).
		Where("id IN ? AND status = ?", roleIDs, 1). // 1 = 启用状态
		Pluck("id", &validIDList).Error
	if err != nil {
		zap.L().Error("查询有效角色失败", zap.Error(err))
		err = errors.NewDBError("查询有效角色失败")
		return
	}

	return
}

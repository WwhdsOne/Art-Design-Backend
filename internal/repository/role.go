package repository

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/constant"
	"Art-Design-Backend/pkg/errorTypes"
	"Art-Design-Backend/pkg/utils"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB // 用户表数据库连接
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: db.Table(constant.RoleTableName),
	}
}

func (r *RoleRepository) CreateRole(c context.Context, role *entity.Role) (err error) {
	if err = r.db.WithContext(c).Create(role).Error; err != nil {
		// 检查是否是唯一约束冲突错误
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			// 从MySQL错误信息中提取冲突字段
			// MySQL错误格式: "Error 1062: Duplicate entry 'value' for key 'field_name'"
			if field := utils.ExtractMySQLUniqueField(err.Error()); field != "" {
				zap.L().Error("创建角色失败: 字段 "+field+" 已存在", zap.Error(err))
				return errorTypes.NewGormError("创建角色失败: 字段 " + field + " 已存在")
			}
		}
		zap.L().Error("创建角色失败", zap.Error(err))
		return errorTypes.NewGormError("创建角色失败")
	}
	return
}

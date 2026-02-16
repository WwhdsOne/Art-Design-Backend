package db

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/internal/model/query"
	"Art-Design-Backend/pkg/errors"
	"context"

	"gorm.io/gorm"
)

type OperationLogDB struct {
	db *gorm.DB // 用户表数据库连接
}

func NewOperationLogDB(db *gorm.DB) *OperationLogDB {
	return &OperationLogDB{
		db: db,
	}
}

func (o *OperationLogDB) CreateOperationLog(c context.Context, operationLog *entity.OperationLog) (err error) {
	if err = DB(c, o.db).Create(operationLog).Error; err != nil {
		err = errors.WrapDBError(err, "创建操作日志失败")
		return
	}
	return
}

func (o *OperationLogDB) GetOperationLogPage(
	c context.Context,
	log *query.OperationLog,
) (logPage []*entity.OperationLog, total int64, err error) {

	db := DB(c, o.db)

	queryConditions := db.Model(&entity.OperationLog{})

	// ====== 条件过滤 ======

	if log.OperatorID != 0 {
		queryConditions = queryConditions.Where("operator_id = ?", log.OperatorID)
	}

	if log.Method != "" {
		queryConditions = queryConditions.Where("method = ?", log.Method)
	}

	if log.Path != "" {
		queryConditions = queryConditions.Where("path LIKE ?", "%"+log.Path+"%")
	}

	if log.Status != 0 {
		queryConditions = queryConditions.Where("status = ?", log.Status)
	}

	if log.IP != "" {
		queryConditions = queryConditions.Where("ip LIKE ?", "%"+log.IP+"%")
	}

	if log.Browser != "" {
		queryConditions = queryConditions.Where("browser = ?", log.Browser)
	}

	if log.OS != "" {
		queryConditions = queryConditions.Where("os = ?", log.OS)
	}

	// ====== 时间区间查询（重点） ======
	if log.StartTime != nil {
		queryConditions = queryConditions.Where("created_at >= ?", *log.StartTime)
	}

	if log.EndTime != nil {
		queryConditions = queryConditions.Where("created_at <= ?", *log.EndTime)
	}

	// ====== 查询总数 ======
	if err = queryConditions.Count(&total).Error; err != nil {
		err = errors.WrapDBError(err, "获取操作日志分页总数失败")
		return
	}

	// ====== 分页查询 ======
	if err = queryConditions.
		Order("id DESC").
		Scopes(log.Paginate()).
		Find(&logPage).Error; err != nil {

		err = errors.WrapDBError(err, "获取操作日志分页数据失败")
		return
	}

	return
}

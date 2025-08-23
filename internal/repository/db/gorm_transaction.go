package db

import (
	"Art-Design-Backend/pkg/errors"
	"context"

	"gorm.io/gorm"
)

type GormTransactionManager struct {
	db *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) *GormTransactionManager {
	return &GormTransactionManager{
		db: db,
	}
}

type dbKey struct{}

func DB(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	db := ctx.Value(dbKey{})
	if db == nil {
		return fallback.WithContext(ctx)
	}
	return db.(*gorm.DB).WithContext(ctx)
}

func (s *GormTransactionManager) Transaction(ctx context.Context, f func(context.Context) error) (err error) {
	tx := DB(ctx, s.db).Begin()
	if tx.Error != nil {
		return errors.WrapDBError(tx.Error, "开启事务失败")
	}
	c := context.WithValue(ctx, dbKey{}, tx)
	err = f(c)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		return errors.WrapDBError(err, "提交事务失败")
	}
	return
}

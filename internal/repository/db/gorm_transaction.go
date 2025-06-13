package db

import (
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

func (s *GormTransactionManager) Transaction(ctx context.Context, f func(context.Context) error) error {
	tx := DB(ctx, s.db).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	c := context.WithValue(ctx, dbKey{}, tx)
	err := f(c)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

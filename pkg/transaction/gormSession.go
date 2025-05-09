package transaction

import (
	"context"
	"gorm.io/gorm"
)

type GormSession struct {
	db  *gorm.DB
	ctx context.Context
}

func NewGormSession(db *gorm.DB) *GormSession {
	return &GormSession{
		db:  db,
		ctx: context.Background(),
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

func (s *GormSession) Transaction(ctx context.Context, f func(context.Context) error) error {
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

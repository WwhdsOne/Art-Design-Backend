package transaction

import (
	"context"
	"gorm.io/gorm"
)

type Session interface {
	Begin(ctx context.Context) (Session, error)
	Transaction(ctx context.Context, f func(context.Context) error) error
	Rollback() error
	Commit() error
	Context() context.Context
}

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

func (s *GormSession) Begin(ctx context.Context) (Session, error) {
	tx := DB(ctx, s.db).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &GormSession{
		ctx: context.WithValue(ctx, dbKey{}, tx),
		db:  tx,
	}, nil
}

func (s *GormSession) Rollback() error {
	return s.db.Rollback().Error
}

func (s *GormSession) Commit() error {
	return s.db.Commit().Error
}

func (s *GormSession) Context() context.Context {
	return s.ctx
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

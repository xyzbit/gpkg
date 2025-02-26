package ctxwrap

import (
	"context"

	"gorm.io/gorm"
)

type gormDBKey struct{}

func NewGormDBContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, gormDBKey{}, db)
}

func FromGormDBContext(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(gormDBKey{}).(*gorm.DB)
	if !ok {
		return nil
	}
	return db
}

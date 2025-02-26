package gormx

import (
	"context"

	"github.com/xyzbit/gpkg/ctxwrap"
	"gorm.io/gorm"
)

type DBMaker interface {
	DB(ctx context.Context) *gorm.DB
}

// Transaction 封装事务方法，service 层屏蔽具体 gorm 对象
//
// 注意：方法中的 context 需要使用 txCtx 作为入参
func Transaction(ctx context.Context, maker DBMaker, fc func(txCtx context.Context) error) error {
	db := maker.DB(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		ctx = ctxwrap.NewGormDBContext(ctx, tx)
		return fc(ctx)
	})
}

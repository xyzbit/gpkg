package gormx

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"codeup.aliyun.com/qimao/pkg/contrib/biz/ctxwrap"
	"codeup.aliyun.com/qimao/pkg/contrib/driver/mysql"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func (r *Repo) DB(ctx context.Context) *gorm.DB {
	upstreamDB := ctxwrap.FromGormDBContext(ctx)
	if upstreamDB != nil {
		return upstreamDB
	}
	return r.db.WithContext(ctx)
}

func (r *Repo) Reset() {
	_ = repo.db.Exec("DROP TABLE test_tx").Error
	_ = repo.db.AutoMigrate(&TestTx{})
}

func (r *Repo) Create(ctx context.Context, tx *TestTx) (*TestTx, error) {
	err := r.DB(ctx).Create(tx).Error
	return tx, err
}

func (r *Repo) Update(ctx context.Context, tx *TestTx) error {
	return r.DB(ctx).Where("id = ?", tx.ID).Updates(tx).Error
}

func (r *Repo) Tx(ctx context.Context) error {
	return Transaction(ctx, repo, func(txCtx context.Context) error {
		t, err := r.Create(txCtx, &TestTx{Content: "test"})
		if err != nil {
			return err
		}

		t.Content = t.Content + "update"
		err = r.Update(txCtx, t)
		if err != nil {
			return err
		}

		return nil
	})
}

type TestTx struct {
	ID          int64     `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement:true;comment:主键ID" json:"id"`                                                // 主键ID
	Content     string    `gorm:"column:content;type:mediumtext;not null;comment:配置项内容" json:"content"`                                                               // 配置项内容
	Creator     string    `gorm:"column:creator;type:varchar(50);not null;comment:创建者" json:"creator"`                                                                // 创建者
	Operator    string    `gorm:"column:operator;type:varchar(50);not null;comment:最近一次操作者姓名" json:"operator"`                                                        // 最近一次操作者姓名
	CreatedTime time.Time `gorm:"column:created_time;type:datetime;not null;index:idx_created,priority:1;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_time"` // 创建时间
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime;not null;index:idx_updated,priority:1;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_time"` // 更新时间
}

func (*TestTx) TableName() string {
	return "test_tx"
}

var (
	repo *Repo
)

func TestMain(m *testing.M) {
	db, err := mysql.NewGorm(mysql.Config{
		DSN:         "",
		MaxIdle:     100,
		MaxOpen:     100,
		MaxLifeTime: 7200,
		MaxIdleTime: 3600,
		IsDebug:     true,
	})
	if err != nil {
		panic(err)
	}

	repo = &Repo{
		db: db,
	}

	exit := m.Run()
	os.Exit(exit)
}

func TestTransaction(t *testing.T) {
	repo.Reset()

	ctx := context.Background()
	err := Transaction(ctx, repo, func(txCtx context.Context) error {
		tx, err := repo.Create(txCtx, &TestTx{Content: "test"})
		if err != nil {
			return err
		}

		tx.Content = tx.Content + "update"
		err = repo.Update(txCtx, tx)
		if err != nil {
			return err
		}
		return errors.New("test rollback")
	})
	if err != nil {
		t.Log(err)
	}
}

func TestTransactionNested(t *testing.T) {
	repo.Reset()

	ctx := context.Background()
	err := Transaction(ctx, repo, func(txCtx context.Context) error {
		tx, err := repo.Create(txCtx, &TestTx{Content: "test"})
		if err != nil {
			return err
		}

		tx.Content = tx.Content + "update"
		err = repo.Update(txCtx, tx)
		if err != nil {
			return err
		}

		err = repo.Tx(txCtx)
		if err != nil {
			return err
		}

		//return nil
		return errors.New("test nested rollback")
	})
	if err != nil {
		t.Log(err)
	}
}

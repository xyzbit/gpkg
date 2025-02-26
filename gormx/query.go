package gormx

import (
	"strings"

	"github.com/bytedance/sonic"
	"gorm.io/gorm"
)

type Kind string

const (
	kindSelect      Kind = "Select"
	kindEqual       Kind = "Equal"
	kindNot         Kind = "Not"
	kindIn          Kind = "In"
	kindNotIn       Kind = "NotIn"
	kindGt          Kind = "Gt"
	kindGte         Kind = "Gte"
	kindLt          Kind = "Lt"
	kindLte         Kind = "Lte"
	kindLike        Kind = "Like"
	kindBetween     Kind = "Between"
	kindOr          Kind = "Or"
	kindIsNull      Kind = "IsNull"
	kindNotNull     Kind = "NotNull"
	kindOrderBy     Kind = "OrderBy"
	kindCustomOrder Kind = "CustomOrder"
	kindLimit       Kind = "Limit"
	kindOffset      Kind = "Offset"
	kindGroup       Kind = "Group"
)

var execOrder = []Kind{
	kindSelect,
	kindEqual,
	kindNot,
	kindIn,
	kindNotIn,
	kindGt,
	kindGte,
	kindLt,
	kindLte,
	kindLike,
	kindBetween,
	kindIsNull,
	kindNotNull,
	kindOr,
	kindGroup,
	kindOrderBy,
	kindCustomOrder,
	kindLimit,
	kindOffset,
}

type Query struct {
	conMap map[Kind]Condition
}

func NewQuery() *Query {
	return &Query{
		conMap: make(map[Kind]Condition, len(execOrder)),
	}
}

type Page interface {
	GetOrderBy() string
	GetPage() uint64
	GetPageSize() uint64
}

func (q *Query) WithDB(db *gorm.DB) *gorm.DB {
	for _, k := range execOrder {
		cond, ok := q.conMap[k]
		if ok {
			db = cond.Do(db)
		}
	}
	return db
}

func (q *Query) Eq(key string, value any) *Query {
	cond, ok := q.conMap[kindEqual]
	if ok {
		cond.Build(key, value)
		return q
	}

	q.conMap[kindEqual] = newEqualCond().Build(key, value)
	return q
}

func (q *Query) Not(key string, value any) *Query {
	cond, ok := q.conMap[kindNot]
	if ok {
		cond.Build(key, value)
		return q
	}

	q.conMap[kindNot] = newNotCond().Build(key, value)
	return q
}

func (q *Query) In(key string, value any) *Query {
	cond, ok := q.conMap[kindIn]
	if ok {
		cond.Build(key, value)
		return q
	}

	q.conMap[kindIn] = newInCond().Build(key, value)
	return q
}

func (q *Query) NotIn(key string, value any) *Query {
	cond, ok := q.conMap[kindNotIn]
	if ok {
		cond.Build(key, value)
		return q
	}

	q.conMap[kindNotIn] = newNotInCond().Build(key, value)
	return q
}

func (q *Query) Gt(key string, value any) *Query {
	cond, ok := q.conMap[kindGt]
	if ok {
		cond.Build(key, value)
		return q
	}

	q.conMap[kindGt] = newGtCond().Build(key, value)
	return q
}

func (q *Query) Gte(key string, value any) *Query {
	cond, ok := q.conMap[kindGte]
	if ok {
		cond.Build(key, value)
		return q
	}

	q.conMap[kindGte] = newGteCond().Build(key, value)
	return q
}

func (q *Query) Lt(key string, value any) *Query {
	cond, ok := q.conMap[kindLt]
	if ok {
		cond.Build(key, value)
		return q
	}

	q.conMap[kindLt] = newLtCond().Build(key, value)
	return q
}

func (q *Query) Lte(key string, value any) *Query {
	cond, ok := q.conMap[kindLte]
	if ok {
		cond.Build(key, value)
		return q
	}

	q.conMap[kindLte] = newLteCond().Build(key, value)
	return q
}

func (q *Query) Like(key string, value any) *Query {
	cond, ok := q.conMap[kindLike]
	if ok {
		cond.Build(key, value)
		return q
	}

	q.conMap[kindLike] = newLikeCond().Build(key, value)
	return q
}

func (q *Query) LikeWithFunction(key string, value any, fc Function) *Query {
	cond, ok := q.conMap[kindLike]
	if ok {
		cond.Build(key, value, fc)
		return q
	}

	q.conMap[kindLike] = newLikeCond().Build(key, value, fc)
	return q
}

func (q *Query) Between(key string, lower, upper any) *Query {
	cond, ok := q.conMap[kindBetween]
	if ok {
		cond.Build(key, lower, upper)
		return q
	}

	q.conMap[kindBetween] = newBetweenCond().Build(key, lower, upper)
	return q
}

func (q *Query) Or(key string, value any) *Query {
	cond, ok := q.conMap[kindOr]
	if ok {
		cond.Build(key, value)
		return q
	}

	q.conMap[kindOr] = newOrCond().Build(key, value)
	return q
}

func (q *Query) IsNull(key string) *Query {
	cond, ok := q.conMap[kindIsNull]
	if ok {
		cond.Build(key)
		return q
	}

	q.conMap[kindIsNull] = newIsNullCond().Build(key)
	return q
}

func (q *Query) NotNull(key string) *Query {
	cond, ok := q.conMap[kindNotNull]
	if ok {
		cond.Build(key)
		return q
	}

	q.conMap[kindNotNull] = newNotNullCond().Build(key)
	return q
}

func (q *Query) Page(page Page) *Query {
	if page == nil {
		return q
	}

	pageNum, pageSize, orderBy := page.GetPage(), page.GetPageSize(), page.GetOrderBy()

	if pageNum > 0 {
		offset := (pageNum - 1) * pageSize
		q = q.Offset(int(offset))
	}
	if pageSize > 0 {
		q = q.Limit(int(pageSize))
	}

	if orderBy != "" {
		orderByMap := make(map[string]string)
		if err := sonic.UnmarshalString(orderBy, &orderByMap); err != nil {
			return q
		}

		for flied, order := range orderByMap { // descend|ascend => desc|asc
			q = q.OrderBy(flied + " " + strings.TrimSuffix(order, "end"))
			break // current only support one.
		}
	}

	return q
}

func (q *Query) OrderBy(key string) *Query {
	cond, ok := q.conMap[kindOrderBy]
	if ok {
		cond.Build(key)
		return q
	}

	q.conMap[kindOrderBy] = newOrderByCond().Build(key)
	return q
}

func (q *Query) Limit(limit int) *Query {
	cond, ok := q.conMap[kindLimit]
	if ok {
		cond.Build("", limit)
		return q
	}

	q.conMap[kindLimit] = newLimitCond().Build("", limit)
	return q
}

func (q *Query) Offset(offset int) *Query {
	cond, ok := q.conMap[kindOffset]
	if ok {
		cond.Build("", offset)
		return q
	}

	q.conMap[kindOffset] = newOffsetCond().Build("", offset)
	return q
}

func (q *Query) CustomOrder(order string, defaultOrder string) *Query {
	cond, ok := q.conMap[kindCustomOrder]
	if ok {
		cond.Build("", order, defaultOrder)
		return q
	}

	q.conMap[kindCustomOrder] = newCustomOrderCond().Build("", order, defaultOrder)
	return q
}

func (q *Query) Select(fields ...any) *Query {
	cond, ok := q.conMap[kindSelect]
	if ok {
		cond.Build("", fields...)
		return q
	}

	q.conMap[kindSelect] = newSelectCond().Build("", fields...)
	return q
}

func (q *Query) Group(fields ...any) *Query {
	cond, ok := q.conMap[kindGroup]
	if ok {
		cond.Build("", fields...)
		return q
	}

	q.conMap[kindGroup] = newGroupCond().Build("", fields...)
	return q
}

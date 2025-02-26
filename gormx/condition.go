package gormx

import (
	"fmt"
	"strings"

	"codeup.aliyun.com/qimao/pkg/contrib/core/convertor"
	"github.com/bytedance/sonic"
	"gorm.io/gorm"
)

type Condition interface {
	Build(key string, args ...any) Condition
	Do(db *gorm.DB) *gorm.DB
}

type equalCond struct {
	body map[string]any
}

func newEqualCond() *equalCond {
	return &equalCond{
		body: make(map[string]any),
	}
}

func (e *equalCond) Build(key string, args ...any) Condition {
	e.body[key] = args[0]
	return e
}

func (e *equalCond) Do(db *gorm.DB) *gorm.DB {
	for k, v := range e.body {
		db = db.Where(fmt.Sprintf("`%s` = ?", k), v)
	}
	return db
}

type notCond struct {
	body map[string]any
}

func newNotCond() *notCond {
	return &notCond{
		body: make(map[string]any),
	}
}

func (n *notCond) Build(key string, args ...any) Condition {
	n.body[key] = args[0]
	return n
}

func (n *notCond) Do(db *gorm.DB) *gorm.DB {
	for k, v := range n.body {
		db = db.Not(fmt.Sprintf("`%s` = ?", k), v)
	}
	return db
}

type inCond struct {
	body map[string]any
}

func newInCond() *inCond {
	return &inCond{
		body: make(map[string]any),
	}
}

func (i *inCond) Build(key string, args ...any) Condition {
	i.body[key] = args[0]
	return i
}

func (i *inCond) Do(db *gorm.DB) *gorm.DB {
	for k, v := range i.body {
		db = db.Where(fmt.Sprintf("`%s` IN ?", k), v)
	}
	return db
}

type notInCond struct {
	body map[string]any
}

func newNotInCond() *notInCond {
	return &notInCond{
		body: make(map[string]any),
	}
}

func (i *notInCond) Build(key string, args ...any) Condition {
	i.body[key] = args[0]
	return i
}

func (i *notInCond) Do(db *gorm.DB) *gorm.DB {
	for k, v := range i.body {
		db = db.Where(fmt.Sprintf("`%s` NOT IN ?", k), v)
	}
	return db
}

type gtCond struct {
	body map[string]any
}

func newGtCond() *gtCond {
	return &gtCond{
		body: make(map[string]any),
	}
}

func (g *gtCond) Build(key string, args ...any) Condition {
	g.body[key] = args[0]
	return g
}

func (g *gtCond) Do(db *gorm.DB) *gorm.DB {
	for k, v := range g.body {
		db = db.Where(fmt.Sprintf("`%s` > ?", k), v)
	}
	return db
}

type gteCond struct {
	body map[string]any
}

func newGteCond() *gteCond {
	return &gteCond{
		body: make(map[string]any),
	}
}

func (g *gteCond) Build(key string, args ...any) Condition {
	g.body[key] = args[0]
	return g
}

func (g *gteCond) Do(db *gorm.DB) *gorm.DB {
	for k, v := range g.body {
		db = db.Where(fmt.Sprintf("`%s` >= ?", k), v)
	}
	return db
}

type ltCond struct {
	body map[string]any
}

func newLtCond() *ltCond {
	return &ltCond{
		body: make(map[string]any),
	}
}

func (l *ltCond) Build(key string, args ...any) Condition {
	l.body[key] = args[0]
	return l
}

func (l *ltCond) Do(db *gorm.DB) *gorm.DB {
	for k, v := range l.body {
		db = db.Where(fmt.Sprintf("`%s` < ?", k), v)
	}
	return db
}

type lteCond struct {
	body map[string]any
}

func newLteCond() *lteCond {
	return &lteCond{
		body: make(map[string]any),
	}
}

func (l *lteCond) Build(key string, args ...any) Condition {
	l.body[key] = args[0]
	return l
}

func (l *lteCond) Do(db *gorm.DB) *gorm.DB {
	for k, v := range l.body {
		db = db.Where(fmt.Sprintf("`%s` <= ?", k), v)
	}
	return db
}

type likeCond struct {
	body map[string]any
}

type likeVal struct {
	val any
	fc  Function
}

func newLikeCond() *likeCond {
	return &likeCond{
		body: make(map[string]any),
	}
}

func (l *likeCond) Build(key string, args ...any) Condition {
	if len(args) == 2 {
		fc, ok := args[1].(Function)
		if !ok {
			panic("args[1] must be Function")
		}
		l.body[key] = likeVal{
			val: args[0],
			fc:  fc,
		}
	} else {
		l.body[key] = likeVal{
			val: args[0],
			fc:  nil,
		}
	}
	return l
}

func (l *likeCond) Do(db *gorm.DB) *gorm.DB {
	for k, v := range l.body {
		lv := v.(likeVal)
		s := convertor.ToString(lv.val)
		if lv.fc == nil {
			db = db.Where(fmt.Sprintf("`%s` LIKE ?", k), "%"+s+"%")
		} else {
			db = db.Where(fmt.Sprintf("%s LIKE ?", lv.fc.Expression(k)), "%"+lv.fc.ConvertVal(s)+"%")
		}
	}
	return db
}

type betweenCond struct {
	body map[string][]any
}

func newBetweenCond() *betweenCond {
	return &betweenCond{
		body: make(map[string][]any),
	}
}

func (b *betweenCond) Build(key string, args ...any) Condition {
	b.body[key] = args
	return b
}

func (b *betweenCond) Do(db *gorm.DB) *gorm.DB {
	for k, v := range b.body {
		db = db.Where(fmt.Sprintf("`%s` BETWEEN ? AND ?", k), v[0], v[1])
	}
	return db
}

type isNullCond struct {
	body map[string]any
}

func newIsNullCond() *isNullCond {
	return &isNullCond{
		body: make(map[string]any),
	}
}

func (i *isNullCond) Build(key string, args ...any) Condition {
	i.body[key] = args[0]
	return i
}

func (i *isNullCond) Do(db *gorm.DB) *gorm.DB {
	for k, v := range i.body {
		db = db.Where(fmt.Sprintf("`%s` IS NULL", k), v)
	}
	return db
}

type notNullCond struct {
	body map[string]any
}

func newNotNullCond() *notNullCond {
	return &notNullCond{
		body: make(map[string]any),
	}
}

func (i *notNullCond) Build(key string, args ...any) Condition {
	i.body[key] = args[0]
	return i
}

func (i *notNullCond) Do(db *gorm.DB) *gorm.DB {
	for k, v := range i.body {
		db = db.Where(fmt.Sprintf("`%s` IS NOT NULL", k), v)
	}
	return db
}

type orCond struct {
	body map[string]any
}

func newOrCond() *orCond {
	return &orCond{
		body: make(map[string]any),
	}
}

func (o *orCond) Build(key string, args ...any) Condition {
	o.body[key] = args[0]
	return o
}

func (o *orCond) Do(db *gorm.DB) *gorm.DB {
	for k, v := range o.body {
		db = db.Or(k, v)
	}
	return db
}

type orderByCond struct {
	body []string
}

func newOrderByCond() *orderByCond {
	return &orderByCond{
		body: make([]string, 0),
	}
}

func (o *orderByCond) Build(key string, args ...any) Condition {
	o.body = append(o.body, key)
	return o
}

func (o *orderByCond) Do(db *gorm.DB) *gorm.DB {
	for _, v := range o.body {
		db = db.Order(v)
	}
	return db
}

type limitCond struct {
	body int
}

func newLimitCond() *limitCond {
	return &limitCond{
		body: 0,
	}
}

func (l *limitCond) Build(key string, args ...any) Condition {
	l.body = args[0].(int)
	return l
}

func (l *limitCond) Do(db *gorm.DB) *gorm.DB {
	db = db.Limit(l.body)
	return db
}

type offsetCond struct {
	body int
}

func newOffsetCond() *offsetCond {
	return &offsetCond{
		body: 0,
	}
}

func (o *offsetCond) Build(key string, args ...any) Condition {
	o.body = args[0].(int)
	return o
}

func (o *offsetCond) Do(db *gorm.DB) *gorm.DB {
	db = db.Offset(o.body)
	return db
}

type customOrderCond struct {
	order        string
	defaultOrder string
}

func newCustomOrderCond() *customOrderCond {
	return &customOrderCond{}
}

func (c *customOrderCond) Build(_ string, args ...any) Condition {
	c.order = args[0].(string)
	c.defaultOrder = args[1].(string)
	return c
}

func (c *customOrderCond) Do(db *gorm.DB) *gorm.DB {
	if c.order == "" && c.defaultOrder == "" {
		return db
	}

	if c.order == "" {
		db.Order(c.defaultOrder)
		return db
	}

	sort := map[string]string{}
	err := sonic.UnmarshalString(c.order, &sort)
	if err != nil {
		db.Order(c.defaultOrder)
		return db
	}

	for k, v := range sort {
		o := "asc"
		if !strings.HasPrefix(strings.ToLower(v), "asc") {
			o = "desc"
		}
		db.Order(k + " " + o)
	}
	return db
}

type selectCond struct {
	fields []string
}

func newSelectCond() *selectCond {
	return &selectCond{
		fields: make([]string, 0),
	}
}

func (c *selectCond) Build(_ string, args ...any) Condition {
	for _, v := range args {
		s, ok := v.(string)
		if ok {
			c.fields = append(c.fields, s)
		}
	}
	return c
}

func (c *selectCond) Do(db *gorm.DB) *gorm.DB {
	db = db.Select(c.fields)
	return db
}

type groupCond struct {
	fields []string
}

func newGroupCond() *groupCond {
	return &groupCond{
		fields: make([]string, 0),
	}
}

func (c *groupCond) Build(_ string, args ...any) Condition {
	for _, v := range args {
		s, ok := v.(string)
		if ok {
			c.fields = append(c.fields, s)
		}
	}
	return c
}

func (c *groupCond) Do(db *gorm.DB) *gorm.DB {
	for _, v := range c.fields {
		db = db.Group(v)
	}
	return db
}

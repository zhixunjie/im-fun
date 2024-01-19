// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
)

func newGroup(db *gorm.DB, opts ...gen.DOOption) group {
	_group := group{}

	_group.groupDo.UseDB(db, opts...)
	_group.groupDo.UseModel(&model.Group{})

	tableName := _group.groupDo.TableName()
	_group.ALL = field.NewAsterisk(tableName)
	_group.ID = field.NewUint64(tableName, "id")
	_group.Name = field.NewString(tableName, "name")
	_group.CreatedAt = field.NewTime(tableName, "created_at")
	_group.UpdatedAt = field.NewTime(tableName, "updated_at")

	_group.fillFieldMap()

	return _group
}

// group 群组表
type group struct {
	groupDo

	ALL       field.Asterisk
	ID        field.Uint64 // 自增id,主键
	Name      field.String // 群组名称
	CreatedAt field.Time   // 创建时间
	UpdatedAt field.Time   // 更新时间

	fieldMap map[string]field.Expr
}

func (g group) Table(newTableName string) *group {
	g.groupDo.UseTable(newTableName)
	return g.updateTableName(newTableName)
}

func (g group) As(alias string) *group {
	g.groupDo.DO = *(g.groupDo.As(alias).(*gen.DO))
	return g.updateTableName(alias)
}

func (g *group) updateTableName(table string) *group {
	g.ALL = field.NewAsterisk(table)
	g.ID = field.NewUint64(table, "id")
	g.Name = field.NewString(table, "name")
	g.CreatedAt = field.NewTime(table, "created_at")
	g.UpdatedAt = field.NewTime(table, "updated_at")

	g.fillFieldMap()

	return g
}

func (g *group) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := g.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (g *group) fillFieldMap() {
	g.fieldMap = make(map[string]field.Expr, 4)
	g.fieldMap["id"] = g.ID
	g.fieldMap["name"] = g.Name
	g.fieldMap["created_at"] = g.CreatedAt
	g.fieldMap["updated_at"] = g.UpdatedAt
}

func (g group) clone(db *gorm.DB) group {
	g.groupDo.ReplaceConnPool(db.Statement.ConnPool)
	return g
}

func (g group) replaceDB(db *gorm.DB) group {
	g.groupDo.ReplaceDB(db)
	return g
}

type groupDo struct{ gen.DO }

type IGroupDo interface {
	gen.SubQuery
	Debug() IGroupDo
	WithContext(ctx context.Context) IGroupDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IGroupDo
	WriteDB() IGroupDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IGroupDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IGroupDo
	Not(conds ...gen.Condition) IGroupDo
	Or(conds ...gen.Condition) IGroupDo
	Select(conds ...field.Expr) IGroupDo
	Where(conds ...gen.Condition) IGroupDo
	Order(conds ...field.Expr) IGroupDo
	Distinct(cols ...field.Expr) IGroupDo
	Omit(cols ...field.Expr) IGroupDo
	Join(table schema.Tabler, on ...field.Expr) IGroupDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IGroupDo
	RightJoin(table schema.Tabler, on ...field.Expr) IGroupDo
	Group(cols ...field.Expr) IGroupDo
	Having(conds ...gen.Condition) IGroupDo
	Limit(limit int) IGroupDo
	Offset(offset int) IGroupDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IGroupDo
	Unscoped() IGroupDo
	Create(values ...*model.Group) error
	CreateInBatches(values []*model.Group, batchSize int) error
	Save(values ...*model.Group) error
	First() (*model.Group, error)
	Take() (*model.Group, error)
	Last() (*model.Group, error)
	Find() ([]*model.Group, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Group, err error)
	FindInBatches(result *[]*model.Group, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Group) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IGroupDo
	Assign(attrs ...field.AssignExpr) IGroupDo
	Joins(fields ...field.RelationField) IGroupDo
	Preload(fields ...field.RelationField) IGroupDo
	FirstOrInit() (*model.Group, error)
	FirstOrCreate() (*model.Group, error)
	FindByPage(offset int, limit int) (result []*model.Group, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IGroupDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	GetByID(id int64) (result *model.Group, err error)
}

// GetByID SELECT * FROM @@table WHERE id=@id
func (g groupDo) GetByID(id int64) (result *model.Group, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("SELECT * FROM group WHERE id=? ")

	var executeSQL *gorm.DB
	executeSQL = g.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (g groupDo) Debug() IGroupDo {
	return g.withDO(g.DO.Debug())
}

func (g groupDo) WithContext(ctx context.Context) IGroupDo {
	return g.withDO(g.DO.WithContext(ctx))
}

func (g groupDo) ReadDB() IGroupDo {
	return g.Clauses(dbresolver.Read)
}

func (g groupDo) WriteDB() IGroupDo {
	return g.Clauses(dbresolver.Write)
}

func (g groupDo) Session(config *gorm.Session) IGroupDo {
	return g.withDO(g.DO.Session(config))
}

func (g groupDo) Clauses(conds ...clause.Expression) IGroupDo {
	return g.withDO(g.DO.Clauses(conds...))
}

func (g groupDo) Returning(value interface{}, columns ...string) IGroupDo {
	return g.withDO(g.DO.Returning(value, columns...))
}

func (g groupDo) Not(conds ...gen.Condition) IGroupDo {
	return g.withDO(g.DO.Not(conds...))
}

func (g groupDo) Or(conds ...gen.Condition) IGroupDo {
	return g.withDO(g.DO.Or(conds...))
}

func (g groupDo) Select(conds ...field.Expr) IGroupDo {
	return g.withDO(g.DO.Select(conds...))
}

func (g groupDo) Where(conds ...gen.Condition) IGroupDo {
	return g.withDO(g.DO.Where(conds...))
}

func (g groupDo) Order(conds ...field.Expr) IGroupDo {
	return g.withDO(g.DO.Order(conds...))
}

func (g groupDo) Distinct(cols ...field.Expr) IGroupDo {
	return g.withDO(g.DO.Distinct(cols...))
}

func (g groupDo) Omit(cols ...field.Expr) IGroupDo {
	return g.withDO(g.DO.Omit(cols...))
}

func (g groupDo) Join(table schema.Tabler, on ...field.Expr) IGroupDo {
	return g.withDO(g.DO.Join(table, on...))
}

func (g groupDo) LeftJoin(table schema.Tabler, on ...field.Expr) IGroupDo {
	return g.withDO(g.DO.LeftJoin(table, on...))
}

func (g groupDo) RightJoin(table schema.Tabler, on ...field.Expr) IGroupDo {
	return g.withDO(g.DO.RightJoin(table, on...))
}

func (g groupDo) Group(cols ...field.Expr) IGroupDo {
	return g.withDO(g.DO.Group(cols...))
}

func (g groupDo) Having(conds ...gen.Condition) IGroupDo {
	return g.withDO(g.DO.Having(conds...))
}

func (g groupDo) Limit(limit int) IGroupDo {
	return g.withDO(g.DO.Limit(limit))
}

func (g groupDo) Offset(offset int) IGroupDo {
	return g.withDO(g.DO.Offset(offset))
}

func (g groupDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IGroupDo {
	return g.withDO(g.DO.Scopes(funcs...))
}

func (g groupDo) Unscoped() IGroupDo {
	return g.withDO(g.DO.Unscoped())
}

func (g groupDo) Create(values ...*model.Group) error {
	if len(values) == 0 {
		return nil
	}
	return g.DO.Create(values)
}

func (g groupDo) CreateInBatches(values []*model.Group, batchSize int) error {
	return g.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (g groupDo) Save(values ...*model.Group) error {
	if len(values) == 0 {
		return nil
	}
	return g.DO.Save(values)
}

func (g groupDo) First() (*model.Group, error) {
	if result, err := g.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Group), nil
	}
}

func (g groupDo) Take() (*model.Group, error) {
	if result, err := g.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Group), nil
	}
}

func (g groupDo) Last() (*model.Group, error) {
	if result, err := g.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Group), nil
	}
}

func (g groupDo) Find() ([]*model.Group, error) {
	result, err := g.DO.Find()
	return result.([]*model.Group), err
}

func (g groupDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Group, err error) {
	buf := make([]*model.Group, 0, batchSize)
	err = g.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (g groupDo) FindInBatches(result *[]*model.Group, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return g.DO.FindInBatches(result, batchSize, fc)
}

func (g groupDo) Attrs(attrs ...field.AssignExpr) IGroupDo {
	return g.withDO(g.DO.Attrs(attrs...))
}

func (g groupDo) Assign(attrs ...field.AssignExpr) IGroupDo {
	return g.withDO(g.DO.Assign(attrs...))
}

func (g groupDo) Joins(fields ...field.RelationField) IGroupDo {
	for _, _f := range fields {
		g = *g.withDO(g.DO.Joins(_f))
	}
	return &g
}

func (g groupDo) Preload(fields ...field.RelationField) IGroupDo {
	for _, _f := range fields {
		g = *g.withDO(g.DO.Preload(_f))
	}
	return &g
}

func (g groupDo) FirstOrInit() (*model.Group, error) {
	if result, err := g.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Group), nil
	}
}

func (g groupDo) FirstOrCreate() (*model.Group, error) {
	if result, err := g.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Group), nil
	}
}

func (g groupDo) FindByPage(offset int, limit int) (result []*model.Group, count int64, err error) {
	result, err = g.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = g.Offset(-1).Limit(-1).Count()
	return
}

func (g groupDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = g.Count()
	if err != nil {
		return
	}

	err = g.Offset(offset).Limit(limit).Scan(result)
	return
}

func (g groupDo) Scan(result interface{}) (err error) {
	return g.DO.Scan(result)
}

func (g groupDo) Delete(models ...*model.Group) (result gen.ResultInfo, err error) {
	return g.DO.Delete(models)
}

func (g *groupDo) withDO(do gen.Dao) *groupDo {
	g.DO = *do.(*gen.DO)
	return g
}

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

func newContact(db *gorm.DB, opts ...gen.DOOption) contact {
	_contact := contact{}

	_contact.contactDo.UseDB(db, opts...)
	_contact.contactDo.UseModel(&model.Contact{})

	tableName := _contact.contactDo.TableName()
	_contact.ALL = field.NewAsterisk(tableName)
	_contact.ID = field.NewUint64(tableName, "id")
	_contact.OwnerID = field.NewUint64(tableName, "owner_id")
	_contact.OwnerType = field.NewUint32(tableName, "owner_type")
	_contact.PeerID = field.NewUint64(tableName, "peer_id")
	_contact.PeerType = field.NewUint32(tableName, "peer_type")
	_contact.PeerAck = field.NewUint32(tableName, "peer_ack")
	_contact.LastMsgID = field.NewUint64(tableName, "last_msg_id")
	_contact.LastDelMsgID = field.NewUint64(tableName, "last_del_msg_id")
	_contact.VersionID = field.NewUint64(tableName, "version_id")
	_contact.SortKey = field.NewUint64(tableName, "sort_key")
	_contact.Status = field.NewUint32(tableName, "status")
	_contact.Labels = field.NewString(tableName, "labels")
	_contact.CreatedAt = field.NewTime(tableName, "created_at")
	_contact.UpdatedAt = field.NewTime(tableName, "updated_at")

	_contact.fillFieldMap()

	return _contact
}

// contact 会话表（通信双方各有一行记录）
type contact struct {
	contactDo

	ALL          field.Asterisk
	ID           field.Uint64 // 自增id,主键
	OwnerID      field.Uint64 // 会话拥有者
	OwnerType    field.Uint32 // 用户类型（owner_id）
	PeerID       field.Uint64 // 联系人（对方用户）
	PeerType     field.Uint32 // 用户类型（peer_id）
	PeerAck      field.Uint32 // peer_id是否给owner发过消息，0：未发过，1：发过
	LastMsgID    field.Uint64 // 聊天记录中，最新一条发送的私信id
	LastDelMsgID field.Uint64 // 聊天记录中，最后一次删除联系人时的私信id
	VersionID    field.Uint64 // 版本id（用于拉取会话框）
	SortKey      field.Uint64 // 会话展示顺序（按顺序展示会话）可修改顺序，如：置顶操作
	Status       field.Uint32 // 联系人状态，0：正常，1：被删除
	Labels       field.String // 会话标签，json字符串
	CreatedAt    field.Time   // 创建时间
	UpdatedAt    field.Time   // 更新时间

	fieldMap map[string]field.Expr
}

func (c contact) Table(newTableName string) *contact {
	c.contactDo.UseTable(newTableName)
	return c.updateTableName(newTableName)
}

func (c contact) As(alias string) *contact {
	c.contactDo.DO = *(c.contactDo.As(alias).(*gen.DO))
	return c.updateTableName(alias)
}

func (c *contact) updateTableName(table string) *contact {
	c.ALL = field.NewAsterisk(table)
	c.ID = field.NewUint64(table, "id")
	c.OwnerID = field.NewUint64(table, "owner_id")
	c.OwnerType = field.NewUint32(table, "owner_type")
	c.PeerID = field.NewUint64(table, "peer_id")
	c.PeerType = field.NewUint32(table, "peer_type")
	c.PeerAck = field.NewUint32(table, "peer_ack")
	c.LastMsgID = field.NewUint64(table, "last_msg_id")
	c.LastDelMsgID = field.NewUint64(table, "last_del_msg_id")
	c.VersionID = field.NewUint64(table, "version_id")
	c.SortKey = field.NewUint64(table, "sort_key")
	c.Status = field.NewUint32(table, "status")
	c.Labels = field.NewString(table, "labels")
	c.CreatedAt = field.NewTime(table, "created_at")
	c.UpdatedAt = field.NewTime(table, "updated_at")

	c.fillFieldMap()

	return c
}

func (c *contact) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := c.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (c *contact) fillFieldMap() {
	c.fieldMap = make(map[string]field.Expr, 14)
	c.fieldMap["id"] = c.ID
	c.fieldMap["owner_id"] = c.OwnerID
	c.fieldMap["owner_type"] = c.OwnerType
	c.fieldMap["peer_id"] = c.PeerID
	c.fieldMap["peer_type"] = c.PeerType
	c.fieldMap["peer_ack"] = c.PeerAck
	c.fieldMap["last_msg_id"] = c.LastMsgID
	c.fieldMap["last_del_msg_id"] = c.LastDelMsgID
	c.fieldMap["version_id"] = c.VersionID
	c.fieldMap["sort_key"] = c.SortKey
	c.fieldMap["status"] = c.Status
	c.fieldMap["labels"] = c.Labels
	c.fieldMap["created_at"] = c.CreatedAt
	c.fieldMap["updated_at"] = c.UpdatedAt
}

func (c contact) clone(db *gorm.DB) contact {
	c.contactDo.ReplaceConnPool(db.Statement.ConnPool)
	return c
}

func (c contact) replaceDB(db *gorm.DB) contact {
	c.contactDo.ReplaceDB(db)
	return c
}

type contactDo struct{ gen.DO }

type IContactDo interface {
	gen.SubQuery
	Debug() IContactDo
	WithContext(ctx context.Context) IContactDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IContactDo
	WriteDB() IContactDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IContactDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IContactDo
	Not(conds ...gen.Condition) IContactDo
	Or(conds ...gen.Condition) IContactDo
	Select(conds ...field.Expr) IContactDo
	Where(conds ...gen.Condition) IContactDo
	Order(conds ...field.Expr) IContactDo
	Distinct(cols ...field.Expr) IContactDo
	Omit(cols ...field.Expr) IContactDo
	Join(table schema.Tabler, on ...field.Expr) IContactDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IContactDo
	RightJoin(table schema.Tabler, on ...field.Expr) IContactDo
	Group(cols ...field.Expr) IContactDo
	Having(conds ...gen.Condition) IContactDo
	Limit(limit int) IContactDo
	Offset(offset int) IContactDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IContactDo
	Unscoped() IContactDo
	Create(values ...*model.Contact) error
	CreateInBatches(values []*model.Contact, batchSize int) error
	Save(values ...*model.Contact) error
	First() (*model.Contact, error)
	Take() (*model.Contact, error)
	Last() (*model.Contact, error)
	Find() ([]*model.Contact, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Contact, err error)
	FindInBatches(result *[]*model.Contact, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Contact) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IContactDo
	Assign(attrs ...field.AssignExpr) IContactDo
	Joins(fields ...field.RelationField) IContactDo
	Preload(fields ...field.RelationField) IContactDo
	FirstOrInit() (*model.Contact, error)
	FirstOrCreate() (*model.Contact, error)
	FindByPage(offset int, limit int) (result []*model.Contact, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IContactDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	GetByID(id int64) (result *model.Contact, err error)
}

// GetByID SELECT * FROM @@table WHERE id=@id
func (c contactDo) GetByID(id int64) (result *model.Contact, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("SELECT * FROM contact WHERE id=? ")

	var executeSQL *gorm.DB
	executeSQL = c.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (c contactDo) Debug() IContactDo {
	return c.withDO(c.DO.Debug())
}

func (c contactDo) WithContext(ctx context.Context) IContactDo {
	return c.withDO(c.DO.WithContext(ctx))
}

func (c contactDo) ReadDB() IContactDo {
	return c.Clauses(dbresolver.Read)
}

func (c contactDo) WriteDB() IContactDo {
	return c.Clauses(dbresolver.Write)
}

func (c contactDo) Session(config *gorm.Session) IContactDo {
	return c.withDO(c.DO.Session(config))
}

func (c contactDo) Clauses(conds ...clause.Expression) IContactDo {
	return c.withDO(c.DO.Clauses(conds...))
}

func (c contactDo) Returning(value interface{}, columns ...string) IContactDo {
	return c.withDO(c.DO.Returning(value, columns...))
}

func (c contactDo) Not(conds ...gen.Condition) IContactDo {
	return c.withDO(c.DO.Not(conds...))
}

func (c contactDo) Or(conds ...gen.Condition) IContactDo {
	return c.withDO(c.DO.Or(conds...))
}

func (c contactDo) Select(conds ...field.Expr) IContactDo {
	return c.withDO(c.DO.Select(conds...))
}

func (c contactDo) Where(conds ...gen.Condition) IContactDo {
	return c.withDO(c.DO.Where(conds...))
}

func (c contactDo) Order(conds ...field.Expr) IContactDo {
	return c.withDO(c.DO.Order(conds...))
}

func (c contactDo) Distinct(cols ...field.Expr) IContactDo {
	return c.withDO(c.DO.Distinct(cols...))
}

func (c contactDo) Omit(cols ...field.Expr) IContactDo {
	return c.withDO(c.DO.Omit(cols...))
}

func (c contactDo) Join(table schema.Tabler, on ...field.Expr) IContactDo {
	return c.withDO(c.DO.Join(table, on...))
}

func (c contactDo) LeftJoin(table schema.Tabler, on ...field.Expr) IContactDo {
	return c.withDO(c.DO.LeftJoin(table, on...))
}

func (c contactDo) RightJoin(table schema.Tabler, on ...field.Expr) IContactDo {
	return c.withDO(c.DO.RightJoin(table, on...))
}

func (c contactDo) Group(cols ...field.Expr) IContactDo {
	return c.withDO(c.DO.Group(cols...))
}

func (c contactDo) Having(conds ...gen.Condition) IContactDo {
	return c.withDO(c.DO.Having(conds...))
}

func (c contactDo) Limit(limit int) IContactDo {
	return c.withDO(c.DO.Limit(limit))
}

func (c contactDo) Offset(offset int) IContactDo {
	return c.withDO(c.DO.Offset(offset))
}

func (c contactDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IContactDo {
	return c.withDO(c.DO.Scopes(funcs...))
}

func (c contactDo) Unscoped() IContactDo {
	return c.withDO(c.DO.Unscoped())
}

func (c contactDo) Create(values ...*model.Contact) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Create(values)
}

func (c contactDo) CreateInBatches(values []*model.Contact, batchSize int) error {
	return c.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (c contactDo) Save(values ...*model.Contact) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Save(values)
}

func (c contactDo) First() (*model.Contact, error) {
	if result, err := c.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Contact), nil
	}
}

func (c contactDo) Take() (*model.Contact, error) {
	if result, err := c.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Contact), nil
	}
}

func (c contactDo) Last() (*model.Contact, error) {
	if result, err := c.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Contact), nil
	}
}

func (c contactDo) Find() ([]*model.Contact, error) {
	result, err := c.DO.Find()
	return result.([]*model.Contact), err
}

func (c contactDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Contact, err error) {
	buf := make([]*model.Contact, 0, batchSize)
	err = c.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (c contactDo) FindInBatches(result *[]*model.Contact, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return c.DO.FindInBatches(result, batchSize, fc)
}

func (c contactDo) Attrs(attrs ...field.AssignExpr) IContactDo {
	return c.withDO(c.DO.Attrs(attrs...))
}

func (c contactDo) Assign(attrs ...field.AssignExpr) IContactDo {
	return c.withDO(c.DO.Assign(attrs...))
}

func (c contactDo) Joins(fields ...field.RelationField) IContactDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Joins(_f))
	}
	return &c
}

func (c contactDo) Preload(fields ...field.RelationField) IContactDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Preload(_f))
	}
	return &c
}

func (c contactDo) FirstOrInit() (*model.Contact, error) {
	if result, err := c.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Contact), nil
	}
}

func (c contactDo) FirstOrCreate() (*model.Contact, error) {
	if result, err := c.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Contact), nil
	}
}

func (c contactDo) FindByPage(offset int, limit int) (result []*model.Contact, count int64, err error) {
	result, err = c.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = c.Offset(-1).Limit(-1).Count()
	return
}

func (c contactDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = c.Count()
	if err != nil {
		return
	}

	err = c.Offset(offset).Limit(limit).Scan(result)
	return
}

func (c contactDo) Scan(result interface{}) (err error) {
	return c.DO.Scan(result)
}

func (c contactDo) Delete(models ...*model.Contact) (result gen.ResultInfo, err error) {
	return c.DO.Delete(models)
}

func (c *contactDo) withDO(do gen.Dao) *contactDo {
	c.DO = *do.(*gen.DO)
	return c
}

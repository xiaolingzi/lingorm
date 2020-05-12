package drivers

import (
	"github.com/xiaolingzi/lingorm/internal/common"
)

// IQuery the interface of Query
type IQuery interface {
	Table(table interface{}) ITableQuery
	Find(table interface{}, where interface{}, orderBy interface{}, slicePtr ...interface{}) (interface{}, error)
	FindTop(table interface{}, where interface{}, orderBy interface{}, top int, slicePtr ...interface{}) (interface{}, error)
	First(table interface{}, where interface{}, orderBy interface{}, structPtr ...interface{}) (interface{}, error)
	FindPage(table interface{}, where interface{}, orderBy interface{}, pageIndex int, pageSize int, slicePtr ...interface{}) (common.PageResult, error)
	Insert(model interface{}) (int, error)
	BatchInsert(modelList []interface{}) (int, error)
	Update(model interface{}) (int, error)
	BatchUpdate(modelList []interface{}) (int, error)
	UpdateBy(table interface{}, params []interface{}, where IWhere) (int, error)
	Delete(model interface{}) (int, error)
	DeleteBy(table interface{}, where IWhere) (int, error)

	QueryBuilder() IQueryBuilder
	NativeQuery() INativeQuery
	CreateWhere() IWhere
	CreateOderBy() IOrderBy
	CreateGroupBy() IGroupBy

	Begin() error
	Commit() error
	Rollback() error
}

// IQueryBuilder the interface of QueryBuilder
type IQueryBuilder interface {
	Select(args ...interface{}) IQueryBuilder
	From(table interface{}) IQueryBuilder
	LeftJoin(table interface{}, whereOrConditions ...interface{}) IQueryBuilder
	RightJoin(table interface{}, whereOrConditions ...interface{}) IQueryBuilder
	InnerJoin(table interface{}, whereOrConditions ...interface{}) IQueryBuilder
	Where(whereOrConditions ...interface{}) IQueryBuilder
	GroupBy(args ...interface{}) IQueryBuilder
	OrderBy(args ...interface{}) IQueryBuilder
	Limit(count int) IQueryBuilder
	Find(slicePtr ...interface{}) (interface{}, error)
	First(structPtr ...interface{}) (interface{}, error)
	FindPage(pageIndex int, pageSize int, slicePtr ...interface{}) (common.PageResult, error)
	FindCount() (int, error)

	CurrentSQL() string
}

// ITableQuery the interface of TableQuery
type ITableQuery interface {
	Table(databaseConfigKey string, table interface{}, transactionKey string) ITableQuery
	Select(args ...interface{}) ITableQuery
	Where(args ...interface{}) ITableQuery
	OrderBy(args ...interface{}) ITableQuery
	GroupBy(args ...interface{}) ITableQuery
	Limit(count int) ITableQuery
	Find(slicePtr ...interface{}) (interface{}, error)
	FindPage(pageIndex int, pageSize int, slicePtr ...interface{}) (common.PageResult, error)
	First(structPtr ...interface{}) (interface{}, error)
	FindCount() (int, error)
	CurrentSQL() string
}

// INativeQuery the interface of NativeQuery
type INativeQuery interface {
	Execute(sql string, params map[string]interface{}) (int, int, error)
	Find(sql string, params map[string]interface{}, slicePtr ...interface{}) (interface{}, error)
	FindPage(pageIndex int, pageSize int, sql string, params map[string]interface{}, slicePtr ...interface{}) (common.PageResult, error)
	FindCount(sql string, params map[string]interface{}) (int, error)
	First(sql string, params map[string]interface{}, slicePtr ...interface{}) (interface{}, error)
}

// IWhere the interface of Where
type IWhere interface {
	Or(args ...interface{}) IWhere
	GetOr(args ...interface{}) string
	And(args ...interface{}) IWhere
	GetAnd(args ...interface{}) string
	CurrentSQL() string
}

// IGroupBy the interface of GroupBy
type IGroupBy interface {
	By(args ...interface{}) IGroupBy
}

// IOrderBy the interface of OrderBy
type IOrderBy interface {
	By(args ...interface{}) IOrderBy
}

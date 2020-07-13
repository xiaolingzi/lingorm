package lingorm

import (
	"strings"

	"github.com/xiaolingzi/lingorm/internal/config"
	"github.com/xiaolingzi/lingorm/internal/drivers"
	"github.com/xiaolingzi/lingorm/internal/drivers/mysql"
)

// IQuery query interface
type IQuery interface {
	drivers.IQuery
}

// IQueryBuilder the interface of QueryBuilder
type IQueryBuilder interface {
	drivers.IQueryBuilder
}

// ITableQuery the interface of TableQuery
type ITableQuery interface {
	drivers.ITableQuery
}

// INativeQuery query interface
type INativeQuery interface {
	drivers.INativeQuery
}

// IWhere the interface of Where
type IWhere interface {
	drivers.IWhere
}

// IGroupBy the interface of GroupBy
type IGroupBy interface {
	drivers.IGroupBy
}

// IOrderBy the interface of OrderBy
type IOrderBy interface {
	drivers.IOrderBy
}

// DB the instance of Query
func DB(databaseConfigKey string) IQuery {
	driver := strings.ToLower(config.GetDatabaseDriver(databaseConfigKey))
	if driver == "msyql" {
		return mysql.NewQuery(databaseConfigKey)
	}
	return mysql.NewQuery(databaseConfigKey)
}

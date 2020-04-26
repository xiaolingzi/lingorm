package mysql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xiaolingzi/lingorm/internal/common"
	"github.com/xiaolingzi/lingorm/internal/drivers"
	"github.com/xiaolingzi/lingorm/model"
)

// QueryBuilder struct of the builder
type QueryBuilder struct {
	DatabaseConfigKey string
	TransactionKey    string
	SQL               string
	Params            map[string]interface{}
	SelectSQL         string
	FromSQL           string
	JoinSQL           string
	WhereSQL          string
	GroupBySQL        string
	OrderBySQL        string
	LimitCount        int
}

// NewQueryBuilder return a new QueryBuilder instance
func NewQueryBuilder(databaseConfigKey string, transactionKey string) *QueryBuilder {
	var builder QueryBuilder
	builder.DatabaseConfigKey = databaseConfigKey
	builder.TransactionKey = transactionKey
	return &builder
}

// Select columns selected
func (b *QueryBuilder) Select(args ...interface{}) drivers.IQueryBuilder {
	b.SelectSQL = NewColumn().GetSelectColumns(args...)
	return b
}

// From the table selected from
func (b *QueryBuilder) From(table interface{}) drivers.IQueryBuilder {
	tableName, _ := model.NewMapping().GetSQLTableName(table)
	b.FromSQL = tableName
	return b
}

// LeftJoin left join
func (b *QueryBuilder) LeftJoin(table interface{}, whereOrConditions ...interface{}) drivers.IQueryBuilder {
	where := (NewWhere().And(b.WhereSQL).And(whereOrConditions...)).(*Where)
	b.Params = b.mergeParams(b.Params, where.Params)

	tableName, _ := model.NewMapping().GetSQLTableName(table)
	b.JoinSQL += fmt.Sprintf(" LEFT JOIN %s ON %s", tableName, where.SQL)

	return b
}

// RightJoin right join
func (b *QueryBuilder) RightJoin(table interface{}, whereOrConditions ...interface{}) drivers.IQueryBuilder {
	where := (NewWhere().And(b.WhereSQL).And(whereOrConditions...)).(*Where)
	b.Params = b.mergeParams(b.Params, where.Params)

	tableName, _ := model.NewMapping().GetSQLTableName(table)
	b.JoinSQL += fmt.Sprintf(" RIGHT JOIN %s ON %s", tableName, where.SQL)

	return b
}

// InnerJoin inner join
func (b *QueryBuilder) InnerJoin(table interface{}, whereOrConditions ...interface{}) drivers.IQueryBuilder {
	where := (NewWhere().And(b.WhereSQL).And(whereOrConditions...)).(*Where)
	b.Params = b.mergeParams(b.Params, where.Params)

	tableName, _ := model.NewMapping().GetSQLTableName(table)
	b.JoinSQL += fmt.Sprintf(" INNER JOIN %s ON %s", tableName, where.SQL)

	return b
}

// Where where
func (b *QueryBuilder) Where(whereOrConditions ...interface{}) drivers.IQueryBuilder {
	where := (NewWhere().And(b.WhereSQL).And(whereOrConditions...)).(*Where)
	b.WhereSQL = where.SQL
	b.Params = b.mergeParams(b.Params, where.Params)
	return b
}

// GroupBy group by
func (b *QueryBuilder) GroupBy(args ...interface{}) drivers.IQueryBuilder {
	group := (NewGroupBy().By(b.GroupBySQL).By(args...)).(*GroupBy)
	b.GroupBySQL = group.SQL
	return b
}

// OrderBy order by
func (b *QueryBuilder) OrderBy(args ...interface{}) drivers.IQueryBuilder {
	orderStr := (NewOrderBy().By(args...)).(*OrderBy).SQL
	orderStr = b.OrderBySQL + "," + orderStr
	b.OrderBySQL = strings.Trim(orderStr, ",")
	return b
}

// Limit return the limited top rows
func (b *QueryBuilder) Limit(count int) drivers.IQueryBuilder {
	b.LimitCount = count
	return b
}

// Find return all the rows that meet query criteria
func (b *QueryBuilder) Find(args ...interface{}) (interface{}, error) {
	sql := b.getSelectSQL()
	return NewNativeQuery(b.DatabaseConfigKey, b.TransactionKey).Find(sql, b.Params, args...)
}

// First return the first row that meet query criteria
func (b *QueryBuilder) First(args ...interface{}) (interface{}, error) {
	b.LimitCount = 0
	sql := b.getSelectSQL()
	return NewNativeQuery(b.DatabaseConfigKey, b.TransactionKey).First(sql, b.Params, args...)
}

// FindPage return the page result
func (b *QueryBuilder) FindPage(pageIndex int, pageSize int, args ...interface{}) (common.PageResult, error) {
	sql := b.getSelectSQL()
	return NewNativeQuery(b.DatabaseConfigKey, b.TransactionKey).FindPage(pageIndex, pageSize, sql, b.Params, args...)
}

// FindCount return the number of rows that meet query criteria
func (b *QueryBuilder) FindCount() (int, error) {
	sql := b.getSelectSQL()
	return NewNativeQuery(b.DatabaseConfigKey, b.TransactionKey).FindCount(sql, b.Params)
}

// CurrentSQL return the current sql
func (b *QueryBuilder) CurrentSQL() string {
	return b.getSelectSQL()
}

func (b *QueryBuilder) getSelectSQL() string {
	if b.SelectSQL == "" {
		b.SelectSQL = "*"
	}
	sql := fmt.Sprintf("SELECT %s FROM %s", b.SelectSQL, b.FromSQL)
	if b.JoinSQL != "" {
		sql += " " + strings.TrimSpace(b.JoinSQL)
	}
	if b.WhereSQL != "" {
		sql += " WHERE " + b.WhereSQL
	}
	if b.GroupBySQL != "" {
		sql += " GROUP BY " + b.GroupBySQL
	}
	if b.OrderBySQL != "" {
		sql += " ORDER BY " + b.OrderBySQL
	}
	if b.LimitCount > 0 {
		sql += " LIMIT " + strconv.Itoa(b.LimitCount)
	}
	return sql
}

func (b *QueryBuilder) mergeParams(params1 map[string]interface{}, params2 map[string]interface{}) map[string]interface{} {
	if params1 == nil {
		return params2
	}
	if params2 == nil {
		return params1
	}
	for key, value := range params2 {
		params1[key] = value
	}
	return params1
}

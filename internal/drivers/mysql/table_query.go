package mysql

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/xiaolingzi/lingorm/internal/common"
	"github.com/xiaolingzi/lingorm/internal/drivers"
	"github.com/xiaolingzi/lingorm/model"
)

// TableQuery struct
type TableQuery struct {
	DatabaseConfigKey string
	TransactionKey    string
	SQL               string
	TableSQL          string
	SelectSQL         string
	WhereSQL          string
	GroupBySQL        string
	OrderBySQL        string
	LimitCount        int
	Params            map[string]interface{}
	Model             interface{}
}

// Table the table selected from
func (s *TableQuery) Table(databaseConfigKey string, table interface{}, transactionKey string) drivers.ITableQuery {
	s.DatabaseConfigKey = databaseConfigKey
	s.TransactionKey = transactionKey
	tableName, _ := model.NewMapping().GetSQLTableName(table)
	s.TableSQL = tableName
	modelType := model.NewMapping().GetModelType(table)
	s.Model = reflect.New(modelType).Elem().Interface()

	return s
}

// Select columns for select
func (s *TableQuery) Select(args ...interface{}) drivers.ITableQuery {
	s.SelectSQL = NewColumn().GetSelectColumns(args...)
	return s
}

// Where where
func (s *TableQuery) Where(args ...interface{}) drivers.ITableQuery {
	where := (NewWhere().And(s.WhereSQL).And(args...)).(*Where)
	s.WhereSQL = where.SQL
	s.Params = where.Params
	return s
}

// GroupBy group by
func (s *TableQuery) GroupBy(args ...interface{}) drivers.ITableQuery {
	group := (NewGroupBy().By(s.GroupBySQL).By(args...)).(*GroupBy)
	s.GroupBySQL = group.SQL
	return s
}

// OrderBy order by
func (s *TableQuery) OrderBy(args ...interface{}) drivers.ITableQuery {
	orderStr := (NewOrderBy().By(args...)).(*OrderBy).SQL
	orderStr = s.OrderBySQL + "," + orderStr
	orderStr = strings.Trim(orderStr, ",")
	s.OrderBySQL = orderStr
	return s
}

// Limit the number of top rows
func (s *TableQuery) Limit(count int) drivers.ITableQuery {
	s.LimitCount = count
	return s
}

// Find return all the rows that meet query criteria
func (s *TableQuery) Find(slicePtr ...interface{}) (interface{}, error) {
	sql := s.getSelectSQL()
	if len(slicePtr) > 0 {
		return NewNativeQuery(s.DatabaseConfigKey, s.TransactionKey).Find(sql, s.Params, slicePtr...)
	}
	return NewNativeQuery(s.DatabaseConfigKey, s.TransactionKey).Find(sql, s.Params, s.Model)
}

// FindPage return the page result
func (s *TableQuery) FindPage(pageIndex int, pageSize int, slicePtr ...interface{}) (common.PageResult, error) {
	sql := s.getSelectSQL()
	if len(slicePtr) > 0 {
		return NewNativeQuery(s.DatabaseConfigKey, s.TransactionKey).FindPage(pageIndex, pageSize, sql, s.Params, slicePtr...)
	}
	return NewNativeQuery(s.DatabaseConfigKey, s.TransactionKey).FindPage(pageIndex, pageSize, sql, s.Params, s.Model)
}

// First return the first row that meet query criteria
func (s *TableQuery) First(structPtr ...interface{}) (interface{}, error) {
	s.LimitCount = 0
	sql := s.getSelectSQL()
	if len(structPtr) > 0 {
		return NewNativeQuery(s.DatabaseConfigKey, s.TransactionKey).First(sql, s.Params, structPtr...)
	}
	return NewNativeQuery(s.DatabaseConfigKey, s.TransactionKey).First(sql, s.Params, s.Model)
}

// FindCount return the number of rows that meet query criteria
func (s *TableQuery) FindCount() (int, error) {
	sql := s.getSelectSQL()
	return NewNativeQuery(s.DatabaseConfigKey, s.TransactionKey).FindCount(sql, s.Params)
}

// CurrentSQL return the current sql
func (s *TableQuery) CurrentSQL() string {
	return s.getSelectSQL()
}

func (s *TableQuery) getSelectSQL() string {
	if s.SelectSQL == "" {
		s.SelectSQL = "*"
	}
	sql := "SELECT " + s.SelectSQL + " FROM " + s.TableSQL
	if s.WhereSQL != "" {
		sql += " WHERE " + s.WhereSQL
	}
	if s.GroupBySQL != "" {
		sql += " GROUP BY " + s.GroupBySQL
	}
	if s.OrderBySQL != "" {
		sql += " ORDER BY " + s.OrderBySQL
	}
	if s.LimitCount > 0 {
		sql += " LIMIT " + strconv.Itoa(int(s.LimitCount))
	}
	s.SQL = sql
	return sql
}

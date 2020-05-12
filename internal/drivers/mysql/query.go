package mysql

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/xiaolingzi/lingorm/internal/common"
	"github.com/xiaolingzi/lingorm/internal/drivers"
	"github.com/xiaolingzi/lingorm/model"
)

// Query struct
type Query struct {
	DatabaseConfigKey string
	TransactionKey    string
}

// NewQuery return the instance of Query
func NewQuery(databaseConfigKey string) *Query {
	var query Query
	query.DatabaseConfigKey = databaseConfigKey
	return &query
}

// Table the table selected from
func (q *Query) Table(table interface{}) drivers.ITableQuery {
	var mysqlSelect TableQuery
	return mysqlSelect.Table(q.DatabaseConfigKey, table, q.TransactionKey)
}

// Find return all the rows that meet query criteria
func (q *Query) Find(table interface{}, where interface{}, orderBy interface{}, slicePtr ...interface{}) (interface{}, error) {
	sql, params := q.getSelectSQL(table, where, orderBy)
	if len(slicePtr) > 0 {
		return NewNativeQuery(q.DatabaseConfigKey, q.TransactionKey).Find(sql, params, slicePtr...)
	} else {
		modelType := model.NewMapping().GetModelType(table)
		model := reflect.New(modelType).Elem().Interface()
		return NewNativeQuery(q.DatabaseConfigKey, q.TransactionKey).Find(sql, params, model)
	}

}

// FindTop return the top rows that meet query criteria
func (q *Query) FindTop(table interface{}, where interface{}, orderBy interface{}, top int, slicePtr ...interface{}) (interface{}, error) {
	sql, params := q.getSelectSQL(table, where, orderBy)
	if top > 0 {
		sql += " LIMIT " + strconv.Itoa(top)
	}
	if len(slicePtr) > 0 {
		return NewNativeQuery(q.DatabaseConfigKey, q.TransactionKey).Find(sql, params, slicePtr...)
	} else {
		modelType := model.NewMapping().GetModelType(table)
		model := reflect.New(modelType).Elem().Interface()
		return NewNativeQuery(q.DatabaseConfigKey, q.TransactionKey).Find(sql, params, model)
	}
}

// First return the first row that meet query criteria
func (q Query) First(table interface{}, where interface{}, orderBy interface{}, structPtr ...interface{}) (interface{}, error) {
	result, err := q.FindTop(table, where, orderBy, 1, structPtr...)
	if err != nil && result != nil {
		return result.([]interface{})[0], nil
	}
	return nil, err
}

// FindPage return the page result
func (q *Query) FindPage(table interface{}, where interface{}, orderBy interface{}, pageIndex int, pageSize int, slicePtr ...interface{}) (common.PageResult, error) {
	sql, params := q.getSelectSQL(table, where, orderBy)
	if len(slicePtr) > 0 {
		return NewNativeQuery(q.DatabaseConfigKey, q.TransactionKey).FindPage(pageIndex, pageSize, sql, params, slicePtr...)
	} else {
		modelType := model.NewMapping().GetModelType(table)
		model := reflect.New(modelType).Elem().Interface()
		return NewNativeQuery(q.DatabaseConfigKey, q.TransactionKey).FindPage(pageIndex, pageSize, sql, params, model)
	}
}

// Insert insert data
func (q *Query) Insert(modelInstance interface{}) (int, error) {
	mapData, tableName := model.NewMapping().GetModelData(modelInstance)
	columns := ""
	values := ""
	params := make(map[string]interface{})
	for _, field := range mapData {
		fieldType := reflect.TypeOf(field.Value).String()
		if field.AutoIncrement || field.Value == nil || (fieldType == "string" && field.Value.(string) == "") || (fieldType == "time.Time" && field.Value.(time.Time).IsZero()) {
			continue
		}
		columns += field.ColumnName + ","
		values += ":" + field.FieldName + ","
		params[field.FieldName] = field.Value
	}
	columns = strings.Trim(columns, ",")
	values = strings.Trim(values, ",")
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, values)
	_, id, err := NewNativeQuery(q.DatabaseConfigKey, q.TransactionKey).Execute(sql, params)
	return id, err
}

// BatchInsert batch insert data
func (q *Query) BatchInsert(modelList []interface{}) (int, error) {
	if len(modelList) == 0 {
		return 0, nil
	} else if len(modelList) == 1 {
		return q.Insert(modelList[0])
	}

	mapData, tableName := model.NewMapping().GetModelData(modelList[0])
	columns := ""
	var fieldArr []string
	params := make(map[string]interface{})
	for _, field := range mapData {
		fieldType := reflect.TypeOf(field.Value).String()
		if field.AutoIncrement || field.Value == nil || (fieldType == "string" && field.Value.(string) == "") || (fieldType == "time.Time" && field.Value.(time.Time).IsZero()) {
			continue
		}
		columns += field.ColumnName + ","
		fieldArr = append(fieldArr, field.FieldName)
	}
	columns = strings.Trim(columns, ",")

	values := ""
	for i := 0; i < len(modelList); i++ {
		tempValueStr := ""
		modelType := reflect.ValueOf(modelList[i])
		for j := 0; j < len(fieldArr); j++ {
			key := fieldArr[j] + strconv.Itoa(i)
			tempValueStr += ":" + key + ","
			params[key] = modelType.FieldByName(fieldArr[j]).Interface()
		}
		tempValueStr = strings.Trim(tempValueStr, ",")
		values += fmt.Sprintf("(%s),", tempValueStr)
	}
	values = strings.Trim(values, ",")

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableName, columns, values)
	count, _, err := NewNativeQuery(q.DatabaseConfigKey, q.TransactionKey).Execute(sql, params)
	return count, err
}

// Update update data
func (q *Query) Update(modelInstance interface{}) (affected int, err error) {
	common.NewError().Defer(&err)
	mapData, tableName := model.NewMapping().GetModelData(modelInstance)
	whereStr := ""
	setStr := ""
	params := make(map[string]interface{})
	for _, field := range mapData {
		fieldType := reflect.TypeOf(field.Value).String()
		if field.IsPrimaryKey {
			if whereStr == "" {
				whereStr += field.ColumnName + "=:" + field.FieldName
			} else {
				whereStr += " AND " + field.ColumnName + "=:" + field.FieldName
			}
			params[field.FieldName] = field.Value
		}
		if field.AutoIncrement {
			continue
		}
		if field.Value == nil || (fieldType == "string" && field.Value.(string) == "") || (fieldType == "time.Time" && field.Value.(time.Time).IsZero()) {
			setStr += field.ColumnName + "=NULL,"
		} else {
			setStr += field.ColumnName + "=:" + field.FieldName + ","
			params[field.FieldName] = field.Value
		}
	}
	if whereStr == "" {
		common.NewError().Throw("Update method require at least one primary key")
	}
	setStr = strings.Trim(setStr, ",")
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, setStr, whereStr)
	affected, _, err = NewNativeQuery(q.DatabaseConfigKey, q.TransactionKey).Execute(sql, params)
	return affected, err
}

// BatchUpdate batch update data
func (q *Query) BatchUpdate(modelList []interface{}) (affected int, err error) {
	common.NewError().Defer(&err)
	if len(modelList) == 0 {
		return 0, nil
	} else if len(modelList) == 1 {
		return q.Update(modelList[0])
	}

	mapData, tableName := model.NewMapping().GetModelData(modelList[0])

	var primaryField model.Field
	setStr := ""
	inStr := ""

	primaryKeyCount := 0
	for _, field := range mapData {
		if field.IsPrimaryKey {
			primaryField = field
			primaryKeyCount++
		}
	}
	if primaryKeyCount > 1 {
		common.NewError().Throw("Batch update method only supports one primary key")
	} else if primaryKeyCount <= 0 {
		common.NewError().Throw("Update method require at least one primary key")
	}

	params := make(map[string]interface{})
	fieldSetArr := make(map[string]string)
	for i := 0; i < len(modelList); i++ {
		modelType := reflect.ValueOf(modelList[i])
		tempPrimaryKey := primaryField.FieldName + strconv.Itoa(i)
		inStr += ":" + tempPrimaryKey + ","
		if modelType.FieldByName(primaryField.FieldName) == (reflect.Value{}) {
			common.NewError().Throw("Invalid primary key value")
		}
		params[tempPrimaryKey] = modelType.FieldByName(primaryField.FieldName).Interface()
		for _, field := range mapData {
			if field.AutoIncrement {
				continue
			}
			fieldType := reflect.TypeOf(field.Value).String()
			if field.Value == nil || (fieldType == "string" && field.Value.(string) == "") || (fieldType == "time.Time" && field.Value.(time.Time).IsZero()) {
				fieldSetArr[field.ColumnName] += fmt.Sprintf(" WHEN :%s THEN NULL", tempPrimaryKey)
			} else {
				tempFieldName := field.FieldName + strconv.Itoa(i)
				fieldSetArr[field.ColumnName] += fmt.Sprintf(" WHEN :%s THEN :%s", tempPrimaryKey, tempFieldName)
				params[tempFieldName] = modelType.FieldByName(field.FieldName).Interface()
			}
		}
	}

	for columnName, str := range fieldSetArr {
		setStr += fmt.Sprintf("%s = CASE %s%s ELSE %s END,", columnName, primaryField.ColumnName, str, columnName)
	}
	setStr = strings.Trim(setStr, ",")
	inStr = strings.Trim(inStr, ",")
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s IN(%s)", tableName, setStr, primaryField.ColumnName, inStr)
	affected, _, err = NewNativeQuery(q.DatabaseConfigKey, q.TransactionKey).Execute(sql, params)
	return affected, err
}

// UpdateBy update by where
func (q *Query) UpdateBy(table interface{}, setParams []interface{}, where drivers.IWhere) (affected int, err error) {
	common.NewError().Defer(&err)
	if len(setParams) == 0 {
		return 0, nil
	}
	if where == nil {
		common.NewError().Throw("where condition is missing")
	}

	tableName, _ := model.NewMapping().GetSQLTableName(table)

	tempWhere := where.(*Where)
	setStr := ""
	for i := 0; i < len(setParams); i++ {
		paramType := reflect.TypeOf(setParams[i]).String()
		if paramType == "string" {
			setStr += setParams[i].(string) + ","
		} else if paramType == "common.Condition" {
			sql, params := getCondition(setParams[i].(common.Condition), tempWhere.Params)
			tempWhere.Params = params
			setStr += sql + ","
		}
	}
	setStr = strings.Trim(setStr, ",")
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, setStr, tempWhere.SQL)
	affected, _, err = NewNativeQuery(q.DatabaseConfigKey, q.TransactionKey).Execute(sql, tempWhere.Params)
	return affected, err
}

// Delete delete data
func (q *Query) Delete(modelInstance interface{}) (affected int, err error) {
	common.NewError().Defer(&err)
	mapData, tableName := model.NewMapping().GetModelData(modelInstance)
	whereStr := ""
	params := make(map[string]interface{})
	for _, field := range mapData {
		if field.IsPrimaryKey {
			if whereStr == "" {
				whereStr += field.ColumnName + "=:" + field.FieldName
			} else {
				whereStr += " AND " + field.ColumnName + "=:" + field.FieldName
			}
			params[field.FieldName] = field.Value
		}
	}

	if whereStr == "" {
		common.NewError().Throw("Update method require at least one primary key")
	}

	sql := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, whereStr)
	affected, _, err = NewNativeQuery(q.DatabaseConfigKey, q.TransactionKey).Execute(sql, params)
	return affected, err
}

// DeleteBy delete by where
func (q *Query) DeleteBy(table interface{}, where drivers.IWhere) (affected int, err error) {
	common.NewError().Defer(&err)
	if where == nil {
		common.NewError().Throw("where condition is missing")
	}

	tableName, aliasTableName := model.NewMapping().GetSQLTableName(table)

	tempWhere := where.(*Where)
	sql := fmt.Sprintf("DELETE %s FROM %s WHERE %s", aliasTableName, tableName, tempWhere.SQL)
	affected, _, err = NewNativeQuery(q.DatabaseConfigKey, q.TransactionKey).Execute(sql, tempWhere.Params)
	return affected, err
}

// QueryBuilder reurn query builder
func (q *Query) QueryBuilder() drivers.IQueryBuilder {
	return NewQueryBuilder(q.DatabaseConfigKey, q.TransactionKey)
}

// NativeQuery reurn native query
func (q *Query) NativeQuery() drivers.INativeQuery {
	return NewNativeQuery(q.DatabaseConfigKey, q.TransactionKey)
}

// CreateWhere reurn where
func (q *Query) CreateWhere() drivers.IWhere {
	return NewWhere()
}

// CreateOderBy reurn order by
func (q *Query) CreateOderBy() drivers.IOrderBy {
	return NewOrderBy()
}

// CreateGroupBy reurn group by
func (q *Query) CreateGroupBy() drivers.IGroupBy {
	return NewGroupBy()
}

func (q *Query) Begin() (err error) {
	common.NewError().Defer(&err)
	q.TransactionKey = NewNative(q.DatabaseConfigKey).Begin()
	return err
}

func (q *Query) Commit() (err error) {
	common.NewError().Defer(&err)
	NewNative(q.DatabaseConfigKey).Commit(q.TransactionKey)
	q.TransactionKey = ""
	return err
}

func (q *Query) Rollback() (err error) {
	common.NewError().Defer(&err)
	NewNative(q.DatabaseConfigKey).Rollback(q.TransactionKey)
	q.TransactionKey = ""
	return err
}

func (q *Query) getSelectSQL(table interface{}, args ...interface{}) (string, map[string]interface{}) {
	tableName, _ := model.NewMapping().GetSQLTableName(table)
	sql := "SELECT * FROM " + tableName
	params := make(map[string]interface{})
	if len(args) == 0 {
		return sql, params
	}
	whereSQL := ""
	orderSQL := ""
	for _, arg := range args {
		if arg == nil {
			continue
		}
		argType := reflect.TypeOf(arg).String()
		if argType == "*mysql.Where" {
			where := arg.(*Where)
			whereSQL = where.SQL
			params = where.Params
		} else if argType == "*mysql.OrderBy" {
			order := arg.(*OrderBy)
			orderSQL += "," + order.SQL
		}
	}
	if whereSQL != "" {
		sql += " WHERE " + whereSQL
	}
	if orderSQL != "" {
		orderSQL = strings.Trim(orderSQL, ",")
		sql += " ORDER BY " + orderSQL
	}

	return sql, params
}

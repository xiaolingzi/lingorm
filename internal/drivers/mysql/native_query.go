package mysql

import (
	"math"
	"strconv"

	"github.com/xiaolingzi/lingorm/internal/common"
	"github.com/xiaolingzi/lingorm/model"
)

// NativeQuery struct
type NativeQuery struct {
	DatabaseConfigKey string
	TransactionKey    string
}

// NewNativeQuery the instance of NativeQuery
func NewNativeQuery(databaseConfigKey string, transactionKey string) *NativeQuery {
	var query NativeQuery
	query.DatabaseConfigKey = databaseConfigKey
	query.TransactionKey = transactionKey
	return &query
}

//Excute excute the native sql
func (q *NativeQuery) Execute(sql string, params map[string]interface{}) (affected int, id int, err error) {
	defer common.NewError().Defer(&err)
	count, id := NewNative(q.DatabaseConfigKey).Execute(sql, params, q.TransactionKey)
	return int(count), int(id), nil
}

// Find return all the rows that meet query criteria
func (q *NativeQuery) Find(sql string, params map[string]interface{}, slicePtr ...interface{}) (result interface{}, err error) {
	defer common.NewError().Defer(&err)
	result = q.getData(sql, params, slicePtr...)
	return result, nil
}

// FindPage return the page result
func (q *NativeQuery) FindPage(pageIndex int, pageSize int, sql string, params map[string]interface{}, slicePtr ...interface{}) (result common.PageResult, err error) {
	defer common.NewError().Defer(&err)
	result.PageIndex = pageIndex
	result.PageSize = pageSize

	result.TotalCount, err = q.FindCount(sql, params)
	if err != nil {
		return result, err
	}

	result.TotalPages = int(math.Ceil(float64(result.TotalCount) / float64(result.PageSize)))

	sqlData := "SELECT * FROM (" + sql + ") tmp LIMIT " + strconv.Itoa((pageIndex-1)*pageSize) + ", " + strconv.Itoa(pageSize)
	data := q.getData(sqlData, params, slicePtr...)
	result.Data = data
	return result, nil
}

// FindCount return the number of rows that meet query criteria
func (q *NativeQuery) FindCount(sql string, params map[string]interface{}) (count int, err error) {
	defer common.NewError().Defer(&err)
	sqlCount := "SELECT count(*) as num FROM (" + sql + ") tmp"
	countResult := q.getData(sqlCount, params)
	result := countResult.(([]map[string]string))
	count, _ = strconv.Atoi(result[0]["num"])
	return count, nil
}

// First return the first row that meet query criteria
func (q *NativeQuery) First(sql string, params map[string]interface{}, structPtr ...interface{}) (interface{}, error) {
	sql = sql + " LIMIT 1"
	list := q.getData(sql, params, structPtr...)
	return list.([]interface{})[0], nil
}

func (q *NativeQuery) getData(sql string, params map[string]interface{}, resultObj ...interface{}) interface{} {
	data := NewNative(q.DatabaseConfigKey).FetchAll(sql, params, q.TransactionKey)

	if len(resultObj) <= 0 {
		return data
	}
	if len(data) == 0 {
		return nil
	}

	result := model.NewMapping().GetMappingData(data, resultObj[0])
	return result
}

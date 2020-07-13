package mysql

import (
	"reflect"
	"strings"

	"github.com/xiaolingzi/lingorm/internal/drivers"
	"github.com/xiaolingzi/lingorm/model"
)

// OrderBy struct
type OrderBy struct {
	SQL string
}

// NewOrderBy the instance of OrderBy
func NewOrderBy() *OrderBy {
	var order OrderBy
	return &order
}

// By order by
func (order *OrderBy) By(args ...interface{}) drivers.IOrderBy {
	if len(args) == 0 {
		return order
	}
	orderList := []string{"ASC", "DESC"}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		refValue := reflect.ValueOf(arg)
		valueType := refValue.Type().String()
		if valueType == "string" {
			order.SQL += "," + arg.(string)
		} else if valueType == "*model.Field" {
			tempValue := arg.(*model.Field)
			order.SQL += "," + tempValue.AliasTableName + "." + tempValue.ColumnName + " " + orderList[tempValue.OrderBy]
		} else if valueType == "model.Field" {
			tempValue := arg.(model.Field)
			order.SQL += "," + tempValue.AliasTableName + "." + tempValue.ColumnName + " " + orderList[tempValue.OrderBy]
		} else if valueType == "*mysql.OrderBy" {
			order.SQL += "," + arg.(*OrderBy).SQL
		}
	}
	order.SQL = strings.Trim(order.SQL, ",")
	return order
}

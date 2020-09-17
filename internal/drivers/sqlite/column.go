package sqlite

import (
	"reflect"
	"strings"

	"github.com/xiaolingzi/lingorm/model"
)

type Column struct {
}

func NewColumn() *Column {
	var column Column
	return &column
}

func (column *Column) GetSelectColumns(args ...interface{}) string {
	if len(args) == 0 {
		return "*"
	}
	result := ""
	for i := 0; i < len(args); i++ {
		columnStr := ""
		argType := reflect.TypeOf(args[i]).String()
		if argType == "string" {
			columnStr += args[i].(string)
		} else if argType == "model.Field" {
			arg := args[i].(model.Field)
			columnStr += arg.AliasTableName + "." + arg.ColumnName
			if arg.IsDistinct {
				columnStr = "DISTINCT " + columnStr
			}
			for j := 0; j < len(arg.ColumnsFunc); j++ {
				columnStr = arg.ColumnsFunc[j] + "(" + columnStr + ")"
			}
			if arg.AliasFieldName != "" {
				columnStr += " AS " + arg.AliasFieldName
			}
		} else if argType == "*model.Field" {
			arg := args[i].(*model.Field)
			columnStr += arg.AliasTableName + "." + arg.ColumnName
			if arg.IsDistinct {
				columnStr = "DISTINCT " + columnStr
			}
			for j := 0; j < len(arg.ColumnsFunc); j++ {
				columnStr = arg.ColumnsFunc[j] + "(" + columnStr + ")"
			}
			if arg.AliasFieldName != "" {
				columnStr += " AS " + arg.AliasFieldName
			}
		}
		result += columnStr + ","
	}
	result = strings.Trim(result, ",")
	if result == "" {
		result = "*"
	}
	return result
}

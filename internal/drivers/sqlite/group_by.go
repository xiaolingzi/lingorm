package sqlite

import (
	"reflect"
	"strings"

	"github.com/xiaolingzi/lingorm/internal/drivers"
	"github.com/xiaolingzi/lingorm/model"
)

// GroupBy the GroupBy struct
type GroupBy struct {
	SQL string
}

// NewGroupBy a
func NewGroupBy() *GroupBy {
	var group GroupBy
	return &group
}

// By group by
func (group *GroupBy) By(args ...interface{}) drivers.IGroupBy {
	if len(args) == 0 {
		return group
	}
	groupStr := ""
	for i := 0; i < len(args); i++ {
		argType := reflect.TypeOf(args[i]).String()
		if argType == "string" {
			groupStr = args[i].(string) + ","
		} else if argType == "model.Field" {
			arg := args[i].(model.Field)
			groupStr = arg.AliasTableName + "." + arg.ColumnName + ","
		} else if argType == "sqlite.GroupBy" {
			arg := args[i].(GroupBy)
			groupStr = arg.SQL + ","
		}
	}
	groupStr = group.SQL + "," + groupStr
	groupStr = strings.Trim(groupStr, ",")
	group.SQL = groupStr
	return group
}

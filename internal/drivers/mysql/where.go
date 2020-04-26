package mysql

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xiaolingzi/lingorm/internal/common"
	"github.com/xiaolingzi/lingorm/internal/drivers"
	"github.com/xiaolingzi/lingorm/model"
)

var paramsIndex int = 0

// Where struct
type Where struct {
	SQL    string
	Params map[string]interface{}
}

// NewWhere the instance of Where
func NewWhere() *Where {
	var where Where
	where.SQL = ""
	where.Params = make(map[string]interface{})
	return &where
}

// And and
func (p *Where) And(args ...interface{}) drivers.IWhere {
	if len(args) == 0 {
		return p
	}
	sql, params := getExpression(1, args, p.Params)
	var args2 []interface{}
	args2 = append(args2, p.SQL)
	args2 = append(args2, sql)
	p.SQL, params = getExpression(1, args2, params)
	p.Params = params
	return p
}

// GetAnd return the and string
func (p *Where) GetAnd(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}

	sql, params := getExpression(1, args, p.Params)
	p.Params = params
	return sql
}

// Or or
func (p *Where) Or(args ...interface{}) drivers.IWhere {
	if len(args) == 0 {
		return p
	}
	sql, params := getExpression(2, args, p.Params)
	var args2 []interface{}
	args2 = append(args2, p.SQL)
	args2 = append(args2, sql)
	p.SQL, params = getExpression(2, args2, params)
	p.Params = params
	return p
}

// GetOr return the or string
func (p *Where) GetOr(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}

	sql, params := getExpression(2, args, p.Params)
	p.Params = params
	return sql
}

// CurrentSQL return the current where sql
func (p *Where) CurrentSQL() string {
	return p.SQL
}

func getExpression(whereType int, args []interface{}, params map[string]interface{}) (string, map[string]interface{}) {
	if len(args) == 0 {
		return "", params
	}
	sql := ""
	count := len(args)
	for i := 0; i < len(args); i++ {
		tempSQL := ""
		argType := reflect.TypeOf(args[i]).String()
		if argType == "string" {
			tempSQL = args[i].(string)
		} else if argType == "*mysql.Where" {
			where := (args[i].(*Where))
			tempSQL = where.SQL
			for key, value := range where.Params {
				params[key] = value
			}
		} else if argType == "common.Condition" {
			tempSQL, params = getCondition(args[i].(common.Condition), params)
		}

		if tempSQL == "" {
			continue
		}

		if count > 1 {
			reg, _ := regexp.Compile("\\([^\\(\\)]*\\)")
			tempStr := reg.ReplaceAllString(tempSQL, "")
			for {
				if !reg.MatchString(tempStr) {
					break
				}
				tempStr = reg.ReplaceAllString(tempStr, "")
			}

			if (whereType == 1 && strings.Contains(tempStr, " OR ")) || (whereType == 2 && strings.Contains(tempStr, " AND ")) {
				tempSQL = "(" + tempSQL + ")"
			}
		}

		if sql == "" {
			sql = tempSQL
		} else {
			if whereType == 1 {
				sql += " AND " + tempSQL
			} else {
				sql += " OR " + tempSQL
			}
		}
	}

	// if len(sql) > 2 && sql[0] == '(' && sql[len(sql)-1] == ')' {
	// 	sql = sql[1 : len(sql)-2]
	// }

	return sql, params
}

func getCondition(condition common.Condition, params map[string]interface{}) (string, map[string]interface{}) {
	operator := getOperator(condition.Operator)
	sql := ""
	fieldName := condition.AliasTableName + "." + condition.ColumnName
	if condition.Operator == common.OperatorNull || condition.Operator == common.OperatorNotNull {
		sql = fieldName + " " + operator
	} else if condition.Value == nil {
		return sql, params
	} else {
		refValue := reflect.ValueOf(condition.Value)
		if refValue.Type().String() == "time.Time" {
			condition.Value = condition.Value.(time.Time).Format(common.TimeLayout)
		}
		if refValue.Type().String() == "model.Field" {
			tempValue := condition.Value.(model.Field)
			sql = fieldName + " " + operator + " " + tempValue.AliasTableName + "." + tempValue.ColumnName
		} else if condition.Operator == common.OperatorIn || condition.Operator == common.OperatorNotIn {
			inStr := ""
			tempValue := condition.Value
			var inValues []interface{}
			valKind := refValue.Kind()
			if valKind == reflect.String {
				tempInValues := strings.Split(tempValue.(string), ",")
				for _, v := range tempInValues {
					inValues = append(inValues, v)
				}

			} else if valKind == reflect.Slice || valKind == reflect.Array {
				for i := 0; i < refValue.Len(); i++ {
					if refValue.Index(i).Type().Name() == "Time" {
						inVlue := refValue.Index(i).Interface().(time.Time).Format(common.TimeLayout)
						inValues = append(inValues, inVlue)
					} else {
						inValues = append(inValues, refValue.Index(i).Interface())
					}
				}
			} else {
				return sql, params
			}
			for _, v := range inValues {
				if v != nil && v != "" {
					key := getParamKey()
					params[key] = v
					inStr += ":" + key + ","
				}
			}
			inStr = strings.Trim(inStr, ",")
			sql = fieldName + " " + operator + "(" + inStr + ")"
		} else if condition.Operator == common.OperatorFindInSet {
			key := getParamKey()
			params[key] = condition.Value
			sql = operator + "(:" + key + "," + fieldName + ")"
		} else {
			key := getParamKey()
			params[key] = condition.Value
			sql = fieldName + " " + operator + " :" + key
		}
	}

	return sql, params
}

func getOperator(operator string) string {
	switch operator {
	case common.OperatorEqual:
		return "="
	case common.OperatorNotEqual:
		return "<>"
	case common.OperatorGreaterThan:
		return ">"
	case common.OperatorGreaterEqual:
		return ">="
	case common.OperatorLessThan:
		return "<"
	case common.OperatorLessEqual:
		return "<="
	case common.OperatorIn:
		return "IN"
	case common.OperatorNotIn:
		return "NOT IN"
	case common.OperatorLike:
		return "like"
	case common.OperatorFindInSet:
		return "FIND_IN_SET"
	case common.OperatorNull:
		return "IS NULL"
	case common.OperatorNotNull:
		return "IS NOT NULL"
	default:
		return ""
	}
}

func getParamKey() string {
	if paramsIndex > 1000000 {
		paramsIndex = 0
	} else {
		paramsIndex++
	}
	key := "p" + strconv.Itoa(paramsIndex)
	return key
}

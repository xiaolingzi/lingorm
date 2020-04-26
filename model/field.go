package model

import (
	"github.com/xiaolingzi/lingorm/internal/common"
)

// Field struct
type Field struct {
	DB            string
	Table         string
	FieldName     string
	ColumnName    string
	ColumnType    string
	Length        int
	IsPrimaryKey  bool
	AutoIncrement bool
	Value         interface{}

	AliasTableName string
	AliasFieldName string

	IsDistinct  bool
	ColumnsFunc []string
	OrderBy     int
}

// F function
func (f *Field) F(functionName string) *Field {
	f.ColumnsFunc = append(f.ColumnsFunc, functionName)
	return f
}

//Alias column alias name
func (f *Field) Alias(aliasName string) *Field {
	f.AliasFieldName = aliasName
	return f
}

// Max max function of sql
func (f *Field) Max() *Field {
	f.ColumnsFunc = append(f.ColumnsFunc, "MAX")
	return f
}

// Min min function of sql
func (f *Field) Min() *Field {
	f.ColumnsFunc = append(f.ColumnsFunc, "MIN")
	return f
}

// Count count function of sql
func (f *Field) Count() *Field {
	f.ColumnsFunc = append(f.ColumnsFunc, "COUNT")
	return f
}

// Sum sum function of sql
func (f *Field) Sum() *Field {
	f.ColumnsFunc = append(f.ColumnsFunc, "SUM")
	return f
}

// Distinct distinct function of sql
func (f *Field) Distinct() *Field {
	f.IsDistinct = true
	return f
}

// EQ equal
func (f *Field) EQ(val interface{}) common.Condition {
	result := common.NewCondition().GetCondition(f.ColumnName, f.AliasTableName, val, common.OperatorEqual)
	return result
}

// NEQ not equal
func (f *Field) NEQ(val interface{}) common.Condition {
	result := common.NewCondition().GetCondition(f.ColumnName, f.AliasTableName, val, common.OperatorNotEqual)
	return result
}

// GT greater than
func (f *Field) GT(val interface{}) common.Condition {
	result := common.NewCondition().GetCondition(f.ColumnName, f.AliasTableName, val, common.OperatorGreaterThan)
	return result
}

// GE greater than or equal
func (f *Field) GE(val interface{}) common.Condition {
	result := common.NewCondition().GetCondition(f.ColumnName, f.AliasTableName, val, common.OperatorGreaterEqual)
	return result
}

// LT less than
func (f *Field) LT(val interface{}) common.Condition {
	result := common.NewCondition().GetCondition(f.ColumnName, f.AliasTableName, val, common.OperatorLessThan)
	return result
}

// LE less than or equal
func (f *Field) LE(val interface{}) common.Condition {
	result := common.NewCondition().GetCondition(f.ColumnName, f.AliasTableName, val, common.OperatorLessEqual)
	return result
}

// LIKE like
func (f *Field) LIKE(val interface{}) common.Condition {
	result := common.NewCondition().GetCondition(f.ColumnName, f.AliasTableName, val, common.OperatorLike)
	return result
}

// IN in
func (f *Field) IN(val interface{}) common.Condition {
	result := common.NewCondition().GetCondition(f.ColumnName, f.AliasTableName, val, common.OperatorIn)
	return result
}

// NIN not in
func (f *Field) NIN(val interface{}) common.Condition {
	result := common.NewCondition().GetCondition(f.ColumnName, f.AliasTableName, val, common.OperatorNotIn)
	return result
}

// FIS FIND_IN_SET
func (f *Field) FIS(val interface{}) common.Condition {
	result := common.NewCondition().GetCondition(f.ColumnName, f.AliasTableName, val, common.OperatorFindInSet)
	return result
}

// IsNull is null
func (f *Field) IsNull() common.Condition {
	result := common.NewCondition().GetCondition(f.ColumnName, f.AliasTableName, false, common.OperatorNull)
	return result
}

// IsNotNull is not null
func (f *Field) IsNotNull() common.Condition {
	result := common.NewCondition().GetCondition(f.ColumnName, f.AliasTableName, true, common.OperatorNotNull)
	return result
}

// ASC sort in ascending order
func (f *Field) ASC() *Field {
	f.OrderBy = 0
	return f
}

// DESC sort in descending order
func (f *Field) DESC() *Field {
	f.OrderBy = 1
	return f
}

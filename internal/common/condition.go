package common

// Condition struct
type Condition struct {
	AliasTableName string
	ColumnName     string
	Operator       string
	Value          interface{}
}

// NewCondition the instance of Condition
func NewCondition() *Condition {
	var condition Condition
	return &condition
}

// GetCondition return conditon
func (c *Condition) GetCondition(columnName string, aliasTableName string, val interface{}, operator string) Condition {
	c.ColumnName = columnName
	c.Operator = operator
	c.Value = val
	c.AliasTableName = aliasTableName
	return *c
}

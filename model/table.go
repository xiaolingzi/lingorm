package model

// Table return model table
func Table(table interface{}, model interface{}, tableName string, databaseName string) interface{} {
	return NewMapping().GetTable(table, model, tableName, databaseName)
}

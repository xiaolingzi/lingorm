package model

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/xiaolingzi/lingorm/internal/common"
	"github.com/xiaolingzi/lingorm/internal/utils/cryptography"
)

// Mapping struct
type Mapping struct {
}

type mappingCacheCollection struct {
	FieldMappingCache map[string]map[string]Field
	TableCache        map[string]interface{}
	ModelTypeCache    map[string]reflect.Type
}

var mappingCache mappingCacheCollection
var tableIndex int = 0

// NewMapping return instance of Mapping
func NewMapping() *Mapping {
	var p Mapping
	return &p
}

// DocParser doc parser
func (p *Mapping) DocParser(mapResult map[string]string, modelInstance interface{}) (interface{}, error) {
	var err error = nil
	refValue := reflect.ValueOf(modelInstance)
	modelType := refValue.Type()
	result := reflect.New(modelType).Elem()

	if len(mapResult) == 0 {
		return result.Interface(), err
	}

	refType := reflect.TypeOf(modelInstance)
	mappings, err := p.getFieldMappings(refType, refValue)
	if err != nil {
		return result.Interface(), err
	}

	fieldCount := refType.NumField()
	for i := 0; i < fieldCount; i++ {
		columnName := mappings[refType.Field(i).Name].ColumnName
		if _, ok := mapResult[columnName]; ok {
			tempValue := convertStringToType(mapResult[columnName], refType.Field(i).Type)
			refTargetValue := reflect.ValueOf(tempValue)
			if refTargetValue.IsValid() {
				if refTargetValue.Type().ConvertibleTo(refValue.Field(i).Type()) {
					result.Field(i).Set(refTargetValue.Convert(refValue.Field(i).Type()))
					continue
				} else {
					err = fmt.Errorf("Could not convert argument of field %s from %s to %s", refType.Field(i).Name, refTargetValue.Type(), refValue.Field(i).Type())
				}
			}
		}
		result.Field(i).Set(reflect.Zero(refValue.Field(i).Type()))
	}
	return result.Interface(), err
}

func (p *Mapping) getFieldMappings(refType reflect.Type, refValue reflect.Value) (map[string]Field, error) {
	var err error
	result := make(map[string]Field)
	key := cryptography.MD5(refType.PkgPath() + "/" + refType.Name())
	if _, ok := mappingCache.FieldMappingCache[key]; !ok {
		if !refValue.MethodByName("Table").IsValid() {
			p.setFieldMappingCache(refType, refValue, key)
			// err = errors.New("Invalid model")
			// log.Fatal(err)
			// return result, err
		} else {
			refValue.MethodByName("Table").Call(nil)
		}

	}

	result = mappingCache.FieldMappingCache[key]
	return result, err
}

//GetModelData return model data and table name
func (p *Mapping) GetModelData(modelInstance interface{}) (map[string]Field, string) {
	dataMap := make(map[string]Field)
	tableName := ""

	refType := reflect.TypeOf(modelInstance)
	refValue := reflect.ValueOf(modelInstance)
	mappings, err := p.getFieldMappings(refType, refValue)
	if err != nil {
		log.Fatalln(err)
	}
	for fieldName, field := range mappings {
		fieldValue := refValue.FieldByName(fieldName).Interface()
		columnName := field.ColumnName
		field.Value = fieldValue
		dataMap[columnName] = field
	}

	key := cryptography.MD5(refType.PkgPath() + "/" + refType.Name())
	table := mappingCache.TableCache[key]

	tbRefValue := reflect.ValueOf(table)

	tableName = tbRefValue.FieldByName("TTTableName").String()
	if tableName == "" {
		log.Fatal("Invalid model")
	}

	databaseName := tbRefValue.FieldByName("TTDatabaseName").String()
	if databaseName != "" {
		tableName = databaseName + "." + tableName
	}

	return dataMap, tableName
}

// GetTable return table instande
func (p *Mapping) GetTable(table interface{}, modelInstance interface{}, tableName string, databaseName string) interface{} {
	var result interface{}
	refType := reflect.TypeOf(modelInstance)
	key := cryptography.MD5(refType.PkgPath() + "/" + refType.Name())
	if _, ok := mappingCache.TableCache[key]; !ok {
		p.setMappingCache(table, key, refType, tableName, databaseName)
	}
	result = mappingCache.TableCache[key]
	return result
}

// GetSQLTableName return the sql table name
func (p *Mapping) GetSQLTableName(table interface{}) (string, string) {
	refValue := reflect.ValueOf(table)

	result := ""
	tableName := refValue.FieldByName("TTTableName").String()
	if tableName == "" {
		log.Fatal("Invalid model")
	}
	result = tableName

	databaseName := refValue.FieldByName("TTDatabaseName").String()
	if databaseName != "" {
		result = databaseName + "." + tableName
	}

	aliasName := refValue.FieldByName("TTAlias").String()
	if aliasName != "" {
		result = result + " " + aliasName
	}
	return result, aliasName
}

// GetModelType get the type of model by table instance
func (p *Mapping) GetModelType(table interface{}) reflect.Type {
	refType := reflect.TypeOf(table)
	tableKey := cryptography.MD5(refType.PkgPath() + "/" + refType.Name())
	return mappingCache.ModelTypeCache[tableKey]
}

func (p *Mapping) setMappingCache(tb interface{}, key string, modelType reflect.Type, tableName string, databaseName string) {
	refType := reflect.TypeOf(tb)
	refValue := reflect.ValueOf(tb)
	table := reflect.New(refValue.Type()).Elem()
	aliasTableName := getTableAlias()

	table.FieldByName("TTDatabaseName").SetString(databaseName)
	table.FieldByName("TTTableName").SetString(tableName)
	table.FieldByName("TTAlias").SetString(aliasTableName)

	fieldMapping := make(map[string]Field)
	fieldCount := modelType.NumField()
	for i := 0; i < fieldCount; i++ {
		fieldTag := modelType.Field(i).Tag
		fieldName := modelType.Field(i).Name
		var tempField Field

		tempField.AliasTableName = aliasTableName
		tempField.FieldName = fieldName
		tempField.ColumnName = fieldName

		columnName := fieldTag.Get("column")
		if columnName != "" {
			tempField.ColumnName = columnName
		}

		isPrimaryKey := strings.ToLower(fieldTag.Get("primary_key"))
		if isPrimaryKey == "true" || isPrimaryKey == "1" {
			tempField.IsPrimaryKey = true
		} else {
			tempField.IsPrimaryKey = false
		}

		autoIncrement := strings.ToLower(fieldTag.Get("auto_increment"))
		if autoIncrement == "true" || autoIncrement == "1" {
			tempField.AutoIncrement = true
		} else {
			tempField.AutoIncrement = false
		}

		refTargetValue := reflect.ValueOf(tempField)
		if table.FieldByName(fieldName).IsValid() {
			table.FieldByName(fieldName).Set(refTargetValue.Convert(table.FieldByName(fieldName).Type()))
		}
		fieldMapping[fieldName] = tempField
	}

	if mappingCache.FieldMappingCache == nil {
		mappingCache.FieldMappingCache = make(map[string]map[string]Field)
	}
	mappingCache.FieldMappingCache[key] = fieldMapping

	if mappingCache.TableCache == nil {
		mappingCache.TableCache = make(map[string]interface{})
	}
	mappingCache.TableCache[key] = table.Interface()

	if mappingCache.ModelTypeCache == nil {
		mappingCache.ModelTypeCache = make(map[string]reflect.Type)
	}

	tableKey := cryptography.MD5(refType.PkgPath() + "/" + refType.Name())
	mappingCache.ModelTypeCache[tableKey] = modelType
}

func (p *Mapping) setFieldMappingCache(refType reflect.Type, refValue reflect.Value, key string) {
	fieldMapping := make(map[string]Field)
	fieldCount := refType.NumField()
	for i := 0; i < fieldCount; i++ {
		fieldTag := refType.Field(i).Tag
		fieldName := refType.Field(i).Name
		var tempField Field

		tempField.FieldName = fieldName
		tempField.ColumnName = fieldName
		columnName := fieldTag.Get("column")
		if columnName != "" {
			tempField.ColumnName = columnName
		}

		isPrimaryKey := strings.ToLower(fieldTag.Get("primary_key"))
		if isPrimaryKey == "true" || isPrimaryKey == "1" {
			tempField.IsPrimaryKey = true
		} else {
			tempField.IsPrimaryKey = false
		}
		autoIncrement := strings.ToLower(fieldTag.Get("auto_increment"))
		if autoIncrement == "true" || autoIncrement == "1" {
			tempField.AutoIncrement = true
		} else {
			tempField.AutoIncrement = false
		}

		fieldMapping[fieldName] = tempField
	}

	if mappingCache.FieldMappingCache == nil {
		mappingCache.FieldMappingCache = make(map[string]map[string]Field)
	}
	mappingCache.FieldMappingCache[key] = fieldMapping
}

func convertStringToType(val string, fiedType reflect.Type) interface{} {
	var result interface{}

	switch fiedType.Kind() {
	case reflect.Int:
		result, _ = strconv.ParseInt(val, 10, 0)
	case reflect.Int8:
		result, _ = strconv.ParseInt(val, 10, 8)
	case reflect.Int32:
		result, _ = strconv.ParseInt(val, 10, 32)
	case reflect.Int64:
		result, _ = strconv.ParseInt(val, 10, 64)
	case reflect.Uint:
		result, _ = strconv.ParseUint(val, 10, 0)
	case reflect.Uint8:
		result, _ = strconv.ParseUint(val, 10, 8)
	case reflect.Uint32:
		result, _ = strconv.ParseUint(val, 10, 32)
	case reflect.Uint64:
		result, _ = strconv.ParseUint(val, 10, 64)
	case reflect.Float32:
		result, _ = strconv.ParseFloat(val, 32)
	case reflect.Float64:
		result, _ = strconv.ParseFloat(val, 64)
	case reflect.Struct:
		if fiedType.String() == "time.Time" {
			result, _ = time.ParseInLocation(common.TimeLayout, val, time.Local)
		}
	default:
		result = val
	}
	return result
}

func getTableAlias() string {
	tableIndex++
	if tableIndex > 1000000 {
		tableIndex = 0
	}
	return "t" + strconv.Itoa(tableIndex)
}

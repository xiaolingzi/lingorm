package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/xiaolingzi/lingorm"
	"github.com/xiaolingzi/lingorm/configs"
)

var query lingorm.INativeQuery
var database, table string

func main() {
	configs.SetEnvConfigs()

	flag.StringVar(&database, "d", "", "database")
	flag.StringVar(&table, "t", "", "table name")
	flag.Parse()
	for {
		if database != "" {
			break
		}
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("input database>")
		input, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		database = strings.Trim(input, "\n")
	}
	// fmt.Println(database)

	query = lingorm.DB("schema").NativeQuery()
	if table != "" {
		generateEntity(database, table)
	} else {
		tableList := getTables(database)
		for i := 0; i < len(tableList); i++ {
			tempTableName := ""
			if _, ok := tableList[i]["Table_name"]; ok {
				tempTableName = tableList[i]["Table_name"]
			} else {
				tempTableName = tableList[i]["TABLE_NAME"]
			}
			generateEntity(database, tempTableName)
		}
	}

}

func generateEntity(db string, tb string) {
	// packageName := getPackageName(db)
	packageName := "models"
	lowerCamelEntityName := getLowerCamelCaseEntityName(tb)
	upperCamelEntityName := getUpperCamelCaseEntityName(tb)

	content := getTemplateContent()
	// fmt.Println(content)
	content = strings.ReplaceAll(content, "{{database_name}}", db)
	content = strings.ReplaceAll(content, "{{table_name}}", tb)
	content = strings.ReplaceAll(content, "{{package}}", packageName)
	content = strings.ReplaceAll(content, "{{lower_camel_entity_name}}", lowerCamelEntityName)
	content = strings.ReplaceAll(content, "{{upper_camel_entity_name}}", upperCamelEntityName)

	columns := getCloumns(db, tb)
	reg, _ := regexp.Compile("(?s)<<<(.*?)>>>")
	content = reg.ReplaceAllStringFunc(content, func(str string) string {
		maxFieldLen := 0
		maxTypeLen := 0
		for i := 0; i < len(columns); i++ {
			column := columns[i]
			columnName := ""
			if _, ok := column["Column_name"]; ok {
				columnName = column["Column_name"]
			} else {
				columnName = column["COLUMN_NAME"]
			}
			column["COLUMN_NAME"] = columnName
			fieldName := getUpperCamelFieldName(columnName)
			if len(fieldName) > maxFieldLen {
				maxFieldLen = len(fieldName)
			}
			column["FIELD_NAME"] = fieldName

			fieldType := getDataType(column["DATA_TYPE"])
			if len(fieldType) > maxTypeLen {
				maxTypeLen = len(fieldType)
			}
			column["FIELD_TYPE"] = fieldType
		}

		regStr := ""
		for i := 0; i < len(columns); i++ {
			tempStr := str
			column := columns[i]
			columnName := column["COLUMN_NAME"]
			upperFieldName := column["FIELD_NAME"]
			lowerFieldName := getLowerCamelCaseFieldName(upperFieldName)
			tempLen := maxFieldLen - len(upperFieldName)
			for j := 0; j < tempLen; j++ {
				upperFieldName += " "
			}

			fieldType := column["FIELD_TYPE"]
			tempLen = maxTypeLen - len(fieldType)
			for k := 0; k < tempLen; k++ {
				fieldType += " "
			}

			columnProperty := `json:"` + lowerFieldName + `"`
			columnProperty += ` comlumn:"` + columnName + `"`
			// columnProperty += ` type:"` + column["DATA_TYPE"] + `"`
			if column["COLUMN_KEY"] == "PRI" {
				columnProperty += ` primary_key:"true"`
			}
			if column["EXTRA"] == "auto_increment" {
				columnProperty += ` auto_increment:"true"`
			}
			tempStr = strings.ReplaceAll(tempStr, "<<<", "")
			tempStr = strings.ReplaceAll(tempStr, ">>>", "")
			tempStr = strings.ReplaceAll(tempStr, "{{upper_camel_field_name}}", upperFieldName)
			tempStr = strings.ReplaceAll(tempStr, "{{lower_camel_field_name}}", lowerFieldName)
			tempStr = strings.ReplaceAll(tempStr, "{{field_type}}", fieldType)
			tempStr = strings.ReplaceAll(tempStr, "{{column_property}}", columnProperty)
			regStr += tempStr + "\n"
		}
		// fmt.Println(regStr)
		regStr = strings.Trim(regStr, "\n")
		return regStr
	})
	saveFile(db, tb, content)
	fmt.Println(tb + " generated!")
}

func getTables(db string) []map[string]string {
	sql := fmt.Sprintf("select TABLE_NAME from TABLES where TABLE_SCHEMA='%s'", db)
	result, _ := query.Find(sql, nil)
	return result.([]map[string]string)
}

func getCloumns(db string, tb string) []map[string]string {
	sql := fmt.Sprintf("select COLUMN_NAME,DATA_TYPE, CHARACTER_MAXIMUM_LENGTH,COLUMN_KEY,EXTRA from COLUMNS where TABLE_SCHEMA='%s' and TABLE_NAME='%s'", db, tb)
	result, _ := query.Find(sql, nil)
	return result.([]map[string]string)
}

func getPackageName(db string) string {
	packageName := strings.ToLower(db)
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	packageName = reg.ReplaceAllString(packageName, "")
	return packageName
}

func getUpperCamelCaseEntityName(tb string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+([a-zA-Z0-9]{1})")
	entityName := reg.ReplaceAllStringFunc(tb, strings.ToUpper)
	reg2, _ := regexp.Compile("[^a-zA-Z0-9]+")
	entityName = reg2.ReplaceAllString(entityName, "")
	entityName = strings.ToUpper(string(entityName[0])) + string(entityName[1:len(entityName)])
	entityName = dealSpecialName(entityName)
	return entityName
}

func getLowerCamelCaseEntityName(tb string) string {
	upperEntityName := getUpperCamelCaseEntityName(tb)
	reg, _ := regexp.Compile("^[A-Z]+")
	entityName := reg.ReplaceAllStringFunc(upperEntityName, strings.ToLower)
	return entityName
}

func getLowerEntityName(tb string) string {
	upperEntityName := getUpperCamelCaseEntityName(tb)
	reg, _ := regexp.Compile("[A-Z]+")
	entityName := reg.ReplaceAllStringFunc(upperEntityName, func(s string) string {
		return "_" + strings.ToLower(s)
	})
	entityName = strings.Trim(entityName, "_")
	return entityName
}

func getUpperCamelFieldName(columnName string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+([a-zA-Z0-9]{1})")
	fieldName := reg.ReplaceAllStringFunc(columnName, strings.ToUpper)
	reg2, _ := regexp.Compile("[^a-zA-Z0-9]+")
	fieldName = reg2.ReplaceAllString(fieldName, "")
	fieldName = strings.ToUpper(string(fieldName[0])) + string(fieldName[1:len(fieldName)])
	fieldName = dealSpecialName(fieldName)
	return fieldName
}

func getLowerCamelCaseFieldName(columnName string) string {
	upperFieldName := getUpperCamelFieldName(columnName)
	reg, _ := regexp.Compile("^[A-Z]+")
	fieldName := reg.ReplaceAllStringFunc(upperFieldName, strings.ToLower)
	return fieldName
}

func dealSpecialName(camelCaseName string) string {
	if camelCaseName == "" {
		return camelCaseName
	}

	tempMap := make(map[string]string)
	tempMap["Id"] = "ID"
	tempMap["Url"] = "URL"
	tempMap["Ip"] = "IP"

	reg, _ := regexp.Compile("[A-Z][^A-Z]+")
	result := reg.ReplaceAllStringFunc(camelCaseName, func(s string) string {
		if _, ok := tempMap[s]; ok {
			return tempMap[s]
		}
		return s
	})
	return result
}

func getDataType(dbType string) string {
	dbType = strings.ToLower(dbType)
	dbTypeMap := make(map[string]string)
	dbTypeMap["char"] = "string"
	dbTypeMap["varchar"] = "string"
	dbTypeMap["nvarchar"] = "string"
	dbTypeMap["text"] = "string"
	dbTypeMap["longtext"] = "string"
	dbTypeMap["tinytext"] = "string"
	dbTypeMap["json"] = "string"
	dbTypeMap["mediumtext"] = "string"
	dbTypeMap["int"] = "int"
	dbTypeMap["smallint"] = "int"
	dbTypeMap["tinyint"] = "int"
	dbTypeMap["bigint"] = "int"
	dbTypeMap["mediumint"] = "int"
	dbTypeMap["datetime"] = "time.Time"
	dbTypeMap["date"] = "time.Time"
	dbTypeMap["time"] = "time.Time"
	dbTypeMap["timestamp"] = "time.Time"
	dbTypeMap["year"] = "string"
	dbTypeMap["double"] = "float64"
	dbTypeMap["decimal"] = "float64"
	dbTypeMap["float"] = "float64"
	dbTypeMap["boolean"] = "bool"
	if _, ok := dbTypeMap[dbType]; ok {
		return dbTypeMap[dbType]
	}
	return ""
}

func getTemplateContent() string {
	dir, _ := os.Getwd()
	filename := dir + "/tools/template.txt"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func saveFile(databaseName string, tableName string, content string) {
	dir, _ := os.Getwd()
	dir = dir + "/tools/" + databaseName
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, 0777)
	}
	lowerEntityName := getLowerEntityName(tableName)
	filename := dir + "/" + lowerEntityName + "_entity.go"
	fmt.Println(filename)
	err = ioutil.WriteFile(filename, []byte(content), 0666)
	if err != nil {
		fmt.Println(err)
	}
}

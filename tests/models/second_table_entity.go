package models

import (
	"github.com/xiaolingzi/lingorm/model"
	"time"
)

//SecondTableEntity entity
type SecondTableEntity struct {
    ID           int       `json:"id" column:"id" primary_key:"true" auto_increment:"true"`
    SecondName   string    `json:"secondName" column:"second_name"`
    SecondNumber int       `json:"secondNumber" column:"second_number"`
    SecondTime   time.Time `json:"secondTime" column:"second_time"`
}

//SecondTableTable table
type SecondTableTable struct {
    TTDatabaseName      string
    TTTableName         string
    TTAlias             string
    ID           model.Field
    SecondName   model.Field
    SecondNumber model.Field
    SecondTime   model.Field
}

//Table return table
func (e SecondTableEntity) Table() SecondTableTable {
	return model.Table(SecondTableTable{}, SecondTableEntity{}, "second_table", "test").(SecondTableTable)
}

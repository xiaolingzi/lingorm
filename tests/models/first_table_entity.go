package models

import (
	"github.com/xiaolingzi/lingorm/model"
	"time"
)

//FirstTableEntity entity
type FirstTableEntity struct {
    ID          int       `json:"id" column:"id" primary_key:"true" auto_increment:"true"`
    FirstName   string    `json:"firstName" column:"first_name"`
    FirstNumber int       `json:"firstNumber" column:"first_number"`
    FirstTime   time.Time `json:"firstTime" column:"first_time"`
}

//FirstTableTable table
type FirstTableTable struct {
    TTDatabaseName      string
    TTTableName         string
    TTAlias             string
    ID          model.Field
    FirstName   model.Field
    FirstNumber model.Field
    FirstTime   model.Field
}

//Table return table
func (e FirstTableEntity) Table() FirstTableTable {
	return model.Table(FirstTableTable{}, FirstTableEntity{}, "first_table", "test").(FirstTableTable)
}

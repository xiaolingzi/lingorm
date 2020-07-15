package mysql

import (
	"strconv"
	"testing"
	"time"

	"github.com/xiaolingzi/lingorm/tests/configs"
	"github.com/xiaolingzi/lingorm/tests/models"
)

func TestCUD(t *testing.T) {
	db := configs.GetDB()
	db.Begin()
	t.Run("TestInsert", func(t *testing.T) {
		entity := models.FirstTableEntity{}
		entity.FirstName = "go name"
		entity.FirstNumber = 1001
		entity.FirstTime = time.Now()

		id, err := db.Insert(entity)
		if err != nil {
			t.Errorf("Insert error")
		}
		if id <= 0 {
			t.Errorf("The value of id %d is not greater than 0", id)
		}

		exists := models.FirstTableEntity{}
		table := models.FirstTableEntity{}.Table()
		db.Table(table).Where(table.ID.EQ(id)).First(&exists)
		if (exists == models.FirstTableEntity{} || exists.FirstName != "go name") {
			t.Errorf("The data is not correct")
		}
	})
	t.Run("TestBatchInsert", func(t *testing.T) {
		var entityList []interface{}
		for i := 1; i < 3; i++ {
			entity := models.FirstTableEntity{}
			entity.FirstName = "go batch name " + strconv.Itoa(i)
			entity.FirstNumber = 2001
			entity.FirstTime = time.Now()
			entityList = append(entityList, entity)
		}

		affected, err := db.BatchInsert(entityList)
		if err != nil {
			t.Errorf("Insert error")
		}
		if affected != 2 {
			t.Errorf("The value of affected %d is not equal to 2", affected)
		}

		var exists []models.FirstTableEntity
		table := models.FirstTableEntity{}.Table()
		db.Table(table).Where(table.FirstName.LIKE("go batch name%")).Find(&exists)
		if exists == nil || len(exists) == 0 || exists[0].FirstName != "go batch name 1" {
			t.Errorf("The data is not correct")
		}
	})

	t.Run("TestUpdate", func(t *testing.T) {
		entity := models.FirstTableEntity{}
		table := models.FirstTableEntity{}.Table()
		db.Table(table).Where(table.FirstNumber.EQ(1001)).First(&entity)
		if (entity == models.FirstTableEntity{}) {
			t.Errorf("Data not found")
		}

		entity.FirstName = "go update name"
		entity.FirstNumber = 1002
		db.Update(entity)

		var exists models.FirstTableEntity
		db.Table(table).Where(table.FirstName.EQ("go update name")).First(&exists)
		if (exists == models.FirstTableEntity{} || exists.FirstName != "go update name") {
			t.Errorf("Data not updated")
		}

	})
	t.Run("TestBatchUpdate", func(t *testing.T) {
		var entityList []models.FirstTableEntity
		table := models.FirstTableEntity{}.Table()
		db.Table(table).Where(table.FirstName.LIKE("go batch name%")).Find(&entityList)
		for i := 0; i < len(entityList); i++ {
			entityList[i].FirstName = "go batch update name"
			entityList[i].FirstNumber = 2003
		}
		db.BatchUpdate(entityList)

		var exists models.FirstTableEntity
		db.Table(table).Where(table.FirstName.EQ("go batch update name")).First(&exists)
		if (exists == models.FirstTableEntity{} || exists.FirstName != "go batch update name") {
			t.Errorf("Data not updated")
		}

	})
	t.Run("TestUpdateBy", func(t *testing.T) {
		table := models.FirstTableEntity{}.Table()
		where := db.CreateWhere().And(table.FirstNumber.EQ(2003), table.FirstName.EQ("go batch update name"))
		var params = []interface{}{
			table.FirstName.EQ("go new update name"),
			table.FirstNumber.EQ(2004),
		}
		affected, err := db.UpdateBy(table, params, where)
		if err != nil {
			t.Errorf("UpdateBy error")
		}
		if affected <= 0 {
			t.Errorf("Nothing updated")
		}

		where = db.CreateWhere().And(table.FirstNumber.EQ(2004), table.FirstName.EQ("go new update name"))
		var entityList []models.FirstTableEntity
		db.Table(table).Where(where).Find(&entityList)
		if entityList == nil || len(entityList) != affected {
			t.Errorf("Data not updated")
		}
	})
	t.Run("Delete", func(t *testing.T) {
		entity := models.FirstTableEntity{}
		table := models.FirstTableEntity{}.Table()
		db.Table(table).Where(table.FirstNumber.EQ(1002)).First(&entity)
		if (entity == models.FirstTableEntity{}) {
			t.Errorf("Data not found")
		}

		affected, err := db.Delete(entity)
		if err != nil {
			t.Errorf("Delete error")
		}

		if affected <= 0 {
			t.Errorf("Nothing deleted")
		}

		var exists []models.FirstTableEntity
		db.Table(table).Where(table.FirstNumber.EQ(1002)).Find(&exists)
		if exists != nil && len(exists) > 0 {
			t.Errorf("Data not deleted")
		}

	})

	t.Run("TestDeleteBy", func(t *testing.T) {
		table := models.FirstTableEntity{}.Table()
		where := db.CreateWhere().And(table.FirstNumber.EQ(2004), table.FirstName.EQ("go new update name"))
		affected, err := db.DeleteBy(table, where)
		if err != nil {
			t.Errorf("DeleteBy error")
		}
		if affected <= 0 {
			t.Errorf("Nothing deleted")
		}

		var entityList []models.FirstTableEntity
		db.Table(table).Where(where).Find(&entityList)
		if entityList != nil && len(entityList) > 0 {
			t.Errorf("Data not deleted")
		}
	})

	db.Rollback()
}

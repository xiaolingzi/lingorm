package tests

import (
	"testing"
	"time"

	"github.com/xiaolingzi/lingorm/tests/configs"
	"github.com/xiaolingzi/lingorm/tests/models"
)

func TestTransaction(t *testing.T) {

	t.Run("TestCommit", func(t *testing.T) {
		db := configs.GetDB()
		db.Begin()
		entity := models.FirstTableEntity{}
		entity.FirstName = "go name"
		entity.FirstNumber = 1001
		entity.FirstTime = time.Now()

		id, err := db.Insert(entity)
		db.Commit()
		if err != nil || id <= 0 {
			t.Errorf("Insert error")
		}

		db2 := configs.GetDB()
		table := models.FirstTableEntity{}.Table()
		var exists models.FirstTableEntity
		_, err = db2.Table(table).Where(table.ID.EQ(id)).OrderBy(table.FirstNumber.DESC()).First(&exists)

		if err != nil {
			t.Errorf("Find error")
		}
		if exists.FirstNumber != 1001 {
			t.Errorf("Commit failed")
		}

		db2.Delete(exists)
	})

	t.Run("TestRollback", func(t *testing.T) {
		db := configs.GetDB()
		db.Begin()
		entity := models.FirstTableEntity{}
		entity.FirstName = "go name"
		entity.FirstNumber = 1002
		entity.FirstTime = time.Now()

		id, err := db.Insert(entity)
		db.Rollback()
		if err != nil || id <= 0 {
			t.Errorf("Insert error")
		}

		db2 := configs.GetDB()
		table := models.FirstTableEntity{}.Table()
		var exists models.FirstTableEntity
		_, err = db2.Table(table).Where(table.ID.EQ(id)).OrderBy(table.FirstNumber.DESC()).First(&exists)

		if err != nil {
			t.Errorf("Find error")
		}
		if exists.FirstNumber == 1002 {
			t.Errorf("Rollback failed")
		}

	})
}

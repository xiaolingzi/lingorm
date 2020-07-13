package mysql

import (
	"strconv"
	"testing"
	"time"

	"github.com/xiaolingzi/lingorm/tests/configs"
	"github.com/xiaolingzi/lingorm/tests/models"
)

func TestQuery(t *testing.T) {
	db := configs.GetDB()
	db.Begin()
	var entityList []interface{}
	for i := 1; i < 3; i++ {
		entity := models.FirstTableEntity{}
		entity.FirstName = "go query name " + strconv.Itoa(i)
		entity.FirstNumber = 1000 + i
		entity.FirstTime = time.Now()
		entityList = append(entityList, entity)
	}

	affected, err := db.BatchInsert(entityList)
	if err != nil || affected != 2 {
		t.Errorf("Inert data error")
	}

	t.Run("TestFind", func(t *testing.T) {
		table := models.FirstTableEntity{}.Table()
		where := db.CreateWhere()
		where.And(table.FirstName.EQ("go query name 1"), table.FirstNumber.EQ(1001))
		where.OrAnd(table.FirstName.EQ("go query name 2"), table.FirstNumber.EQ(1002))

		orderBy := db.CreateOrderBy().By(table.FirstNumber.DESC())

		var list []models.FirstTableEntity
		result, err := db.Find(table, where, orderBy, &list)
		if err != nil {
			t.Errorf("Find error")
		}
		if len(list) <= 0 || list[0].FirstNumber != 1002 || result == nil {
			t.Errorf("Find result invalid")
		}
	})

	t.Run("TestFindPage", func(t *testing.T) {
		table := models.FirstTableEntity{}.Table()
		where := db.CreateWhere()
		where.And(table.FirstName.EQ("go query name 1"), table.FirstNumber.EQ(1001))
		where.OrAnd(table.FirstName.EQ("go query name 2"), table.FirstNumber.EQ(1002))

		orderBy := db.CreateOrderBy().By(table.FirstNumber.DESC())

		var list []models.FirstTableEntity
		result, err := db.FindPage(table, where, orderBy, 1, 1, &list)
		if err != nil {
			t.Errorf("FindPage error")
		}
		if len(list) <= 0 || list[0].FirstNumber != 1002 || result.TotalPages != 2 {
			t.Errorf("FindPage result invalid")
		}
	})

	t.Run("TestFindTop", func(t *testing.T) {
		table := models.FirstTableEntity{}.Table()
		where := db.CreateWhere()
		where.And(table.FirstName.EQ("go query name 1"), table.FirstNumber.EQ(1001))
		where.OrAnd(table.FirstName.EQ("go query name 2"), table.FirstNumber.EQ(1002))

		orderBy := db.CreateOrderBy().By(table.FirstNumber.DESC())

		var list []models.FirstTableEntity
		result, err := db.FindTop(table, 1, where, orderBy, &list)
		if err != nil {
			t.Errorf("FindTop error")
		}
		if len(list) != 1 || list[0].FirstNumber != 1002 || result == nil {
			t.Errorf("FindTop result invalid")
		}
	})

	t.Run("TestFirst", func(t *testing.T) {
		table := models.FirstTableEntity{}.Table()
		where := db.CreateWhere()
		where.And(table.FirstName.EQ("go query name 1"), table.FirstNumber.EQ(1001))
		where.OrAnd(table.FirstName.EQ("go query name 2"), table.FirstNumber.EQ(1002))

		orderBy := db.CreateOrderBy().By(table.FirstNumber.DESC())

		var entity models.FirstTableEntity
		result, err := db.First(table, where, orderBy, &entity)
		if err != nil {
			t.Errorf("First error")
		}

		if (entity == models.FirstTableEntity{} || entity.FirstNumber != 1002 || result == nil) {
			t.Errorf("First result invalid")
		}
	})

	t.Run("TestFindCount", func(t *testing.T) {
		table := models.FirstTableEntity{}.Table()
		where := db.CreateWhere()
		where.And(table.FirstName.EQ("go query name 1"), table.FirstNumber.EQ(1001))
		where.OrAnd(table.FirstName.EQ("go query name 2"), table.FirstNumber.EQ(1002))

		result, err := db.FindCount(table, where)
		if err != nil {
			t.Errorf("FindCount error")
		}

		if result != 2 {
			t.Errorf("FindCount result invalid")
		}
	})

	db.Rollback()
}

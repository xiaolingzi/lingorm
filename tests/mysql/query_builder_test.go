package mysql

import (
	"strconv"
	"testing"
	"time"

	"github.com/xiaolingzi/lingorm/tests/configs"
	"github.com/xiaolingzi/lingorm/tests/models"
)

type Result struct {
	FirstName   string    `json:"firstName" column:"first_name"`
	FirstNumber int       `json:"firstNumber" column:"first_number"`
	Num         time.Time `json:"firstTime" column:"first_time"`
	SecondName  string    `json:"secondName" column:"second_name"`
}

func TestQueryBuilder(t *testing.T) {
	db := configs.GetDB()
	db.Begin()
	var firstList []interface{}
	for i := 1; i < 3; i++ {
		entity := models.FirstTableEntity{}
		entity.FirstName = "first name " + strconv.Itoa(i)
		entity.FirstNumber = 1000 + i
		entity.FirstTime = time.Now()
		firstList = append(firstList, entity)
	}

	affected, err := db.BatchInsert(firstList)
	if err != nil || affected != 2 {
		t.Errorf("Inert data error")
	}

	var secondList []interface{}
	for i := 1; i < 3; i++ {
		entity := models.SecondTableEntity{}
		entity.SecondName = "second name " + strconv.Itoa(i)
		entity.SecondNumber = 1000 + i
		entity.SecondTime = time.Now()
		secondList = append(secondList, entity)
	}

	affected, err = db.BatchInsert(secondList)
	if err != nil || affected != 2 {
		t.Errorf("Inert data error")
	}

	t.Run("TestFind", func(t *testing.T) {
		firstTable := models.FirstTableEntity{}.Table()
		where := db.CreateWhere()
		where.And(firstTable.FirstName.EQ("first name 1"), firstTable.FirstNumber.EQ(1001))
		where.OrAnd(firstTable.FirstName.EQ("first name 2"), firstTable.FirstNumber.EQ(1002))

		orderBy := db.CreateOrderBy().By(firstTable.FirstNumber.DESC())

		secondTable := models.SecondTableEntity{}.Table()

		builder := db.QueryBuilder()
		builder = builder.Select(firstTable.FirstNumber, firstTable.FirstName.Max().Alias("first_name"), firstTable.ID.Count().Alias("num"), secondTable.SecondName.F("MAX").Alias("second_name")).
			From(firstTable).
			LeftJoin(secondTable, firstTable.FirstNumber.EQ(secondTable.SecondNumber)).
			GroupBy(firstTable.FirstNumber).
			Where(where).
			OrderBy(orderBy).
			Limit(1)

		var list []Result
		result, err := builder.Find(&list)

		if err != nil {
			t.Errorf("Find error")
		}
		if len(list) != 1 || list[0].FirstNumber != 1002 || result == nil {
			t.Errorf("Find result invalid")
		}
	})

	t.Run("TestFindPage", func(t *testing.T) {
		firstTable := models.FirstTableEntity{}.Table()
		secondTable := models.SecondTableEntity{}.Table()
		where := db.CreateWhere()
		where.And(firstTable.FirstName.EQ("first name 1"), firstTable.FirstNumber.EQ(1001))
		where.OrAnd(firstTable.FirstName.EQ("first name 2"), firstTable.FirstNumber.EQ(1002))

		builder := db.QueryBuilder()
		builder = builder.Select(firstTable.FirstNumber, firstTable.FirstName.I().Max().Alias("first_name"), firstTable.ID.I().Count().Alias("num"), secondTable.SecondName.I().F("MAX").Alias("second_name")).
			From(firstTable).
			RightJoin(secondTable, firstTable.FirstNumber.EQ(secondTable.SecondNumber)).
			GroupBy(firstTable.FirstNumber).
			Where(where).
			OrderBy(firstTable.FirstNumber.DESC())

		var list []Result
		result, err := builder.FindPage(1, 1, &list)
		if err != nil {
			t.Errorf("FindPage error")
		}
		if len(list) <= 0 || list[0].FirstNumber != 1002 || result.TotalPages != 2 {
			t.Errorf("FindPage result invalid")
		}
	})

	t.Run("TestFirst", func(t *testing.T) {
		firstTable := models.FirstTableEntity{}.Table()
		secondTable := models.SecondTableEntity{}.Table()
		where := db.CreateWhere()
		where.And(firstTable.FirstName.EQ("first name 1"), firstTable.FirstNumber.EQ(1001))
		where.OrAnd(firstTable.FirstName.EQ("first name 2"), firstTable.FirstNumber.EQ(1002))

		builder := db.QueryBuilder()
		builder = builder.Select(firstTable.FirstNumber, firstTable.FirstName.Max().Alias("first_name"), firstTable.ID.Count().Alias("num"), secondTable.SecondName.F("MAX").Alias("second_name")).
			From(firstTable).
			InnerJoin(secondTable, firstTable.FirstNumber.EQ(secondTable.SecondNumber)).
			GroupBy(firstTable.FirstNumber).
			Where(where).
			OrderBy(firstTable.FirstNumber.DESC())

		var entity Result
		result, err := builder.First(&entity)

		if err != nil {
			t.Errorf("First error")
		}

		if (entity == Result{} || entity.FirstNumber != 1002 || result == nil) {
			t.Errorf("First result invalid")
		}
	})

	t.Run("TestFindCount", func(t *testing.T) {
		firstTable := models.FirstTableEntity{}.Table()
		secondTable := models.SecondTableEntity{}.Table()
		where := db.CreateWhere()
		where.And(firstTable.FirstName.EQ("first name 1"), firstTable.FirstNumber.EQ(1001))
		where.OrAnd(firstTable.FirstName.EQ("first name 2"), firstTable.FirstNumber.EQ(1002))

		builder := db.QueryBuilder()
		builder = builder.Select(firstTable.FirstNumber, firstTable.FirstName.Max().Alias("first_name"), firstTable.ID.Count().Alias("num"), secondTable.SecondName.F("MAX").Alias("second_name")).
			From(firstTable).
			InnerJoin(secondTable, firstTable.FirstNumber.EQ(secondTable.SecondNumber)).
			GroupBy(firstTable.FirstNumber).
			Where(where)

		result, err := builder.FindCount()

		if err != nil {
			t.Errorf("FindCount error")
		}

		if result != 2 {
			t.Errorf("FindCount result invalid")
		}
	})

	db.Rollback()
}

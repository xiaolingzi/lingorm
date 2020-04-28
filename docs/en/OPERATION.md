# Data CUD

## Insert

### Insert Single row

``` go
db := lingorm.DB("testdb1")
myCompany := company.CompanyEntity{}
myCompany.CompanyName = "go test 1"
myCompany.CreatedAt = time.Now()
myCompany.UpdatedAt = time.Now()
id, err := db.Insert(myCompany)
```

### Batch Insert

``` go
db := lingorm.DB("testdb1")

myCompany1 := company.CompanyEntity{}
myCompany1.CompanyName = "go test 2"
myCompany1.CreatedAt = time.Now()
myCompany1.UpdatedAt = time.Now()

myCompany2 := company.CompanyEntity{}
myCompany2.CompanyName = "go test 3"
myCompany2.CreatedAt = time.Now()
myCompany2.UpdatedAt = time.Now()

list := []interface{}{myCompany1, myCompany2}
affected, err := db.BatchInsert(list)
```

## Update

### Update Single row

``` go
affected, err := db.Update(myCompany)
```

### Batch Update

``` go
affected, err := db.BatchUpdate(list)
```

### Update By

``` go
db := lingorm.DB("testdb1")
table := company.CompanyEntity{}.Table()
where := db.CreateWhere().Or(table.ID.EQ(37), table.ID.EQ(8))
var params = []interface{}{
    table.CompanyName.EQ("new company name"),
    table.ShortName.EQ("abc"),
}
affected, err := db.UpdateBy(table, params, where)
```

## Delete

### Delete Single row

``` go
affected, err := db.Delete(myCompany)
```

### Delete By

``` go
table := company.CompanyEntity{}.Table()
where := db.CreateWhere().Or(table.ID.EQ(38), table.ID.EQ(39))
affected, err := db.DeleteBy(table, where)
```

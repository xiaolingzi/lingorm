# 数据的增删改

## 插入

### 单条插入

``` go
db := lingorm.DB("testdb1")
myCompany := company.CompanyEntity{}
myCompany.CompanyName = "go test 1"
myCompany.CreatedAt = time.Now()
myCompany.UpdatedAt = time.Now()
id, err := db.Insert(myCompany)
```

### 批量插入

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

## 更新

### 单条更新

``` go
affected, err := db.Update(myCompany)
```

### 批量更新

``` go
affected, err := db.BatchUpdate(list)
```

### 条件更新

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

## 删除

### 单条删除

``` go
affected, err := db.Delete(myCompany)
```

### 条件删除

``` go
table := company.CompanyEntity{}.Table()
where := db.CreateWhere().Or(table.ID.EQ(38), table.ID.EQ(39))
affected, err := db.DeleteBy(table, where)
```

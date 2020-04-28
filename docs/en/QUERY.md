# Data Query

## Where Clause

Example：

``` go
db := lingorm.DB("testdb1")
table := company.CompanyEntity{}.Table()
where := db.CreateWhere()
where.Or(table.ID.EQ(38), table.ID.EQ(39))
where.And(table.CompanyName.LIKE("name"))
fmt.Println(where.CurrentSQL()) // print (t1.id = :p1 OR t1.id = :p2) AND t1.company_name like :p3
```

With 'CreateWhere', you can create a WHERE clause object. There are four member methods in it. There are 'Or', 'And', 'GetOr' and 'GetAnd'. 'Or' and 'And' return the where object, while 'GetOr' and 'GetAnd' return the where clause string.

``` go
db := lingorm.DB("testdb1")
table := company.CompanyEntity{}.Table()
where := db.CreateWhere()
where.Or(table.ID.EQ(38), table.ID.EQ(39))
where.And(table.CompanyName.LIKE("%name%"), where.GetOr(table.ShortName.EQ("name"), table.ShortName.LIKE("a%")))
fmt.Println(where.CurrentSQL()) // print (t1.id = :p1 OR t1.id = :p2) AND t1.company_name like :p5 AND (t1.short_name = :p3 OR t1.short_name like :p4)
```

The available operators are as follows:

GT | Greater than. Example: table.ID.GE(39)
GE | Greater and equal than.
LT| Less than.
LE | Less and equal than.
EQ | Equal.
NEQ | Not Equal.
LIKE
IN
NIN | Not in.

## Simple Query

1) Use the 'Table' function

``` go
db := lingorm.DB("testdb1")
table := company.CompanyEntity{}.Table()
result, err := db.Table(table).Select(table.ID, table.CompanyName).Where(table.IsDeleted.EQ(0), table.ID.GE(5)).OrderBy(table.ID.DESC()).Find()
```

This 'Where' method can accept the instance created by 'CreateWhere' as argument.
In addition to the 'Find' method. All the available functions are as follows:

```go
First() (interface{}, error) // return the first row
Find() (interface{}, error) // return all rows
FindPage(pageIndex int, pageSize int) (common.PageResult, error) // return page result
FindCount() (int, error) // return the number of rows
```

2）直接查询

``` go
db := lingorm.DB("testdb1")
table := company.CompanyEntity{}.Table()
where := db.CreateWhere().And(table.IsDeleted.EQ(0), table.ID.GE(5))
orderBy := db.CreateOderBy().By(table.ID.DESC(), table.CreatedAt.ASC)
result, err := db.Find(table, where, orderBy)
```

Other functions available:

```go
Find(table interface{}, args ...interface{}) (interface{}, error)
FindTop(table interface{}, top int, args ...interface{}) (interface{}, error)
First(table interface{}, args ...interface{}) (interface{}, error)
FindPage(table interface{}, pageIndex int, pageSize int, args ...interface{}) (common.PageResult, error)
```

'where' and 'order by' instances are supported for args

## Query Builder

Query Builder support for multi-table joint queries. For example:

``` go
db := lingorm.DB("testdb1")
companyTable := company.CompanyEntity{}.Table()
departmentTable := company.DepartmentEntity{}.Table()
builder := db.QueryBuilder()
builder.Select(departmentTable.CompanyID, departmentTable.ID.Count().Alias("Num"), companyTable.CompanyName.Max().Alias("companyName")).
    From(departmentTable).
    LeftJoin(companyTable, departmentTable.CompanyID.EQ(companyTable.ID)).
    Where(departmentTable.IsDeleted.EQ(0)).
    GroupBy(departmentTable.CompanyID).
    OrderBy(departmentTable.CompanyID.ASC())
result, err := builder.Find()
```

The type of 'Find' function's return value is '[]map[string]string'. If you want to map the result to a struct, you can do it like this:

``` go
type DepartmentResult struct {
    CompanyID   int    `column:"company_id"`
    CompanyName string `column:"company_name"`
    Num         int    `column:"num"`
}

result, err := builder.Find(DepartmentResult{})
```

Also, the where object created by 'CreateWhere' can be used as the 'Where' function's argument.

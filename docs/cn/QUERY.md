# 数据查询

## 查询条件构建

我们先来看一个示例：

``` go
db := lingorm.DB("testdb1")
table := company.CompanyEntity{}.Table()
where := db.CreateWhere()
where.Or(table.ID.EQ(38), table.ID.EQ(39))
where.And(table.CompanyName.LIKE("name"))
fmt.Println(where.CurrentSQL()) // 输出 (t1.id = :p1 OR t1.id = :p2) AND t1.company_name like :p3
```

通过CreateWhere方法可以构建一个where对象，where对象提供Or和And方法，这两个方法就分别对应数据库中的or和and。而且这两个方法都是接受多个参数。
除此之外，where对象还提供了GetOr和GetAnd方法，参数跟Or和And一样，但是他们是返回当前的条件字符串，用于一些复杂的条件拼接。例如：

``` go
db := lingorm.DB("testdb1")
table := company.CompanyEntity{}.Table()
where := db.CreateWhere()
where.Or(table.ID.EQ(38), table.ID.EQ(39))
where.And(table.CompanyName.LIKE("%name%"), where.GetOr(table.ShortName.EQ("name"), table.ShortName.LIKE("a%")))
fmt.Println(where.CurrentSQL()) // 输出 (t1.id = :p1 OR t1.id = :p2) AND t1.company_name like :p5 AND (t1.short_name = :p3 OR t1.short_name like :p4)
```

条件中比较运算符有如下：
GT 大于。比如 table.ID.GE(39)
GE 大于等于
LT 小于
LE 小于等于
EQ 等于
NEQ 不等于
LIKE 模糊匹配
IN 对应sql中的in
NIN 对应sql中的not in

## 简单查询

1）通过Table方法进行链式调用

``` go
db := lingorm.DB("testdb1")
table := company.CompanyEntity{}.Table()
var result []company.CompanyEntity
_, err := db.Table(table).Select(table.ID, table.CompanyName).Where(table.IsDeleted.EQ(0), table.ID.GE(5)).OrderBy(table.ID.DESC()).Find(&result)

// 或者
//result, err := db.Table(table).Select(table.ID, table.CompanyName).Where(table.IsDeleted.EQ(0), table.ID.GE(5)).OrderBy(table.ID.DESC()).Find()
```

其中Where方法可以还可以接收CreateWhere创建的where对象以支持复杂的查询条件。Find方法返回的是一个列表，已经映射到模型中，可以进行类型转换，或者直接将结果变量指针传给Find方法。
除了Find方法之外，还有以下几个方法：

``` go
Find(slicePtr ...interface{}) (interface{}, error) // 返回符合条件的第一条数据
FindPage(pageIndex int, pageSize int, slicePtr ...interface{}) (common.PageResult, error) // 返回符合条件的数据列表
First(structPtr ...interface{}) (interface{}, error) //返回分页数据
FindCount() (int, error) // 返回数量
```

2）直接查询

``` go
db := lingorm.DB("testdb1")
table := company.CompanyEntity{}.Table()
where := db.CreateWhere().And(table.IsDeleted.EQ(0), table.ID.GE(5))
orderBy := db.CreateOderBy().By(table.ID.DESC(), table.CreatedAt.ASC)
var result []company.CompanyEntity
_, err := db.Find(table, where, orderBy, &result)
// 或者
// result, err := db.Find(table, where, orderBy)
```

类似的也有以下几个方法

``` go
Find(table interface{}, where interface{}, orderBy interface{}, slicePtr ...interface{}) (interface{}, error)
FindTop(table interface{}, top int, where interface{}, orderBy interface{}, slicePtr ...interface{}) (interface{}, error)
First(table interface{}, where interface{}, orderBy interface{}, structPtr ...interface{}) (interface{}, error)
FindPage(table interface{}, where interface{}, orderBy interface{}, pageIndex int, pageSize int, slicePtr ...interface{}) (common.PageResult, error)
```

## 复杂查询

复杂查询支持多表联合查询。示例：

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

其中Find方法默认返回的是[]map[string]string类型，如果需要映射到结构体，则将定义的结构体实例传递给Find方法，比如

``` go
type DepartmentResult struct {
    CompanyID   int    `column:"company_id"`
    CompanyName string `column:"company_name"`
    Num         int    `column:"num"`
}

var result []DepartmentResult
_, err := builder.Find(&result)
```

同样，Where方法也支持通过CreateWhere创建的where对象，以支持复杂的查询条件。

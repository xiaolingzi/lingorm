# 原生SQL支持

示例：

``` go
db := lingorm.DB("testdb1")
native := db.NativeQuery()
sql := "select * from company where id =:id"
params := make(map[string]interface{})
params["id"] = 1
var result []company.CompanyEntity
_, err := native.Find(sql, params, &result)

// 或者
// result, err := native.Find(sql, params)
```

如果需要将结果映射到结构体，就如下：

``` go
result, err := native.Find(sql, params, company.CompanyEntity{})
```

*需要注意的是，这里的参数化查询采用的命名参数的方式而不是问号。

所有支持的方法如下：

```go
Execute(sql string, params map[string]interface{}) (int, int, error) // 执行增、删、改时使用，分别返回影响条数、最后ID和错误
Find(sql string, params map[string]interface{}, slicePtr ...interface{}) (interface{}, error) // 查询列表时使用
FindPage(pageIndex int, pageSize int, sql string, params map[string]interface{}, slicePtr ...interface{}) (common.PageResult, error) // 查询分页数据时使用
FindCount(sql string, params map[string]interface{}) (int, error) // 查询数量时使用
First(sql string, params map[string]interface{}, structPtr ...interface{}) (interface{}, error) // 返回符合条件的第一条
```

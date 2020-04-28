# Native SQL

For example:

``` go
db := lingorm.DB("testdb1")
native := db.NativeQuery()
sql := "select * from company where id =:id"
params := make(map[string]interface{})
params["id"] = 1
result, err := native.Find(sql, params)
```

If you want to map the results to a structure object, you can do it like this:

``` go
result, err := native.Find(sql, params, company.CompanyEntity{})
```

*Note that you should use the named parameters in the sql instead of a question mark.

All supported methods are as follows:

```go
Execute(sql string, params map[string]interface{}) (int, int, error) // execute the insert, update and delete sql， return the affected rows, the last inserted id and error.
Find(sql string, params map[string]interface{}, entity ...interface{}) (interface{}, error) // return all rows
FindPage(pageIndex int, pageSize int, sql string, params map[string]interface{}, entity ...interface{}) (common.PageResult, error) // return page result
FindCount(sql string, params map[string]interface{}) (int, error) // return the number of rows
First(sql string, params map[string]interface{}, entity ...interface{}) (interface{}, error) // return the first row
```

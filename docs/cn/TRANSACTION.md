# 事务的使用

示例如下：

```go
db := lingorm.DB("testdb1")
err := db.Begin()
if err != nil {
    // do something here...
}
// do something here...
err = db.Commit()
// err = db.Rollback()
```

需要注意的是，同一个事务必须在同一个lingorm.DB("testdb1")创建的实例下执行以上三个方法。

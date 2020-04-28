# Transaction

Here's an example:

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

It is important to note that the same transaction must be in the same db connection. The three methods are executed under an instance created by the DB ( "testdb1" ) .

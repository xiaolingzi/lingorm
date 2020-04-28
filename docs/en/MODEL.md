# Model

## Examples and Instructions

``` go
package models

import (
    "github.com/xiaolingzi/lingorm/model"
    "time"
)

//CompanyEntity entity
type CompanyEntity struct {
    CompanyName         string    `json:"companyName" comlumn:"company_name"`
    CreatedAt           time.Time `json:"createdAt" comlumn:"created_at"`
    DeletedAt           time.Time `json:"deletedAt" comlumn:"deleted_at"`
    ID                  int       `json:"id" comlumn:"id" primary_key:"true" auto_increment:"true"`
    IsDeleted           int       `json:"isDeleted" comlumn:"is_deleted"`
    Logo                string    `json:"logo" comlumn:"logo"`
    ShortName           string    `json:"shortName" comlumn:"short_name"`
    UpdatedAt           time.Time `json:"updatedAt" comlumn:"updated_at"`
}

//CompanyTable table
type CompanyTable struct {
    TTDatabaseName      string
    TTTableName         string
    TTAlias             string
    CompanyName         model.Field
    CreatedAt           model.Field
    DeletedAt           model.Field
    ID                  model.Field
    IsDeleted           model.Field
    Logo                model.Field
    ShortName           model.Field
    UpdatedAt           model.Field
}

//Table return table
func (e CompanyEntity) Table() CompanyTable {
    return model.Table(CompanyTable{}, CompanyEntity{}, "company", "mydb").(CompanyTable)
}
```

As shown in the code above.

First, you have to import the model package, it is "github.com/xxx/lingorm/model".

Second, define a struct type that map to a database table. Like the 'CompanyEntity' above.
1）json. It is the property name after json serialization, It's mainly up to the serialization plug-ins used.
2）column. It is the name of the table fields.
3）primary_key. It indicates whether the field is a primary key, true or 1 means yes, and no setting or setting another value means No.
4）auto_increment. It indicates whether the field is self-added, true or 1 means yes, and no setting or setting another value means No.

At last, Define a struct type for table and a 'Table' function which return a table instance. You may be confused about this. As for most orm frameworks, define the second part is enough. But it does improve our efficiency. For example:

```go
where.Or(table.ID.EQ(38), table.ID.EQ(39))
```

The 'company' and 'mydb' in 'model.Table(CompanyTable{}, CompanyEntity{}, "company", "mydb")' are the table name and database name. And the database name can be empty if you don't need cross-database queries.

## Model Generation

As the code above shows, it's still a bit complicated to write. So it is highly recommended to use code generation tools to generate code for the model. There is an example in the tools diretory, and you can change it to meet your needs.

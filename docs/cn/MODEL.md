# 模型的定义

## 示例及说明

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

如果上面代码所示，模型的代码包含四个部分

第一部分是包的导入那块，需要导入 github.com/xxx/lingorm/model 包。

第二部分是CompanyEntity结构体那部分，这部分是与数据库表相对应定义的结构体，其中的标签字段说明如下：
1）json是json序列化后的属性名称，如果与属性名一样则不需要，主要根据所采用的序列化插件来决定；
2）column是对应的数据表中的字段名
3）primary_key表示该字段是否为主键，设置为true或者1则表示是，不设置或者设置其它值则表示否；
4）auto_increment表示该字段是否为自增字段，设置为true或者1则表示是，不设置或者设置其它值则表示否；

第三第四部分主要是定义orm对应的table。这里大家可能会有些困惑，对于大多数orm框架来说有第二部分就足够，那这多余的定义的作用在哪里？我们知道很多orm框架中不管是查询的字段还是查询的条件都是通过拼接字符串的方式进行，这样如果有错误的话需要到执行阶段才能发现，而这多余的定义就是为了解决这个问题，有了它在使用的过程中就不用再拼接字符串，具体在看了后面的例子之后就会明白了。

特别要注意的是，第四部分model.Table(CompanyTable{}, CompanyEntity{}, "company", "mydb")中的company和mydb分别是指表名和数据库名。也就是说表名和数据库名字是通过显示指定的，这样模型结构体的名称理论可以按你想要的方式进行取名，指定数据库名的话则可以支持在同一台服务器上进行跨库查询，不需要的话传空字符串即可。

## 模型的生成

由于多了第三第四部分的内容，整个模型的定义其实也不简单呢，所以建议最好还是通过生成工具去生成以提高效率。代码中了tools目录有一个生成模型的例子，大家可以按自己的需求进行相应的修改或者自己写一个。

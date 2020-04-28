# 数据库配置

数据库的连接配置是通过json文件进行配置，格式如下：

```json
{
    "testdb1":{
        "host":"192.168.0.22",
        "user":"db_user",
        "password":"password",
        "database":"dbname1",
        "charset":"utf8",
        "driver":"mysql"
    },
    "testdb2":{
        "host":"192.168.0.22",
        "user":"db_user",
        "password":"password",
        "database":"dbname2",
        "charset":"utf8",
        "driver":"mysql"
    },
    "testdb3":{
        "driver":"mysql",
        "database":"dbname3",
        "charset":"utf8",
        "user":"db_user",
        "password":"password",
        "servers":[
        {
            "host":"192.168.0.110",
            "mode":"w",
            "wweight":3
        },
        {
            "host":"192.168.0.111",
            "mode":"rw",
            "w_weight":1,
            "rweight":1
        },
        {
            "host":"192.168.0.112",
            "mode":"r",
            "rweight":3
        },
        {
            "host":"192.168.0.113",
            "user":"db_user",
            "password":"password",
            "database":"dbname113",
            "charset":"utf8",
            "mode":"r",
            "rweight":2
        }]
    }
}
```

1. 单台服务器
对于只有一台服务器的情况只需按testdb1和testdb2的方式进行配置即可。

2. 多台服务器集群
一读多写、多读多写的情况配置如testdb3那样，将host属性改为servers属性，servers里面配置每个数据库信息：
1）读写配置通过mode属性，r为只读，w为只写，rw为读写
2）数据库名称、用户密码、编码等配置可以继承父节点配置，也可以自定义
3）多台配置通过权重随机的算法来决定连接那台服务器，读权重属性为rweight、写权重的属性为wweight

3. 配置文件路径
数据库配置文件可以放在程序有可读权限的任何位置，然后通过环境变量LINGORM_CONFIG来设置其所在路径即可，例如：
`os.Setenv("LINGORM_CONFIG", "/your/path/database.json")`

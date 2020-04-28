# Database Configuration

The connection configuration of the database is configured via the json file. the json file format is as follows：

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

1. Single server
For the case of only one server, you can simply configure it in the way of 'testdb1' and 'testdb2'.

2. Cluster
The configuration of the cluster, like 'testdb3', and use the 'servers' property to config the servers instead of 'host'.
1）Configure the database readable and writeable through the 'mode' property, r is read-only, w is write-only, and rw is read and write
2）Database, user, passwords, encodings, and other configurations can inherit from the parent node. Or you can customize them in child node.
3）The read and write weights of the database are configured through the 'rweight' and 'wweight' properties.

3. Config the json file's path
The database config file can be placed anywhere the program has readable permissions, and then set the path it is in by the environment variable 'LINGORM_CONFIG', for example:
`os.Setenv("LINGORM_CONFIG", "/your/path/database.json")`

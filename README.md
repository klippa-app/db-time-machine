# db-time-machine

```

$ dbtm [FLAGS] 

FLAGS:

--config=<path>            defaults to `$PWD/.dbtm.yaml`
--<configField>=<value>    to override config values

dbtm list
dbtm switch
dbtm prune

-- 

```


```
dbname := "mydb"

if development {

    os.exec("dbtm --username='user' --password='password' ")
    dbname, err := dbtm.get({
        Username: "user",
        Password: "pass",
    })
    if err {}
} 

db := fmt.Printf("postgresql://user:pass@example.com/%s", dbname)

psql.connect(db)
```


```
prefix = "mydb"
migration = {
    directory = "./migrations"
    format = "!\d{15}_"
    command = "migrate up"
}
database = {
    dialect = "postgreql"
}
```

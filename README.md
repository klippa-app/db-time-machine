# DB Time Machine

DB Time Machine is an experimental tool for automatically creating, managing, and swapping
development databases based on migration file hashes.

Often at Klippa we have found ourselves in the situation where either as a team, or 
just by ourselves we're developing two features that both a database migration.

If we switch between such branches, while our code and database schema changes,
the migrations that have already been applied to our local database does not.

Usually the only recourse is to either rollback the migrations before switching, and
likely loose any data that was created since, or drop the database entirely and 
recreated it from scratch on the new branch.

This process looses us a lot of time and happens frequently enough that it's often
complained about.

DB Time Machine aims to fix this (at least for postgres databases), by tracking the hashes
of migration files and creating a new database automatically when new migrations are detected, 
cloning from the nearest parent hash, and finally running all the migrations.

If a database already exists for a hash, then that database is used as-is.

In theory this means that while developing you will never have to drop the database, rollback, 
run migrations or reseed the database ever again.

There are ofcourse some caveats to this, namely that data "loss" is nolonger a 
choice by the developer but an automatic side effect of migration hashes changing.
Hopefully though with enough history this should not be a common occurance.

# Usage

While DBTM is intended to be database and language agnostic, currently only
postgres databases are supported. However new database dialects can be
implemented using the driver interface exported from the `dbtm/db` package.

DBTM is written in golang it can be imported and used directly in
golang during your database connection setup:

```go
    // This config will be merged with the project `.dbtm.toml` config, this can be used
    // to override things like the connection URI if you have runtime specific settings.
    cfg := dbtm.Config{}

    // dbtm.TimeTravel will do all the necessary instantiation  
    dbName, err = dbtm.TimeTravel(dialect.Postgres(), cfg)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Time travelling to %s\n", dbName)

    // use the db name to build your URI and connect to the database as normal.
	dbURI := fmt.Printf("postgres://user:password@example.com/%s", dbName)

```

As for other langauges, DBTM can also be used as a standalone CLI tool. Once configured,
running the command `dbtm` will return the name of a database that can be used in 
your connection URI. The command also takes various arguments for overriding the project
config for runtime configuration. 
This allows dbtm to be used in a variety of ways, you can call it directly from your language
of choice with a call like `os.Exec('dbtm')`, or in your startup script and export the result
as an ENV variable. Whatever works best for you.

Unfortunately however, custom database dialects can only be implemented using the
go API. You could build a custom binary for your usecase, or make an pull request for
your dialect, but otherwise support for that is outside of the scope of our MVP.

# Configuration

Configuration is in toml, and by default dbtm will look for `.dbtm.toml` or `dbtm.toml`
in that order. both the golang library and the cli tool accept an override, which you can
use to specify any other `.toml` file of your choice.

```toml
# The Prefix is prepended to database names.
Prefix = "dbtm"

[Connection]
# The URI to use for connecting to the DB.
# `{}` will be substituted for the database name.
URI = "postgres://postgres@localhost/{}?sslmode=disable"
# The database/schema to use for the inital connection.
Database = "postgres" 


[Migration]
# The directory relative to the execution path where migrations are found.
Directory = "./migrations" 
# The regex format for matching migration file names. 
Format = "^[0-9]{4}_"
# The command to run all pending migrations. 
# Must accept the DB name as the first and only Parameter.
# Must be only the command or script. no args, spaces, or env variables, or logic.
Command = "./migrate-up"
```

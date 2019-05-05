# GOMysql [![Build Status](https://travis-ci.org/illuminasy/gomysql.svg?branch=master)](https://travis-ci.org/illuminasy/gomysql) [![Coverage Status](https://coveralls.io/repos/github/illuminasy/gomysql/badge.svg?branch=master)](https://coveralls.io/github/illuminasy/gomysql?branch=master) [![GoDoc](https://godoc.org/github.com/illuminasy/gomysql?status.svg)](https://godoc.org/github.com/illuminasy/gomysql) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/illuminasy/gomysql/blob/master/LICENSE.md)

GOMysql
A Simple logging package, prefixes logging level.
Also supports error logging to bugsnag.

Currently supported middlewares:
 
# Usage

Get the library:

    $ go get -v github.com/illuminasy/gomysql

Mysql DB Migration
```go

import github.com/illuminasy/gomysql/migrate

//MigrationDir default directory "./internal/database/migration/sql"
// copy all your migration sql there
// check github.com/golang-migrate/migrate for more details
config := gomysql.Config{
	DBHost: "localhost'
	DBPort: "3306"
	DBUser: "root"
	DBPass: "root"
	DBName: "testdb"
	MigrationDir "somedir" // use this to use custom directory
}

client := gomysql.GetClient(config)
err := client.Migrate()
if err != nil {
	fmt.Println(err)
}

```

Mysql 
```go
config := gomysql.Config{
	DBHost: "localhost'
	DBPort: "3306"
	DBUser: "root"
	DBPass: "root"
	DBName: "testdb"
	MigrationDir "somedir" // use this to use custom directory
}

// To check connection
status, err := c.ConnCheck()

```
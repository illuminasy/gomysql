package gomysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// MaxOpenConns maximum open connections
const MaxOpenConns = 10

// MaxIdleConns maximum idle connections
const MaxIdleConns = 20

// Config ...
type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPass       string
	DBName       string
	DBSSL        string
	DBTimeout    string
	DBCharset    string
	DBCollation  string
	MigrationDir string
}

// Client ...
type Client struct {
	config Config
}

// Datastore ...
type Datastore interface {
	getConfig() *Config
	ConnCheck() (bool, error)
	Query(query string, params []interface{}) (*sql.Rows, error)
	QueryRow(query string, params []interface{}) (*sql.Rows, error)
	Exec(query string, params []interface{}) (sql.Result, error)

	StmtExec(stmt *sql.Stmt, params []interface{}) (sql.Result, error)

	Migrate() error
}

// GetClient ...
func GetClient(c Config) *Client {
	return &Client{
		config: c,
	}
}

// ConnCheck check if connection exists
func (c Client) ConnCheck() (bool, error) {
	dbConn, err := getDbConn(c.config)
	defer dbConn.Close()
	if err != nil {
		return false, err
	}

	err = dbConn.Ping()
	if err == nil {
		return true, nil
	}

	return false, err
}

// Query queries multiple rows
func (c Client) Query(query string, params []interface{}) (*sql.Rows, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	return dbConn.Query(query, params...)
}

// QueryRow Execute query and returns single row
func (c Client) QueryRow(query string, params []interface{}) (*sql.Row, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	return dbConn.QueryRow(query, params...), nil
}

// Exec ...
func (c Client) Exec(query string, params []interface{}) (sql.Result, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	return dbConn.Exec(query, params...)
}

// StmtExec Execute query
func (c Client) StmtExec(stmt *sql.Stmt, params []interface{}) (sql.Result, error) {
	return stmt.Exec(params...)
}

func getDbConn(config Config) (*sql.DB, error) {
	var db *sql.DB

	var psqlSource = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s%s",
		config.DBUser,
		config.DBPass,
		config.DBHost,
		config.DBPort,
		config.DBName,
		buildParameters(config),
	)

	db, err := sql.Open("postgres", psqlSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)

	return db, nil
}

func buildParameters(config Config) string {
	var dbParameters = "?"

	if len(config.DBSSL) > 0 {
		dbParameters += "tls=true&"
	}

	if len(config.DBTimeout) > 0 {
		dbParameters += fmt.Sprintf("timeout=%s&", config.DBTimeout)
	}

	if len(config.DBCharset) > 0 {
		dbParameters += fmt.Sprintf("charset=%s&", config.DBCharset)
	}

	if len(config.DBCollation) > 0 {
		dbParameters += fmt.Sprintf("collation=%s&", config.DBCollation)
	}

	if dbParameters == "?" {
		dbParameters = ""
	} else if dbParameters[len(dbParameters)-1:] == "&" {
		dbParameters = dbParameters[0 : len(dbParameters)-1]
	}

	return dbParameters
}

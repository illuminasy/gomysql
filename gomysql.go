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
	ParaseTime   string
}

// Client ...
type Client struct {
	config Config
}

// Datastore ...
type Datastore interface {
	getConfig() *Config
	ConnCheck() (bool, error)
	GetStats()

	Query(query string, params []interface{}) (*sql.Rows, error)
	QueryRow(query string, params []interface{}) (*sql.Row, error)
	Exec(query string, params []interface{}) (sql.Result, error)

	Prepare(query string) (*sql.Stmt, error)
	StmtQuery(stmt *sql.Stmt, params []interface{}) (*sql.Rows, error)
	StmtQueryRow(stmt *sql.Stmt, params []interface{}) *sql.Row
	StmtExec(stmt *sql.Stmt, params []interface{}) (sql.Result, error)

	BeginTx() (*sql.Tx, error)
	TxPrepare(tx *sql.Tx, query string) (*sql.Stmt, error)
	TxQuery(tx *sql.Tx, query string, params []interface{}) (*sql.Rows, error)
	TxQueryRow(tx *sql.Tx, query string, params []interface{}) *sql.Row
	TxExec(tx *sql.Tx, query string, params []interface{}) (sql.Result, error)
	TxCommit(tx *sql.Tx) error
	TxRollback(tx *sql.Tx) error

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

// GetStats return database stats
func (c Client) GetStats() (sql.DBStats, error) {
	dbConn, err := getDbConn(c.config)
	defer dbConn.Close()
	if err != nil {
		return sql.DBStats{}, err
	}

	return dbConn.Stats(), nil
}

// Query queries multiple rows
func (c Client) Query(query string, params []interface{}) (*sql.Rows, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}

	return dbConn.Query(query, params...)
}

// QueryRow Execute query and returns single row
func (c Client) QueryRow(query string, params []interface{}) (*sql.Row, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}

	return dbConn.QueryRow(query, params...), nil
}

// Exec ...
func (c Client) Exec(query string, params []interface{}) (sql.Result, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}

	return dbConn.Exec(query, params...)
}

// Prepare Prepare query
func (c Client) Prepare(query string) (*sql.Stmt, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}
	return dbConn.Prepare(query)
}

// StmtQuery queries multiple rows
func (c Client) StmtQuery(stmt *sql.Stmt, params []interface{}) (*sql.Rows, error) {
	return stmt.Query(params...)
}

// StmtQueryRow Execute query and returns single row
func (c Client) StmtQueryRow(stmt *sql.Stmt, params []interface{}) *sql.Row {
	return stmt.QueryRow(params...)
}

// StmtExec Execute query
func (c Client) StmtExec(stmt *sql.Stmt, params []interface{}) (sql.Result, error) {
	return stmt.Exec(params...)
}

// TxBegin Start and return transaction
func (c Client) TxBegin() (*sql.Tx, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}

	return dbConn.Begin()
}

// TxQuery queries multiple rows
func (c Client) TxQuery(tx *sql.Tx, query string, params []interface{}) (*sql.Rows, error) {
	return tx.Query(query, params...)
}

// TxQueryRow Execute query and returns single row
func (c Client) TxQueryRow(tx *sql.Tx, query string, params []interface{}) *sql.Row {
	return tx.QueryRow(query, params...)
}

// TxExec Execute query
func (c Client) TxExec(tx *sql.Tx, query string, params []interface{}) (sql.Result, error) {
	return tx.Exec(query, params...)
}

// TxCommit commits transaction
func (c Client) TxCommit(tx *sql.Tx) error {
	return tx.Commit()
}

// TxRollback rollback transaction
func (c Client) TxRollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func getDbConn(config Config) (*sql.DB, error) {
	var db *sql.DB

	var mysqlSource = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s%s",
		config.DBUser,
		config.DBPass,
		config.DBHost,
		config.DBPort,
		config.DBName,
		buildParameters(config),
	)

	db, err := sql.Open("mysql", mysqlSource)
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

	if len(config.ParaseTime) > 0 {
		dbParameters += "parseTime=true&"
	}

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

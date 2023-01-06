package gomysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
)

// MaxOpenConns maximum open connections
const MaxOpenConns = 10

// MaxIdleConns maximum idle connections
const MaxIdleConns = 5

// MaxOpenConnsTime maximum connection timeout
const MaxOpenConnsTime = 30 * time.Second

// MaxIdleConnsTime maximum idle connection timeout
const MaxIdleConnsTime = 1 * time.Second

// Config ...
type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPass           string
	DBName           string
	DBSSL            string
	DBTimeout        string
	DBCharset        string
	DBCollation      string
	MigrationDir     string
	MigrationsTable  string
	ParaseTime       string
	MaxOpenConns     int
	MaxIdleConns     int
	MaxOpenConnsTime time.Duration
	MaxIdleConnsTime time.Duration
	NewrelicEnabled  bool
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

	QueryWithContext(ctx context.Context, query string, params []interface{}) (*sql.Rows, error)
	QueryRowWithContext(ctx context.Context, query string, params []interface{}) (*sql.Row, error)
	ExecWithContext(ctx context.Context, query string, params []interface{}) (sql.Result, error)

	Prepare(query string) (*sql.Stmt, error)
	StmtQuery(stmt *sql.Stmt, params []interface{}) (*sql.Rows, error)
	StmtQueryRow(stmt *sql.Stmt, params []interface{}) *sql.Row
	StmtExec(stmt *sql.Stmt, params []interface{}) (sql.Result, error)

	PrepareWithContext(ctx context.Context, query string) (*sql.Stmt, error)
	StmtQueryWithContext(ctx context.Context, stmt *sql.Stmt, params []interface{}) (*sql.Rows, error)
	StmtQueryRowWithContext(ctx context.Context, stmt *sql.Stmt, params []interface{}) *sql.Row
	StmtExecWithContext(ctx context.Context, stmt *sql.Stmt, params []interface{}) (sql.Result, error)

	BeginTx() (*sql.Tx, error)
	TxPrepare(tx *sql.Tx, query string) (*sql.Stmt, error)
	TxQuery(tx *sql.Tx, query string, params []interface{}) (*sql.Rows, error)
	TxQueryRow(tx *sql.Tx, query string, params []interface{}) *sql.Row
	TxExec(tx *sql.Tx, query string, params []interface{}) (sql.Result, error)
	TxCommit(tx *sql.Tx) error
	TxRollback(tx *sql.Tx) error

	BeginTxWithContext(ctx context.Context) (*sql.Tx, error)
	TxPrepareWithContext(ctx context.Context, tx *sql.Tx, query string) (*sql.Stmt, error)
	TxQueryWithContext(ctx context.Context, tx *sql.Tx, query string, params []interface{}) (*sql.Rows, error)
	TxQueryRowWithContext(ctx context.Context, tx *sql.Tx, query string, params []interface{}) *sql.Row
	TxExecWithContext(ctx context.Context, tx *sql.Tx, query string, params []interface{}) (sql.Result, error)

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
	if err != nil {
		return false, err
	}

	defer dbConn.Close()

	err = dbConn.Ping()
	if err == nil {
		return true, nil
	}

	return false, err
}

// ConnCheck check if connection exists
func (c Client) ConnCheckWithContext(ctx context.Context) (bool, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return false, err
	}

	defer dbConn.Close()

	err = dbConn.PingContext(ctx)
	if err == nil {
		return true, nil
	}

	return false, err
}

// GetStats return database stats
func (c Client) GetStats() (sql.DBStats, error) {
	dbConn, err := getDbConn(c.config)
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

// QueryWithContext queries multiple rows
func (c Client) QueryWithContext(ctx context.Context, query string, params []interface{}) (*sql.Rows, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}

	return dbConn.QueryContext(ctx, query, params...)
}

// QueryRow Execute query and returns single row
func (c Client) QueryRow(query string, params []interface{}) (*sql.Row, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}

	return dbConn.QueryRow(query, params...), nil
}

// QueryRowWithContext Execute query and returns single row
func (c Client) QueryRowWithContext(ctx context.Context, query string, params []interface{}) (*sql.Row, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}

	return dbConn.QueryRowContext(ctx, query, params...), nil
}

// Exec ...
func (c Client) Exec(query string, params []interface{}) (sql.Result, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}

	return dbConn.Exec(query, params...)
}

// ExecWithContext ...
func (c Client) ExecWithContext(ctx context.Context, query string, params []interface{}) (sql.Result, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}

	return dbConn.ExecContext(ctx, query, params...)
}

// Prepare Prepare query
func (c Client) Prepare(query string) (*sql.Stmt, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}

	return dbConn.Prepare(query)
}

// PrepareWithContext Prepare query
func (c Client) PrepareWithContext(ctx context.Context, query string) (*sql.Stmt, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}

	return dbConn.PrepareContext(ctx, query)
}

// StmtQuery queries multiple rows
func (c Client) StmtQuery(stmt *sql.Stmt, params []interface{}) (*sql.Rows, error) {
	return stmt.Query(params...)
}

// StmtQueryWithContext queries multiple rows
func (c Client) StmtQueryWithContext(ctx context.Context, stmt *sql.Stmt, params []interface{}) (*sql.Rows, error) {
	return stmt.QueryContext(ctx, params...)
}

// StmtQueryRow Execute query and returns single row
func (c Client) StmtQueryRow(stmt *sql.Stmt, params []interface{}) *sql.Row {
	return stmt.QueryRow(params...)
}

// StmtQueryRowWithContext Execute query and returns single row
func (c Client) StmtQueryRowWithContext(ctx context.Context, stmt *sql.Stmt, params []interface{}) *sql.Row {
	return stmt.QueryRowContext(ctx, params...)
}

// StmtExec Execute query
func (c Client) StmtExec(stmt *sql.Stmt, params []interface{}) (sql.Result, error) {
	return stmt.Exec(params...)
}

// StmtExecWithContext Execute query
func (c Client) StmtExecWithContext(ctx context.Context, stmt *sql.Stmt, params []interface{}) (sql.Result, error) {
	return stmt.ExecContext(ctx, params...)
}

// TxBegin Start and return transaction
func (c Client) TxBegin() (*sql.Tx, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}

	return dbConn.Begin()
}

// TxBeginWithContext Start and return transaction
func (c Client) TxBeginWithContext(ctx context.Context) (*sql.Tx, error) {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return nil, err
	}

	return dbConn.BeginTx(ctx, nil)
}

// TxQuery queries multiple rows
func (c Client) TxQuery(tx *sql.Tx, query string, params []interface{}) (*sql.Rows, error) {
	return tx.Query(query, params...)
}

// TxQueryWithContext queries multiple rows
func (c Client) TxQueryWithContext(ctx context.Context, tx *sql.Tx, query string, params []interface{}) (*sql.Rows, error) {
	return tx.QueryContext(ctx, query, params...)
}

// TxQueryRow Execute query and returns single row
func (c Client) TxQueryRow(tx *sql.Tx, query string, params []interface{}) *sql.Row {
	return tx.QueryRow(query, params...)
}

// TxQueryRowWithContext Execute query and returns single row
func (c Client) TxQueryRowWithContext(ctx context.Context, tx *sql.Tx, query string, params []interface{}) *sql.Row {
	return tx.QueryRowContext(ctx, query, params...)
}

// TxExec Execute query
func (c Client) TxExec(tx *sql.Tx, query string, params []interface{}) (sql.Result, error) {
	return tx.Exec(query, params...)
}

// TxExecWithContext Execute query
func (c Client) TxExecWithContext(ctx context.Context, tx *sql.Tx, query string, params []interface{}) (sql.Result, error) {
	return tx.ExecContext(ctx, query, params...)
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

	driverName := "mysql"

	if config.NewrelicEnabled {
		driverName = "nrmysql"
	}

	db, err := sql.Open(driverName, mysqlSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	moc := MaxOpenConns

	if config.MaxOpenConns > 0 {
		moc = config.MaxOpenConns
	}

	mic := MaxIdleConns

	if config.MaxIdleConns > 0 {
		mic = config.MaxIdleConns
	}

	moct := MaxOpenConnsTime

	if config.MaxOpenConnsTime > 0 {
		moct = config.MaxOpenConnsTime
	}

	mict := MaxIdleConnsTime

	if config.MaxIdleConnsTime > 0 {
		mict = config.MaxIdleConnsTime
	}

	db.SetMaxOpenConns(moc)
	db.SetMaxIdleConns(mic)
	db.SetConnMaxLifetime(moct)
	db.SetConnMaxIdleTime(mict)

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

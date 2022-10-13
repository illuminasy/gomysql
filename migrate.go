package gomysql

import (
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/file"
)

// DefaultMigrationDir Default directory for migration sql files
const DefaultMigrationDir = "./internal/database/migration/sql"

// Migrate start db migration up
func (c Client) Migrate() error {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	var mDir = DefaultMigrationDir
	if len(c.config.MigrationDir) > 0 {
		mDir = c.config.MigrationDir
	}

	var migrationDir = flag.String("migration.files", mDir, "Directory where the migration files are located?")

	// Run migrations
	driver, err := mysql.WithInstance(dbConn, &mysql.Config{})
	if err != nil {
		return err
	}

	fsrc, err := (&file.File{}).Open(fmt.Sprintf("file://%s", *migrationDir))
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"file",
		fsrc,
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	// m.Force(1)
	// if err != nil {
	// 	log.Fatalf("force migration failed... %v", err)
	// }

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// CleanUp start db migration down
func (c Client) CleanUp() error {
	dbConn, err := getDbConn(c.config)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	var mDir = DefaultMigrationDir
	if len(c.config.MigrationDir) > 0 {
		mDir = c.config.MigrationDir
	}

	var migrationDir = flag.String("migration.files", mDir, "Directory where the migration files are located?")

	// Run migrations
	driver, err := mysql.WithInstance(dbConn, &mysql.Config{})
	if err != nil {
		return err
	}

	fsrc, err := (&file.File{}).Open(fmt.Sprintf("file://%s", *migrationDir))
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"file",
		fsrc,
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	// m.Force(1)
	// if err != nil {
	// 	log.Fatalf("force migration failed... %v", err)
	// }

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

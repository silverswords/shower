package sql

import "database/sql"

const (
	mysqlDropDatabase = iota
	mysqlDropTable
)

var DropSQLStrings = []string{
	`DROP DATABASE IF EXISTS %s `,
	`DROP TABLE IF EXISTS %s `,
}

// DropTableFirst create role table.
func DropTableFirst(db *sql.DB) error {
	_, err := db.Exec(DropSQLStrings[mysqlDropTable])
	return err
}

// DropDatabaseFirst create role table.
func DropDatabaseFirst(db *sql.DB) error {
	_, err := db.Exec(DropSQLStrings[mysqlDropDatabase])
	return err
}

/*
 * Revision History:
 *     Initial: 2020/10/05        Abserari
 */

package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/abserari/shower/utils/salt"
)

const (
	mysqlUserCreateDatabase = iota
	mysqlUserCreateTable
	mysqlUserInsert
	mysqlUserLogin
	mysqlUserModifyEmail
	mysqlUserModifyMobile
	mysqlUserGetPassword
	mysqlUserModifyPassword
	mysqlUserModifyActive
	mysqlUserGetIsActive
)

const (
	DBName    = "admin"
	TableName = "admin"
)

var (
	errInvalidMysql = errors.New("affected 0 rows")
	errLoginFailed  = errors.New("invalid username or password")

	adminSQLString = []string{
		fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s ;`, DBName),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
			admin_id    BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			name     	VARCHAR(512) UNIQUE NOT NULL DEFAULT ' ',
			password 	VARCHAR(512) NOT NULL DEFAULT ' ',
			mobile   	VARCHAR(32) UNIQUE DEFAULT NULL,
			email    	VARCHAR(128) UNIQUE DEFAULT NULL,
			active   	BOOLEAN DEFAULT TRUE,
			created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (admin_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`, DBName, TableName),
		fmt.Sprintf(`INSERT INTO %s.%s (name,password,active)  VALUES (?,?,?)`, DBName, TableName),
		fmt.Sprintf(`SELECT admin_id,password FROM %s.%s WHERE name = ? LOCK IN SHARE MODE`, DBName, TableName),
		fmt.Sprintf(`UPDATE %s.%s SET email=? WHERE admin_id = ? LIMIT 1`, DBName, TableName),
		fmt.Sprintf(`UPDATE %s.%s SET mobile=? WHERE admin_id = ? LIMIT 1`, DBName, TableName),
		fmt.Sprintf(`SELECT password FROM %s.%s WHERE admin_id = ?  LOCK IN SHARE MODE`, DBName, TableName),
		fmt.Sprintf(`UPDATE %s.%s SET password = ? WHERE admin_id = ? LIMIT 1`, DBName, TableName),
		fmt.Sprintf(`UPDATE %s.%s SET active = ? WHERE admin_id = ? LIMIT 1`, DBName, TableName),
		fmt.Sprintf(`SELECT active FROM %s.%s WHERE admin_id = ? LOCK IN SHARE MODE`, DBName, TableName),
	}
)

// CreateDatabase create admin table.
func CreateDatabase(db *sql.DB) error {
	_, err := db.Exec(adminSQLString[mysqlUserCreateDatabase])
	if err != nil {
		return err
	}

	return nil
}

// CreateTable create admin table.
func CreateTable(db *sql.DB, name, password *string) error {
	_, err := db.Exec(adminSQLString[mysqlUserCreateTable])
	if err != nil {
		return err
	}

	//
	err = CreateAdmin(db, name, password)
	if err != nil {
		// don't error when create admin user twice.
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil
		}
		return err
	}
	return nil
}

//CreateAdmin create an administrative user
func CreateAdmin(db *sql.DB, name, password *string) error {
	hash, err := salt.Generate(password)
	if err != nil {
		return err
	}

	result, err := db.Exec(adminSQLString[mysqlUserInsert], name, hash, true)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

//Login the administrative user logins
func Login(db *sql.DB, name, password *string) (uint32, error) {
	var (
		id  uint32
		pwd string
	)

	err := db.QueryRow(adminSQLString[mysqlUserLogin], name).Scan(&id, &pwd)
	if err != nil {
		return 0, err
	}

	if !salt.Compare([]byte(pwd), password) {
		return 0, errLoginFailed
	}

	return id, nil
}

// ModifyEmail the administrative user updates email
func ModifyEmail(db *sql.DB, id uint32, email *string) error {

	result, err := db.Exec(adminSQLString[mysqlUserModifyEmail], email, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

// ModifyMobile the administrative user updates mobile
func ModifyMobile(db *sql.DB, id uint32, mobile *string) error {

	result, err := db.Exec(adminSQLString[mysqlUserModifyMobile], mobile, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

// ModifyPassword the administrative user updates password
func ModifyPassword(db *sql.DB, id uint32, password, newPassword *string) error {
	var (
		pwd string
	)

	err := db.QueryRow(adminSQLString[mysqlUserGetPassword], id).Scan(&pwd)
	if err != nil {
		return err
	}

	if !salt.Compare([]byte(pwd), password) {
		return errLoginFailed
	}

	hash, err := salt.Generate(newPassword)
	if err != nil {
		return err
	}

	_, err = db.Exec(adminSQLString[mysqlUserModifyPassword], hash, id)

	return err
}

//ModifyAdminActive the administrative user updates active
func ModifyAdminActive(db *sql.DB, id uint32, active bool) error {
	result, err := db.Exec(adminSQLString[mysqlUserModifyActive], active, id)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil

}

//IsActive return user.Active and nil if query success.
func IsActive(db *sql.DB, id uint32) (bool, error) {
	var (
		isActive bool
	)

	db.QueryRow(adminSQLString[mysqlUserGetIsActive], id).Scan(&isActive)
	return isActive, nil
}

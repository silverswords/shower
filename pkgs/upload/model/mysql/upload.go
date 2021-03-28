/*
 * Revision History:
 *     Initial: 2020/1016       Abserari
 */

package mysql

import (
	"database/sql"
	"errors"
	"time"
)

const (
	mysqlFileCreateTable = iota
	mysqlFileInsert
	mysqlFileQueryByMD5
	mysqlDeleteByPath
)

var (
	//ErrNoRows -
	ErrNoRows        = errors.New("there is no such data in database")
	errInvalidInsert = errors.New("upload file: insert affected 0 rows")

	sqlString = []string{
		`CREATE TABLE IF NOT EXISTS files (
			 user_id 	INTEGER UNSIGNED NOT NULL,
			 md5 		VARCHAR(512) NOT NULL DEFAULT ' ',
			 path 		VARCHAR(512) NOT NULL DEFAULT ' ',
			 created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			 PRIMARY KEY (md5)
		 ) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO files(user_id,md5,path,created_at) VALUES (?,?,?,?)`,
		`SELECT path FROM files WHERE md5 = ? LOCK IN SHARE MODE`,
		`DELETE FROM files WHERE path = ? LIMIT 1`,
	}
)

// CreateTable create files table.
func CreateTable(db *sql.DB) error {
	_, err := db.Exec(sqlString[mysqlFileCreateTable])

	return err
}

// Insert insert a file
func Insert(db *sql.DB, userID uint32, path, md5 string) error {
	result, err := db.Exec(sqlString[mysqlFileInsert], userID, md5, path, time.Now())
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidInsert
	}

	return nil
}

// QueryByMD5 select by MD5
func QueryByMD5(db *sql.DB, md5 string) (string, error) {
	var (
		path string
	)

	err := db.QueryRow(sqlString[mysqlFileQueryByMD5], md5).Scan(&path)
	if err != nil {
		if err == sql.ErrNoRows {
			return path, ErrNoRows
		}
		return path, err
	}

	return path, nil
}

// Delete Clear the records of file so file could reupload.
func DeleteByPath(db *sql.DB, path string) error {
	_, err := db.Exec(sqlString[mysqlDeleteByPath], path)
	if err != nil {
		return errors.New("errDeleteMysql")
	}

	return nil
}

/*
 * Revision History:
 *     Initial: 2019/03/20       Yang ChengKai
 */

package mysql

import (
	"database/sql"
	"errors"
)

// Message -
type Message struct {
	Mobile string `db:"mobile"`
	Date   int64  `db:"date"`
	Code   string `db:"code"`
	Sign   string `db:"sign"`
}

const (
	mysqlMessageCreateTable = iota
	mysqlMessageInsert
	mysqlMessageGetDate
	mysqlMessageDelete
	mysqlMessageGetCode
	mysqlMessageUGetMobile
)

var (
	messageSQLString = []string{
		`CREATE TABLE IF NOT EXISTS message(
			mobile VARCHAR(32) UNIQUE NOT NULL,
			date  INT(11) DEFAULT 0,
			code VARCHAR(32) ,
			sign VARCHAR(32) UNIQUE NOT NULL
		)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO message(mobile,date,code,sign) VALUES (?,?,?,?)`,
		`SELECT date FROM message WHERE sign = ?`,
		`DELETE FROM message WHERE sign = ? LIMIT 1`,
		`SELECT code FROM message WHERE sign = ?`,
		`SELECT mobile FROM message WHERE sign = ?`,
	}
)

// CreateTable create message table.
func CreateTable(db *sql.DB) error {
	_, err := db.Exec(messageSQLString[mysqlMessageCreateTable])
	return err
}

// Insert Insert a new message.
func Insert(db *sql.DB, mobile string, date int64, code string, sign string) error {
	result, err := db.Exec(messageSQLString[mysqlMessageInsert], mobile, date, code, sign)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("errInvalidInsert")
	}

	return nil
}

// GetDate return message date(unixtime) and nil if no err,or (0,err).
func GetDate(db *sql.DB, sign string) (int64, error) {
	var unixtime int64

	err := db.QueryRow(messageSQLString[mysqlMessageGetDate], sign).Scan(&unixtime)
	if err != nil {
		return 0, errors.New("errQueryDate")
	}

	return unixtime, nil
}

// Delete Clear delete a  message.
func Delete(db *sql.DB, sign string) error {
	_, err := db.Exec(messageSQLString[mysqlMessageDelete], sign)
	if err != nil {
		return errors.New("errDeleteMysql")
	}

	return nil
}

// GetCode return message date and nil or "0"and err.
func GetCode(db *sql.DB, sign string) (string, error) {
	var code string

	err := db.QueryRow(messageSQLString[mysqlMessageGetCode], sign).Scan(&code)
	if err != nil {
		return "0", errors.New("errQueryCode")
	}

	return code, nil
}

//GetMobile return User's mobile like ID or "0" and err
func GetMobile(db *sql.DB, sign string) (string, error) {
	var mobile string

	err := db.QueryRow(messageSQLString[mysqlMessageUGetMobile], mobile).Scan(&mobile)
	if err != nil {
		return "0", err
	}

	return mobile, nil
}

//GetMessage return message
func GetMessage(db *sql.DB, sign string) *Message {
	var msg Message
	msg.Code, _ = GetCode(db, sign)
	msg.Date, _ = GetDate(db, sign)
	msg.Mobile, _ = GetMobile(db, sign)
	msg.Sign = sign

	return &msg
}

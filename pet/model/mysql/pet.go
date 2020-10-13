/*
 * Revision History:
 *     Initial: 2020/1018       Abserari
 */

package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Pet -
type Pet struct {
	PetID          uint64
	AdminID        uint64
	Name           string
	Category       string
	Avatar         string
	Birthday       time.Time
	MedicalCurrent string
	Hobbies        string
	Gender         string
}

const (
	mysqlPetCreateTable = iota
	mysqlPetInsert
	mysqlPetListPetByAdminID
	mysqlPetInfoByID
	mysqlPetUpdateNameByID
	mysqlPetUpdateCategoryByID
	mysqlPetUpdateAvatarByID
	mysqlPetUpdateBirthdayByID
	mysqlPetUpdateMedicalCurrentByID
	mysqlPetUpdateHobbiesByID
	mysqlPetUpdateGenderByID
	mysqlPetDeleteByID
)

var (
	errInvalidNoRowsAffected = errors.New("insert schedule:insert affected 0 rows")

	petSQLString = []string{
		`CREATE TABLE IF NOT EXISTS %s (
petID    BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
adminID    BIGINT UNSIGNED NOT NULL,
name        VARCHAR(512) NOT NULL DEFAULT ' ',
category VARCHAR(255)  NOT NULL DEFAULT ' ',
avatar   VARCHAR(512) NOT NULL DEFAULT ' ',
birthday   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
medicalCurrent     VARCHAR(512) NOT NULL DEFAULT ' ',
hobbies     VARCHAR(512) NOT NULL DEFAULT ' ',
gender VARCHAR(127) NOT NULL DEFAULT '中性',
PRIMARY KEY (petID)
)ENGINE=InnoDB AUTO_INCREMENT=1000000 DEFAULT CHARSET=utf8mb4;`,
		`INSERT INTO  %s (adminID,name,category,avatar,birthday,medicalCurrent,hobbies,gender) VALUES (?,?,?,?,?,?,?,?)`,
		`SELECT * FROM %s WHERE adminID = ? LOCK IN SHARE MODE`,
		`SELECT * FROM %s WHERE petID = ? LIMIT 1 LOCK IN SHARE MODE`,
		`UPDATE %s SET name=? WHERE petID = ? LIMIT 1`,
		`UPDATE %s SET category=? WHERE petID = ? LIMIT 1`,
		`UPDATE %s SET avatar=? WHERE petID = ? LIMIT 1`,
		`UPDATE %s SET birthday=? WHERE petID = ? LIMIT 1`,
		`UPDATE %s SET medicalCurrent=? WHERE petID = ? LIMIT 1`,
		`UPDATE %s SET hobbies=? WHERE petID = ? LIMIT 1`,
		`UPDATE %s SET gender=? WHERE petID = ? LIMIT 1`,
		`DELETE FROM %s WHERE petID = ? LIMIT 1`,
	}
)

// CreateTable -
func CreateTable(db *sql.DB, tableName string) error {
	sql := fmt.Sprintf(petSQLString[mysqlPetCreateTable], tableName)
	_, err := db.Exec(sql)
	return err
}

// InsertPet return  id
func InsertPet(db *sql.DB, tableName string, adminID uint64, name, category, avatar string, birthday time.Time, medical_current, hobbies, gender string) (int, error) {
	sql := fmt.Sprintf(petSQLString[mysqlPetInsert], tableName)
	result, err := db.Exec(sql, adminID, name, category, avatar, birthday, medical_current, hobbies, gender)
	if err != nil {
		return 0, err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return 0, errInvalidNoRowsAffected
	}

	petID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(petID), nil
}

// ListValidPetByUnixDate return schedule list which have valid date
func ListPetByAdminID(db *sql.DB, tableName string, adminID uint64) ([]*Pet, error) {
	var (
		pets []*Pet

		petID          uint64
		name           string
		category       string
		avatar         string
		birthday       time.Time
		medicalCurrent string
		hobbies        string
		gender         string
	)

	sql := fmt.Sprintf(petSQLString[mysqlPetListPetByAdminID], tableName)
	rows, err := db.Query(sql, adminID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&petID, &adminID, &name, &category, &avatar, &birthday, &medicalCurrent, &hobbies, &gender); err != nil {
			return nil, err
		}

		pet := &Pet{
			PetID:          petID,
			AdminID:        adminID,
			Name:           name,
			Category:       category,
			Avatar:         avatar,
			Birthday:       birthday,
			MedicalCurrent: medicalCurrent,
			Hobbies:        hobbies,
			Gender:         gender,
		}

		pets = append(pets, pet)
	}

	return pets, nil
}

// InfoByID squery by id
func InfoByID(db *sql.DB, tableName string, id uint64) (*Pet, error) {
	var pet Pet

	sql := fmt.Sprintf(petSQLString[mysqlPetInfoByID], tableName)
	rows, err := db.Query(sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&pet.PetID, &pet.AdminID, &pet.Name, &pet.Category, &pet.Avatar, &pet.Birthday, &pet.MedicalCurrent, &pet.Hobbies, &pet.Gender); err != nil {
			return nil, err
		}
	}

	return &pet, nil
}

// ModifyEmail the administrative user updates email
func ModifyName(db *sql.DB, id uint64, name string) error {
	result, err := db.Exec(petSQLString[mysqlPetUpdateNameByID], name, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidNoRowsAffected
	}

	return nil
}

// ModifyEmail the administrative user updates email
func ModifyCategory(db *sql.DB, id uint64, category string) error {
	result, err := db.Exec(petSQLString[mysqlPetUpdateCategoryByID], category, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidNoRowsAffected
	}

	return nil
} // ModifyAvatar the administrative user updates email
func ModifyAvatar(db *sql.DB, id uint64, avatar string) error {
	result, err := db.Exec(petSQLString[mysqlPetUpdateAvatarByID], avatar, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidNoRowsAffected
	}

	return nil
}

// ModifyEmail the administrative user updates email
func ModifyBirthday(db *sql.DB, id uint64, birthday time.Time) error {
	result, err := db.Exec(petSQLString[mysqlPetUpdateBirthdayByID], birthday, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidNoRowsAffected
	}

	return nil
}

// ModifyEmail the administrative user updates email
func ModifyMedicalCurrent(db *sql.DB, id uint64, MedicalCurrent string) error {
	result, err := db.Exec(petSQLString[mysqlPetUpdateMedicalCurrentByID], MedicalCurrent, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidNoRowsAffected
	}

	return nil
}

// ModifyEmail the administrative user updates email
func ModifyHobbies(db *sql.DB, id uint64, hobbies string) error {
	result, err := db.Exec(petSQLString[mysqlPetUpdateHobbiesByID], hobbies, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidNoRowsAffected
	}

	return nil
}

// ModifyEmail the administrative user updates email
func ModifyGender(db *sql.DB, id uint64, gender string) error {
	result, err := db.Exec(petSQLString[mysqlPetUpdateGenderByID], gender, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidNoRowsAffected
	}

	return nil
}

// DeleteByID delete by id
func DeleteByID(db *sql.DB, tableName string, id uint64) error {
	sql := fmt.Sprintf(petSQLString[mysqlPetDeleteByID], tableName)
	_, err := db.Exec(sql, id)
	return err
}

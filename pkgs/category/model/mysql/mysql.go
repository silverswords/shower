package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	mysqlCreateTable = iota
	mysqlInsert
	mysqlUpdateStatus
	mysqlUpdateName
	mysqlSelectByParentID
)

var (
	errInvaildInsert         = errors.New("insert comment: insert affected 0 rows")
	errInvalidChangeCategory = errors.New("change status: affected 0 rows")
	categorySQLFormatStr     = []string{
		`CREATE TABLE IF NOT EXISTS %s(
			categoryId INT(11) NOT NULL AUTO_INCREMENT COMMENT '类别id',
			parentId INT(11) DEFAULT NULL  COMMENT '父类别id',
			name VARCHAR(50) DEFAULT NULL COMMENT '类别名称',
			status TINYINT(1) DEFAULT '1' COMMENT '状态1-在售，2-废弃',
			createTime DATETIME DEFAULT current_timestamp COMMENT '创建时间',
			PRIMARY KEY (categoryId),
			INDEX(parentId)
		)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4`,
		`INSERT INTO %s(parentId,name) VALUES (?,?)`,
		`UPDATE %s SET status = ? WHERE categoryId = ? LIMIT 1`,
		`UPDATE %s SET name = ? WHERE categoryId = ? LIMIT 1`,
		`SELECT * FROM %s WHERE parentId = ?`,
	}
)

// Category -
type Category struct {
	CategoryID uint
	ParentID   uint //为0则是根目录
	Name       string
	Status     int8
	CreateTime time.Time
}

// CreateTable -
func CreateTable(db *sql.DB, tableName string) error {
	sql := fmt.Sprintf(categorySQLFormatStr[mysqlCreateTable], tableName)
	_, err := db.Exec(sql)
	return err
}

//InsertCategory 自动设定 id 和 status状态和 创建时间 -
func InsertCategory(db *sql.DB, tableName string, parentID uint, name string) (uint, error) {
	sql := fmt.Sprintf(categorySQLFormatStr[mysqlInsert], tableName)
	result, err := db.Exec(sql, parentID, name)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errInvaildInsert
	}

	categoryId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(categoryId), nil
}

//ChangeCategoryStatus 改变目录状态
func ChangeCategoryStatus(db *sql.DB, tableName string, category uint, status int8) error {
	sql := fmt.Sprintf(categorySQLFormatStr[mysqlUpdateStatus], tableName)
	result, err := db.Exec(sql, status, category)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidChangeCategory
	}

	return nil
}

//ChangeCategoryName 改变目录名称
func ChangeCategoryName(db *sql.DB, tableName string, category uint, name string) error {
	sql := fmt.Sprintf(categorySQLFormatStr[mysqlUpdateName], tableName)
	result, err := db.Exec(sql, name, category)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidChangeCategory
	}

	return nil
}

// LisitChirldrenByParentID -
func LisitChirldrenByParentID(db *sql.DB, tableName string, parentID uint) ([]*Category, error) {
	var (
		categoryID uint
		name       string
		status     int8
		creatTime  time.Time

		categorys []*Category
	)
	sql := fmt.Sprintf(categorySQLFormatStr[mysqlSelectByParentID], tableName)
	rows, err := db.Query(sql, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&categoryID, &parentID, &name, &status, &creatTime); err != nil {
			return nil, err
		}

		category := &Category{
			CategoryID: categoryID,
			ParentID:   parentID,
			Name:       name,
			Status:     status,
			CreateTime: creatTime,
		}
		categorys = append(categorys, category)
	}

	return categorys, nil
}

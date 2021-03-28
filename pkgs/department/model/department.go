/*
 * Revision History:
 *     Initial: 2020/10/05        Abserari
 */

package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	mysqlDepartmentCreateTable = iota
	mysqlDepartmentInsert
)

const (
	mysqlDepartmentMemberCreateTable = iota
	mysqlDepartmentMemberInsert
)

var (
	errInvalidMysql = errors.New("affected 0 rows")

	departmentSQLString = []string{
		`create table if not exists  %s.%s 
			(
				id                int auto_increment primary key,
				code              varchar(30)            null comment '编号',
				organization_code varchar(30)            null comment '组织编号',
				name              varchar(30)            null comment '名称',
				sort              int         default 0  null comment '排序',
				pcode             varchar(30) default '' null comment '上级编号',
				icon              varchar(20)            null comment '图标',
				create_time       varchar(20)            null comment '创建时间',
				path              text                   null comment '上级路径',
				constraint code
					unique (code)
			)comment '部门表' ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ;`,
		// insert one entry
		`INSERT INTO %s.%s (code,organization_code,name,sort,pcode,icon,create_time,path)  VALUES (?,?,?,?,?,?,?,?)`,
	}

	departmentMemberSQLString = []string{
		`create table if not exists %s.%s 
			(
				id                int auto_increment
					primary key,
				code              varchar(30) default '' null comment 'id',
				department_code   varchar(30) default '' null comment '部门id',
				organization_code varchar(30) default '' null comment '组织id',
				account_code      varchar(30) default '' null comment '成员id',
				join_time         varchar(255)           null comment '加入时间',
				is_principal      tinyint(1)             null comment '是否负责人',
				is_owner          tinyint(1)  default 0  null comment '拥有者',
				authorize         varchar(255)           null comment '角色',
				constraint code
					unique (code)
			)
				comment '部门-成员表' ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin comment '部门表';`,
		// insert one entry
		`INSERT INTO %s.%s (code,department_code,organization_code,account_code,join_time,is_principal,is_owner,authorize)  VALUES (?,?,?,?,?,?,?,?)`,
	}
)

type Department struct {
	Id               int64     `json:"id,omitempty"`
	Code             string    `json:"code,omitempty"`
	OrganizationCode string    `json:"organization_code,omitempty"`
	Name             string    `json:"name,omitempty"`
	Sort             string    `json:"sort,omitempty"`
	Pcode            string    `json:"pcode,omitempty"`
	Icon             string    `json:"icon,omitempty"`
	CreateTime       time.Time `json:"create_time,omitempty"`
	Path             string    `json:"path,omitempty"`
}

type DepartmentMember struct {
	Id               int64     `json:"id,omitempty"`
	Code             string    `json:"code,omitempty"`
	DepartmentCode   string    `json:"department_code,omitempty"`
	OrganizationCode string    `json:"organization_code,omitempty"`
	AccountCode      string    `json:"account_code,omitempty"`
	JoinTime         time.Time `json:"join_time,omitempty"`
	IsPrincipal      int8      `json:"is_principal,omitempty"`
	IsOwner          int8      `json:"is_owner,omitempty"`
	Authorize        string    `json:"authorize,omitempty"`
}

// CreateTable create
func CreateDepartmentTable(db *sql.DB, DBName, TableName string) error {
	_, err := db.Exec(fmt.Sprintf(departmentSQLString[mysqlDepartmentCreateTable], DBName, TableName))
	return err
}

func InsertDepartment(db *sql.DB, DBName, TableName string, d *Department) (int64, error) {
	result, err := db.Exec(fmt.Sprintf(departmentSQLString[mysqlDepartmentInsert], DBName, TableName),
		d.Code, d.OrganizationCode, d.Name, d.Sort,
		d.Pcode, d.Icon, d.CreateTime, d.Path)
	if err != nil {
		return 0, err
	}

	if rows, err := result.RowsAffected(); rows == 0 {
		return 0, err
	}

	departmentId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return departmentId, nil

}

//CreateDepartmentMemberTable create an administrative userAuth
func CreateDepartmentMemberTable(db *sql.DB, DBName, TableName string) error {
	result, err := db.Exec(fmt.Sprintf(departmentMemberSQLString[mysqlDepartmentMemberCreateTable], DBName, TableName))
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

//InsertDepartmentMember create an
func InsertDepartmentMember(db *sql.DB, DBName, TableName string, dm *DepartmentMember) (int64, error) {
	result, err := db.Exec(fmt.Sprintf(departmentMemberSQLString[mysqlDepartmentMemberInsert], DBName, TableName),
		dm.Code, dm.DepartmentCode, dm.OrganizationCode,
		dm.AccountCode, dm.JoinTime, dm.IsPrincipal, dm.IsOwner,
		dm.Authorize)
	if err != nil {
		return 0, err
	}

	if rows, err := result.RowsAffected(); rows == 0 {
		return 0, err
	}

	departmentMemberId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return departmentMemberId, nil
}

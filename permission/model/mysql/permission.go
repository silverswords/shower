/*
 * Revision History:
 *     Initial: 2020/10/05        Abserari
 */

package mysql

import (
	"database/sql"
	"errors"
	"time"
)

type (
	//Role -
	Role struct {
		RoleID   uint32
		Name     string
		Intro    string
		Active   bool
		CreateAt string
	}
	//Permission -
	Permission struct {
		URL       string
		RoleID    uint32
		CreatedAt string
	}
	//RelationData -
	RelationData struct {
		AdminID uint32
		RoleID  uint32
	}
)

const (
	mysqlRoleCreateTable = iota
	mysqlRoleInsert
	mysqlRoleModify
	mysqlRoleModifyActive
	mysqlRoleGetList
	mysqlRoleGetByID
	mysqlRoleGetIsActive
)

const (
	mysqlPermissionCreateTable = iota
	mysqlPermissionInstert
	mysqlPermissionDelete
	mysqlPermissonGetRole
	mysqlPermissonGetAll
)

const (
	mysqlRelationCreateTable = iota
	mysqlRelationInsert
	mysqlRelationDelete
	mysqlRelationRoleMap
	mysqlRelationSelectAdmin
	mysqlRelationSelectRole
)

var (
	errInvalidMysql  = errors.New("affected 0 rows")
	errRoleInactive  = errors.New("the role is not activated")

	roleSQLString = []string{
		`CREATE TABLE IF NOT EXISTS role (
			role_id 	INT UNSIGNED NOT NULL AUTO_INCREMENT,
			name		VARCHAR(512) UNIQUE NOT NULL DEFAULT ' ',
			intro		VARCHAR(512) NOT NULL DEFAULT ' ',
			active		BOOLEAN DEFAULT TRUE,
			created_at 	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (role_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO role(name,intro,active) VALUES (?,?,?)`,
		`UPDATE role SET name = ?,intro = ? WHERE role_id = ? LIMIT 1`,
		`UPDATE role SET active = ? WHERE role_id = ? LIMIT 1`,
		`SELECT * FROM role LOCK IN SHARE MODE`,
		`SELECT * FROM role WHERE role_id = ? AND active = true LOCK IN SHARE MODE`,
		`SELECT active FROM role WHERE role_id = ? LOCK IN SHARE MODE`,
	}

	permissionSQLString = []string{
		`CREATE TABLE IF NOT EXISTS permission (
			url			VARCHAR(512) NOT NULL DEFAULT ' ',
			role_id		MEDIUMINT UNSIGNED NOT NULL,
			created_at 	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (url,role_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO permission(url,role_id) VALUES (?,?)`,
		`DELETE FROM permission WHERE role_id = ? AND url = ? LIMIT 1`,
		`SELECT permission.role_id FROM permission, role WHERE url = ? AND role.active = true AND permission.role_id = role.role_id LOCK IN SHARE MODE`,
		`SELECT * FROM permission LOCK IN SHARE MODE`,
	}

	relationSQLString = []string{
		`CREATE TABLE IF NOT EXISTS relation (
			admin_id 	BIGINT UNSIGNED NOT NULL,
			role_id		INT UNSIGNED NOT NULL,
			created_at 	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (admin_id,role_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO relation(admin_id,role_id,created_at) VALUES (?,?,?)`,
		`DELETE FROM relation WHERE admin_id = ? AND role_id = ? LIMIT 1`,
		`SELECT relation.role_id FROM relation, role WHERE relation.admin_id = ? AND role.active = true AND relation.role_id = role.role_id LOCK IN SHARE MODE`,
		`SELECT relation.admin_id FROM admin, relation,role WHERE relation.role_id = ? AND role.active = true AND admin.active AND relation.admin_id = admin.admin_id LOCK IN SHARE MODE`,
		`SELECT relation.role_id FROM relation, role WHERE  role.active = true AND relation.role_id = role.role_id LOCK IN SHARE MODE`,
	}
)

// CreateTable create role table.
func CreateTable(db *sql.DB) error {
	_, err := db.Exec(roleSQLString[mysqlRoleCreateTable])
	if err != nil {
		return err
	}

	_, err = db.Exec(permissionSQLString[mysqlPermissionCreateTable])
	if err != nil {
		return err
	}

	_, err = db.Exec(relationSQLString[mysqlRelationCreateTable])
	if err != nil {
		return err
	}

	return nil
}

// CreateRole create a new role information.
func CreateRole(db *sql.DB, name, intro *string) error {
	result, err := db.Exec(roleSQLString[mysqlRoleInsert], name, intro, true)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

// ModifyRole modify role information.
func ModifyRole(db *sql.DB, id uint32, name, intro *string) error {
	_, err := db.Exec(roleSQLString[mysqlRoleModify], name, intro, id)

	return err
}

// ModifyRoleActive modify role active.
func ModifyRoleActive(db *sql.DB, id uint32, active bool) error {
	_, err := db.Exec(roleSQLString[mysqlRoleModifyActive], active, id)

	return err
}

// RoleList get all role information.
func RoleList(db *sql.DB) ([]*Role, error) {
	var (
		roleid   uint32
		name     string
		intro    string
		active   bool
		createAt string
		roles    []*Role
	)

	rows, err := db.Query(roleSQLString[mysqlRoleGetList])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&roleid, &name, &intro, &active, &createAt); err != nil {
			return nil, err
		}

		r := &Role{
			RoleID:   roleid,
			Name:     name,
			Intro:    intro,
			Active:   active,
			CreateAt: createAt,
		}

		roles = append(roles, r)
	}

	return roles, nil
}

// GetRoleByID get role by id.
func GetRoleByID(db *sql.DB, id uint32) (*Role, error) {
	var (
		r Role
	)

	err := db.QueryRow(roleSQLString[mysqlRoleGetByID], id).Scan(&r.RoleID, &r.Name, &r.Intro, &r.Active, &r.CreateAt)
	return &r, err
}

// AddURLPermission -
func AddURLPermission(db *sql.DB, rid uint32, url string) error {
	roleIsActive, err := IsActive(db, rid)
	if err != nil {
		return err
	}

	if !roleIsActive {
		return errRoleInactive
	}

	_, err = db.Exec(permissionSQLString[mysqlPermissionInstert], url, rid)
	return err
}

// RemoveURLPermission -
func RemoveURLPermission(db *sql.DB, rid uint32, url string) error {
	roleIsActive, err := IsActive(db, rid)
	if err != nil {
		return err
	}

	if !roleIsActive {
		return errRoleInactive
	}

	_, err = db.Exec(permissionSQLString[mysqlPermissionDelete], rid, url)
	return err
}

// URLPermissions lists all the roles of the specified URL.
func URLPermissions(db *sql.DB, url *string) (map[uint32]bool, error) {
	var (
		roleID uint32
		result = make(map[uint32]bool)
	)

	rows, err := db.Query(permissionSQLString[mysqlPermissonGetRole], url)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&roleID); err != nil {
			return nil, err
		}
		result[roleID] = true
	}

	return result, nil
}

// Permissions lists all the roles.
func Permissions(db *sql.DB) (*[]*Permission, error) {
	var (
		roleID    uint32
		url       string
		createdAt string

		result []*Permission
	)

	rows, err := db.Query(permissionSQLString[mysqlPermissonGetAll])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&url, &roleID, &createdAt); err != nil {
			return nil, err
		}
		data := &Permission{
			URL:       url,
			RoleID:    roleID,
			CreatedAt: createdAt,
		}
		result = append(result, data)
	}

	return &result, nil
}

// AddRelation add an relation
func AddRelation(db *sql.DB, aid, rid uint32) error {
	roleIsActive, err := IsActive(db, rid)
	if err != nil {
		return err
	}

	if !roleIsActive {
		return errRoleInactive
	}

	result, err := db.Exec(relationSQLString[mysqlRelationInsert], aid, rid, time.Now())
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

// RemoveRelation delate an relation
func RemoveRelation(db *sql.DB, aid, rid uint32) error {
	_, err := db.Exec(relationSQLString[mysqlRelationDelete], aid, rid)
	return err
}

// AdminGetRoleMap list all the roles of the specified admin and the return form is map.
func AdminGetRoleMap(db *sql.DB, aid uint32) (map[uint32]bool, error) {
	var (
		roleID uint32
		result = make(map[uint32]bool)
	)

	rows, err := db.Query(relationSQLString[mysqlRelationRoleMap], aid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&roleID); err != nil {
			return nil, err
		}
		result[roleID] = true
	}

	return result, nil
}

// AssociatedRoleList list all the roles of the specified admin and the return form is slice.
func AssociatedRoleList(db *sql.DB, aid uint32) ([]*RelationData, error) {
	var (
		roleID uint32
		r      *RelationData
		result []*RelationData
	)

	rows, err := db.Query(relationSQLString[mysqlRelationRoleMap], aid)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&roleID); err != nil {
			return nil, err
		}
		r = &RelationData{
			AdminID: aid,
			RoleID:  roleID,
		}
		result = append(result, r)
	}

	return result, nil
}

//IsActive return Active and nil if query success.
func IsActive(db *sql.DB, id uint32) (bool, error) {
	var (
		isActive bool
	)

	db.QueryRow(roleSQLString[mysqlRoleGetIsActive], id).Scan(&isActive)
	return isActive, nil
}

// GetAdminIDMap list all the roles of the specified admin and the return form is map.
func GetAdminIDMap(db *sql.DB) (map[uint32]bool, error) {
	var (
		AdminID uint32
		result  = make(map[uint32]bool)
	)

	slice1, err := RoleList(db)
	if err != nil {
		return nil, err
	}

	for _, aid := range slice1 {
		rows, err := db.Query(relationSQLString[mysqlRelationSelectAdmin], aid.RoleID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&AdminID); err != nil {
				return nil, err
			}
			result[AdminID] = true
		}
	}

	return result, nil
}

// GetRoleIDMap list all the roles of the specified admin and the return form is map.
func GetRoleIDMap(db *sql.DB) (map[uint32]bool, error) {
	var (
		AdminID uint32
		result  = make(map[uint32]bool)
	)

	rows, err := db.Query(relationSQLString[mysqlRelationSelectRole])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&AdminID); err != nil {
			return nil, err
		}
		result[AdminID] = true
	}

	return result, nil
}

package Department

import (
	"database/sql"
	"log"
	"net/http"

	mysql "github.com/abserari/shower/pkgs/department/model"

	"github.com/gin-gonic/gin"
)

// Controller -
type Controller struct {
	db        *sql.DB
	tableName string
	dbName    string
}

// Register -
func Register(db *sql.DB, tableName string, dbName string, r gin.IRouter) error {
	c := New(db, tableName, dbName)

	if err := c.CreateTable(); err != nil {
		log.Fatal(err)
		return err
	}

	r.POST("/department/create", c.InsertDepartment)
	r.POST("/department-member/create", c.InsertDepartmentMember)
	return nil
}

// New -
func New(db *sql.DB, tableName string, dbName string) *Controller {
	return &Controller{
		db:        db,
		tableName: tableName,
		dbName:    dbName,
	}
}

// CreateTable -
func (con *Controller) CreateTable() error {
	err := mysql.CreateDepartmentTable(con.db, con.dbName, con.tableName)
	if err != nil {
		log.Fatal(err)
	}
	err = mysql.CreateDepartmentMemberTable(con.db, con.dbName, con.tableName)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Insert -
func (con *Controller) InsertDepartment(c *gin.Context) {
	var (
		d *mysql.Department
	)

	if err := c.ShouldBindJSON(d); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}

	_, err := mysql.InsertDepartment(con.db, con.dbName, con.tableName, d)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	return
}

// Insert -
func (con *Controller) InsertDepartmentMember(c *gin.Context) {
	var (
		dm *mysql.DepartmentMember
	)

	if err := c.ShouldBindJSON(dm); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}

	_, err := mysql.InsertDepartmentMember(con.db, con.dbName, con.tableName, dm)

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	return
}

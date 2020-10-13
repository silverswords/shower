package category

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/abserari/shower/category/model/mysql"
	"github.com/gin-gonic/gin"
)

// Controller -
type Controller struct {
	db        *sql.DB
	tableName string
	dbName    string
}

// Config -
type Config struct {
	CategoryDB    string
	CategoryTable string
}

// Register -
func Register(db *sql.DB, tableName string, dbName string, r gin.IRouter) error {
	c := New(db, tableName, dbName)

	if err := c.CreateTable(); err != nil {
		log.Fatal(err)
		return err
	}

	r.POST("/api/v1/category/create", c.Insert)
	r.POST("/api/v1/category/modify/status", c.ChangeCategoryStatus)
	r.POST("/api/v1/category/modify/name", c.ChangeCategoryName)
	r.POST("/api/v1/category/children", c.LisitChirldrenByParentID)
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
	err := mysql.CreateTable(con.db, con.tableName)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Insert -
func (con *Controller) Insert(c *gin.Context) {
	var (
		req struct {
			ParentID uint   `json:"parentId"`
			Name     string `json:"name"`
		}
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}

	_, err := mysql.InsertCategory(con.db, con.tableName, req.ParentID, req.Name)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	return
}

// ChangeCategoryStatus -
func (con *Controller) ChangeCategoryStatus(c *gin.Context) {
	var (
		req struct {
			CategoryID uint `json:"categoryId"`
			Status     int8 `json:"status"`
		}
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}

	err := mysql.ChangeCategoryStatus(con.db, con.tableName, req.CategoryID, req.Status)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	return
}

// ChangeCategoryName -
func (con *Controller) ChangeCategoryName(c *gin.Context) {
	var (
		req struct {
			CategoryID uint   `json:"categoryId"`
			Name       string `json:"name"`
		}
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}

	err := mysql.ChangeCategoryName(con.db, con.tableName, req.CategoryID, req.Name)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	return
}

// LisitChirldrenByParentID -
func (con *Controller) LisitChirldrenByParentID(c *gin.Context) {
	var (
		req struct {
			ParentID uint `json:"parentId"`
		}
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}

	list, err := mysql.LisitChirldrenByParentID(con.db, con.tableName, req.ParentID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "list": list})
	return
}

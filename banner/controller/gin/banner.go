/*
 * Revision History:
 *     Initial: 2019/03/18        Yang ChengKai
 */

package controller

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	mysql "github.com/abserari/shower/banner/model/mysql"
	"github.com/gin-gonic/gin"
)

// BannerController -
type BannerController struct {
	db        *sql.DB
	tableName string
}

// New -
func New(db *sql.DB, tableName string) *BannerController {
	return &BannerController{
		db:        db,
		tableName: tableName,
	}
}

// RegisterRouter -
func (b *BannerController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := mysql.CreateTable(b.db, b.tableName)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/create", b.create)
	r.POST("/delete", b.deleteByID)
	r.POST("/info/id", b.infoByID)
	r.POST("/list/date", b.lisitValidBannerByUnixDate)
}

func (b *BannerController) create(c *gin.Context) {
	var (
		req struct {
			Name      string    `json:"name"      binding:"required"`
			ImagePath string    `json:"imageurl"  binding:"required"`
			EventPath string    `json:"eventurl"  binding:"required"`
			StartDate time.Time `json:"start_date"`
			EndDate   time.Time `json:"end_date"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	id, err := mysql.InsertBanner(b.db, b.tableName, req.Name, req.ImagePath, req.EventPath, req.StartDate, req.EndDate)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ID": id})
}

func (b *BannerController) lisitValidBannerByUnixDate(c *gin.Context) {
	var (
		req struct {
			Unixtime int64 `json:"unixtime"    binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	banners, err := mysql.LisitValidBannerByUnixDate(b.db, b.tableName, req.Unixtime)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "banners": banners})
}

func (b *BannerController) infoByID(c *gin.Context) {
	var (
		req struct {
			ID int `json:"id"     binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	ban, err := mysql.InfoByID(b.db, b.tableName, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ban": ban})
}

func (b *BannerController) deleteByID(c *gin.Context) {
	var (
		req struct {
			ID int `json:"id"    binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.DeleteByID(b.db, b.tableName, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

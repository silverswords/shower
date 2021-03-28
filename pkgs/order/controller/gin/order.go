package order

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/abserari/shower/order/model/mysql"
	"github.com/gin-gonic/gin"
)

// type Stocker interface {
// 	ModifyProductStock(tx *sql.Tx, targetID uint32, num int) error
// }

// type UserChecker interface {
// 	UserCheck(tx *sql.Tx, userid uint64, productID uint32) error
// }

type Config struct {
	OrderDB        string
	OrderTable     string
	ItemTable      string
	ClosedInterval int
	// Stock          Stocker
	// User           UserChecker
}

type Controller struct {
	db  *sql.DB
	Cnf Config
}

func Register(r gin.IRouter, db *sql.DB) error {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}
	cnf := Config{
		OrderDB:        "test",
		OrderTable:     "orderTable",
		ItemTable:      "Items",
		ClosedInterval: 5,
	}
	c := New(db, cnf)

	if err := c.CreateDB(); err != nil {
		log.Fatal(err)
		return err
	}

	if err := c.CreateOrderTable(); err != nil {
		log.Fatal(err)
		return err
	}

	if err := c.CreateItemTable(); err != nil {
		log.Fatal(err)
		return err
	}

	r.POST("/api/v1/order/create", c.Insert)
	r.POST("/api/v1/order/info", c.OrderInfoByOrderID)
	r.POST("/api/v1/order/user", c.LisitOrderByUserIDAndStatus)
	r.POST("/api/v1/order/id", c.OrderIDByOrderCode)

	return nil
}

// New -
func New(db *sql.DB, cnf Config) *Controller {
	return &Controller{
		db:  db,
		Cnf: cnf,
	}
}

// CreateDB -
func (ctl *Controller) CreateDB() error {
	return mysql.CreateDB(ctl.db, ctl.Cnf.OrderDB)
}

// CreateOrderTable -
func (ctl *Controller) CreateOrderTable() error {
	ostore := ctl.Cnf.OrderDB + "." + ctl.Cnf.OrderTable
	return mysql.CreateTable(ctl.db, ostore)
}

// CreateItemTable -
func (ctl *Controller) CreateItemTable() error {
	istore := ctl.Cnf.OrderDB + "." + ctl.Cnf.ItemTable
	return mysql.CreateTable(ctl.db, istore)
}

// Insert -
func (ctl *Controller) Insert(c *gin.Context) {
	var (
		req struct {
			UserID     uint64 `json:"userid"`
			AddressID  string `json:"addressid"`
			TotalPrice uint32 `json:"totalprice"`
			Promotion  string `json:"promotion"`
			Freight    uint32 `json:"freight"`

			Items []mysql.Item `json:"items"`
		}
		rep struct {
			ordercode string
			orderid   uint32
		}
		err error
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}
	promotion, err := strconv.ParseBool(req.Promotion)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}

	times := time.Now()
	rep.ordercode = strconv.Itoa(times.Year()) + strconv.Itoa(int(times.Month())) + strconv.Itoa(times.Day()) + strconv.Itoa(times.Hour()) + strconv.Itoa(times.Minute()) + strconv.Itoa(times.Second()) + strconv.Itoa(int(req.UserID))
	order := mysql.Order{
		OrderCode:  rep.ordercode,
		UserID:     req.UserID,
		AddressID:  req.AddressID,
		TotalPrice: req.TotalPrice,
		Promotion:  promotion,
		Freight:    req.Freight,
		Created:    times,
	}

	rep.orderid, err = mysql.Insert(order, req.Items, ctl.db, ctl.Cnf.ClosedInterval, ctl.Cnf.OrderDB, ctl.Cnf.OrderTable, ctl.Cnf.ItemTable)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, gin.H{"statue": http.StatusOK, "orderid": rep.orderid, "ordercode": rep.ordercode})
	return
}

//optional
func (ctl *Controller) OrderIDByOrderCode(c *gin.Context) {
	var req struct {
		Ordercode string `json:"ordercode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}
	ostore := ctl.Cnf.OrderDB + "." + ctl.Cnf.OrderTable
	id, err := mysql.OrderIDByOrderCode(ctl.db, ostore, req.Ordercode)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, gin.H{"statue": http.StatusOK, "id": id})
	return
}

//full info for One Order
func (ctl *Controller) OrderInfoByOrderID(c *gin.Context) {
	var req struct {
		OrderID uint32 `json:"orderid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}
	ostore := ctl.Cnf.OrderDB + "." + ctl.Cnf.OrderTable
	istore := ctl.Cnf.OrderDB + "." + ctl.Cnf.ItemTable
	rep, err := mysql.SelectByOrderKey(ctl.db, ostore, istore, req.OrderID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, gin.H{"statue": http.StatusOK, "Order": rep.Order, "Orm": rep.Orm})
	return
}

/*
mode:
  Unfinished = 0
  Finished   = 1
  Paid       = 2
  Consigned  = 3
  Canceled   = 4
*/
func (ctl *Controller) LisitOrderByUserIDAndStatus(c *gin.Context) {
	var req struct {
		Userid uint64 `json:"userid"`
		Status uint8  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}
	ostore := ctl.Cnf.OrderDB + "." + ctl.Cnf.OrderTable
	istore := ctl.Cnf.OrderDB + "." + ctl.Cnf.ItemTable
	orders, err := mysql.LisitOrderByUserID(ctl.db, ostore, istore, req.Userid, req.Status)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, gin.H{"statue": http.StatusOK, "Orders": orders})
	return
}

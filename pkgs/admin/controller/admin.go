/*
 * Revision History:
 *     Initial: 2020/10/05        Abserari
 */

package controller

import (
	"database/sql"
	"errors"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"log"
	"net/http"

	"github.com/abserari/shower/admin/model/mysql"
	"github.com/gin-gonic/gin"
)

var (
	// default user
	name     = "Admin"
	password = "111111"

	errActive          = errors.New("the admin is not activated")
	errUserIDNotExists = errors.New("Get Admin ID is not exists")
	errUserIDNotValid  = func(value interface{}) error {
		return errors.New(fmt.Sprintf("Get Admin ID is not valid. Is %s", value))
	}
)

// Controller external service interface
type Controller struct {
	db  *sql.DB
	JWT *jwt.GinJWTMiddleware
}

// New create an external service interface
func New(db *sql.DB) *Controller {
	c := &Controller{
		db: db,
	}
	var err error
	c.JWT, err = c.newJWTMiddleware()
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// RegisterRouter register router. It fatal because there is no service if register failed.
func (con *Controller) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}
	err := mysql.CreateDatabase(con.db)
	if err != nil {
		log.Fatal(err)
	}

	err = mysql.CreateTable(con.db, &name, &password)
	if err != nil {
		log.Fatal(err)
	}
	
	// admin crud API
	r.POST("/create", con.create)
	r.POST("/modify/email", con.modifyEmail)
	r.POST("/modify/mobile", con.modifyMobile)
	r.POST("/modify/password", con.modifyPassword)
	r.POST("/modify/active", con.modifyAdminActive)
}

func (con *Controller) create(ctx *gin.Context) {
	var (
		admin struct {
			Name     string `json:"name"      binding:"required,alphanum,min=5,max=30"`
			Password string `json:"password"  binding:"omitempty,min=5,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	//Default password
	if admin.Password == "" {
		admin.Password = "111111"
	}

	err = mysql.CreateAdmin(con.db, &admin.Name, &admin.Password)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (con *Controller) modifyEmail(ctx *gin.Context) {
	var (
		admin struct {
			AdminID uint32 `json:"admin_id"    binding:"required"`
			Email   string `json:"email"       binding:"required,email"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyEmail(con.db, admin.AdminID, &admin.Email)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (con *Controller) modifyMobile(ctx *gin.Context) {
	var (
		admin struct {
			AdminID uint32 `json:"admin_id"     binding:"required"`
			Mobile  string `json:"mobile"       binding:"required,numeric,len=11"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyMobile(con.db, admin.AdminID, &admin.Mobile)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (con *Controller) modifyPassword(ctx *gin.Context) {
	var (
		admin struct {
			AdminID     uint32 `json:"admin_id"       binding:"required"`
			Password    string `json:"password"       binding:"printascii,min=6,max=30"`
			NewPassword string `json:"new_password"   binding:"printascii,min=6,max=30"`
			Confirm     string `json:"confirm"        binding:"printascii,min=6,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if admin.NewPassword == admin.Password {
		ctx.Error(err)
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	if admin.NewPassword != admin.Confirm {
		ctx.Error(err)
		ctx.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict})
		return
	}

	err = mysql.ModifyPassword(con.db, admin.AdminID, &admin.Password, &admin.NewPassword)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (con *Controller) modifyAdminActive(ctx *gin.Context) {
	var (
		admin struct {
			CheckID     uint32 `json:"check_id"    binding:"required"`
			CheckActive bool   `json:"check_active"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyAdminActive(con.db, admin.CheckID, admin.CheckActive)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

//Login JWT validation
func (con *Controller) Login(ctx *gin.Context) (uint32, error) {
	var (
		admin struct {
			Name     string `json:"name"      binding:"required,alphanum,min=5,max=30"`
			Password string `json:"password"  binding:"printascii,min=6,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		return 0, err
	}

	ID, err := mysql.Login(con.db, &admin.Name, &admin.Password)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

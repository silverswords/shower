/*
 * Revision History:
 *     Initial: 2020/1018       Abserari
 */

package controller

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	mysql "github.com/abserari/shower/pet/model/mysql"
	"github.com/gin-gonic/gin"
)

// PetController -
type PetController struct {
	db        *sql.DB
	tableName string
}

// New -
func New(db *sql.DB, tableName string) *PetController {
	return &PetController{
		db:        db,
		tableName: tableName,
	}
}

// RegisterRouter -
func (b *PetController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := mysql.CreateTable(b.db, b.tableName)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/create", b.create)

	r.POST("/update/name", b.modifyName)
	r.POST("/update/category", b.modifyCategory)
	r.POST("/update/avatar", b.modifyAvatar)
	r.POST("/update/birthday", b.modifyBirthday)
	r.POST("/update/medicalcurrent", b.modifyMedicalCurrent)
	r.POST("/update/hobbies", b.modifyHobbies)
	r.POST("/update/gender", b.modifyGender)

	r.POST("/delete", b.deleteByID)

	r.POST("/info/id", b.infoByID)
	r.POST("/list/adminid", b.listPetByAdminID)
}

func (b *PetController) create(c *gin.Context) {
	var (
		req struct {
			AdminID        uint64    `json:"adminID"    binding:"required"`
			Name           string    `json:"name"      binding:"required"`
			Category       string    `json:"category" `
			Avatar         string    `json:"avatar" `
			Birthday       time.Time `json:"birthday" `
			MedicalCurrent string    `json:"medicalCurrent" `
			Hobbies        string    `json:"hobbies" `
			Gender         string    `json:"gender" `
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	id, err := mysql.InsertPet(b.db, b.tableName, req.AdminID, req.Name, req.Category, req.Avatar, req.Birthday, req.MedicalCurrent, req.Hobbies, req.Gender)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ID": id})
}

func (b *PetController) listPetByAdminID(c *gin.Context) {
	var (
		req struct {
			AdminID uint64 `json:"adminID"    binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	pets, err := mysql.ListPetByAdminID(b.db, b.tableName, req.AdminID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "pets": pets})
}

func (b *PetController) infoByID(c *gin.Context) {
	var (
		req struct {
			ID uint64 `json:"id"     binding:"required"`
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
func (con *PetController) modifyName(ctx *gin.Context) {
	var (
		admin struct {
			PetID uint64 `json:"petID"    binding:"required"`
			Name  string `json:"name"       binding:"required,email"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyName(con.db, admin.PetID, admin.Name)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (con *PetController) modifyCategory(ctx *gin.Context) {
	var (
		admin struct {
			PetID    uint64 `json:"petID"    binding:"required"`
			Category string `json:"category"  binding:"required,email"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyCategory(con.db, admin.PetID, admin.Category)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
func (con *PetController) modifyAvatar(ctx *gin.Context) {
	var (
		admin struct {
			PetID  uint64 `json:"petID"     binding:"required"`
			Avatar string `json:"avatar"       binding:"required,email"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyAvatar(con.db, admin.PetID, admin.Avatar)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
func (con *PetController) modifyBirthday(ctx *gin.Context) {
	var (
		admin struct {
			PetID    uint64    `json:"petID"    binding:"required"`
			Birthday time.Time `json:"birthday"     binding:"required,email"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyBirthday(con.db, admin.PetID, admin.Birthday)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
func (con *PetController) modifyMedicalCurrent(ctx *gin.Context) {
	var (
		admin struct {
			PetID          uint64 `json:"petID"    binding:"required"`
			MedicalCurrent string `json:"medicalCurrent"     binding:"required,email"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyMedicalCurrent(con.db, admin.PetID, admin.MedicalCurrent)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
func (con *PetController) modifyHobbies(ctx *gin.Context) {
	var (
		admin struct {
			PetID   uint64 `json:"petID"    binding:"required"`
			Hobbies string `json:"hobbies"    binding:"required,email"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyHobbies(con.db, admin.PetID, admin.Hobbies)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
func (con *PetController) modifyGender(ctx *gin.Context) {
	var (
		admin struct {
			PetID  uint64 `json:"petID"    binding:"required"`
			Gender string `json:"gender"     binding:"required,email"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyGender(con.db, admin.PetID, admin.Gender)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
func (b *PetController) deleteByID(c *gin.Context) {
	var (
		req struct {
			ID uint64 `json:"id"    binding:"required"`
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

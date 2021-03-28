/*
 * Revision History:
 *     Initial: 2019/03/20        Yang ChengKai
 */

package controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/abserari/shower/smservice/model/mysql"
	service "github.com/abserari/shower/smservice/service"
	"github.com/gin-gonic/gin"
)

// SMController -
type SMController struct {
	ser *service.Controller
}

// New -
func New(db *sql.DB, conf *service.Config) *SMController {
	return &SMController{
		ser: service.NewController(db, conf),
	}
}

// RegisterRouter -
func (s *SMController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := mysql.CreateTable(s.ser.DB)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/send", s.Send)
	r.POST("/check", s.Check)
}

// Send 调度分配出发送短信
func (s *SMController) Send(c *gin.Context) {
	var (
		req struct {
			Mobile string `json:"mobile"`
			Sign   string `json:"sign"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	if err = service.Send(req.Mobile, req.Sign, &s.ser.Conf, s.ser.DB); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

// Check 调度分配检查验证码
func (s *SMController) Check(c *gin.Context) {
	var (
		req struct {
			Code string `json:"code"`
			Sign string `json:"sign"`
		}

		resp struct {
			sign   string
			mobile string
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	resp.sign = req.Sign
	resp.mobile, _ = mysql.GetMobile(s.ser.DB, resp.sign)

	if err = service.Check(req.Code, req.Sign, &s.ser.Conf, s.ser.DB); err != nil {
		s.ser.Conf.OnCheck.OnVerifyFailed(resp.sign, resp.mobile)

		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	s.ser.Conf.OnCheck.OnVerifySucceed(resp.sign, resp.mobile)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

// main use to config and run server.
package main

import (
	"database/sql"
	"log"

	admin "github.com/abserari/shower/pkgs/userAuth/controller"
	permission "github.com/abserari/shower/pkgs/permission/controller/gin"
	pet "github.com/abserari/shower/pkgs/pet/controller/gin"
	smservice "github.com/abserari/shower/pkgs/smservice/controller/gin"
	service "github.com/abserari/shower/pkgs/smservice/service"
	upload "github.com/abserari/shower/pkgs/upload/controller/gin"
	"github.com/abserari/shower/utils/fileserver"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type funcv struct{}

func (v funcv) OnVerifySucceed(targetID, mobile string) {}
func (v funcv) OnVerifyFailed(targetID, mobile string)  {}

func main() {
	var v funcv

	router := gin.Default()

	dbConn, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/project?parseTime=true")
	if err != nil {
		panic(err)
	}

	con := &service.Config{
		Host:           "https://fesms.market.alicloudapi.com/sms/",
		Appcode:        "6f37345cad574f408bff3ede627f7014",
		Digits:         6,
		ResendInterval: 60,
		OnCheck:        v,
	}
	smserviceCon := smservice.New(dbConn, con)
	smserviceCon.RegisterRouter(router.Group("/api/v1/message"))

	adminCon := admin.New(dbConn)
	// login and refresh token.
	router.POST("/api/v1/userAuth/login", adminCon.JWT.LoginHandler)
	router.GET("/api/v1/userAuth/refresh_token", adminCon.JWT.RefreshHandler)
	// start to add token on every API after userAuth.RegisterRouter
	router.Use(adminCon.JWT.MiddlewareFunc())
	// start to check the userAuth active every time.
	router.Use(adminCon.CheckActive())
	adminCon.RegisterRouter(router.Group("/api/v1/userAuth"))

	permissionCon := permission.New(dbConn, adminCon.GetID)
	router.Use(permissionCon.CheckPermission())
	permissionCon.RegisterRouter(router.Group("/api/v1/permission"))

	petCon := pet.New(dbConn, "pet")
	petCon.RegisterRouter(router.Group("/api/v1/pet"))

	uploadCon := upload.New(dbConn, "0.0.0.0:9573", adminCon.GetID)
	uploadCon.RegisterRouter(router.Group("/api/v1/userAuth"))

	go fileserver.StartFileServer("0.0.0.0:9573", "")
	log.Fatal(router.Run(":8000"))
}

// main use to config and run server.
package main

import (
	"database/sql"
	"log"

	permission "github.com/abserari/shower/pkgs/permission/controller/gin"
	pet "github.com/abserari/shower/pkgs/pet/controller/gin"
	upload "github.com/abserari/shower/pkgs/upload/controller/gin"
	admin "github.com/abserari/shower/pkgs/userAuth/controller"
	"github.com/abserari/shower/utils/fileserver"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)


func main() {
	router := gin.Default()

	dbConn, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/project?parseTime=true")
	if err != nil {
		panic(err)
	}

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

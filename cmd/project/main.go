// main use to config and run server.
package main

import (
	"database/sql"
	"log"

	permission "github.com/abserari/shower/pkgs/permission/controller/gin"
	upload "github.com/abserari/shower/pkgs/upload/controller/gin"
	admin "github.com/abserari/shower/pkgs/userAuth/controller"
	"github.com/abserari/shower/utils/fileserver"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// config
const (
	uploadAddressBase        = "0.0.0.0:9573"
	serverAddressBase        = ":8000"
	mysqlConnAddress         = "root:123456@tcp(localhost:3306)/project?parseTime=true"

 	userAuthRouterGroup = "/api/v1/userAuth"
	userAuthRouterGroupLogin = userAuthRouterGroup + "/login"
	userAuthRouterRefreshToken = userAuthRouterGroup +"/refresh_token"
 	permissionRouterGroup = "/api/v1/permission"
 	uploadRouterGroup = "/api/v1/upload"
)

func main() {
	router := gin.Default()

	dbConn, err := sql.Open("mysql", mysqlConnAddress)
	if err != nil {
		panic(err)
	}

	// init controller with db conn
	adminCon := admin.New(dbConn)
	permissionCon := permission.New(dbConn, adminCon.GetID)
	uploadCon := upload.New(dbConn, uploadAddressBase, adminCon.GetID)

	// register router and MiddlewareFunc

	// login and refresh token.
	router.POST(userAuthRouterGroupLogin, adminCon.JWT.LoginHandler)
	router.GET(userAuthRouterRefreshToken, adminCon.JWT.RefreshHandler)

	// start to add token on every API after userAuth.RegisterRouter
	router.Use(adminCon.JWT.MiddlewareFunc())
	// start to check the userAuth active every time.
	router.Use(adminCon.CheckActive())
	// start to check the userAuth permission every time.
	router.Use(permissionCon.CheckPermission())

	adminCon.RegisterRouter(router.Group(userAuthRouterGroup))
	permissionCon.RegisterRouter(router.Group(permissionRouterGroup))
	uploadCon.RegisterRouter(router.Group(uploadRouterGroup))

	// start the fileServer services
	go fileserver.StartFileServer(uploadAddressBase, "")
	log.Fatal(router.Run(serverAddressBase))
}

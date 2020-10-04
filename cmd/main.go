package main

import (
	"database/sql"

	banner "github.com/fengyfei/comet/banner/controller/gin"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()

	dbConn, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}
	bannerCon := banner.New(dbConn, "banner")

	bannerCon.RegisterRouter(router.Group("/api/v1/banner"))

	router.Run(":8000")
}

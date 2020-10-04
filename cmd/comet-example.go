package main

import (
	"database/sql"
	"os"

	comet "github.com/abserari/shower/comet"
	banner "github.com/fengyfei/comet/banner/controller/gin"

	"github.com/gin-gonic/gin"
)

const (
	// export COMET_DATABASE_URL=root:123456@tcp(127.0.0.1:3306)/test
	databaseurl = "COMET_DATABASE_URL"
)

func main() {
	router := gin.Default()

	url := os.Getenv(databaseurl)
	dbConn, err := sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	var banner comet.Comet = banner.New(dbConn, "banner")
	banner.RegisterRouter(router.Group("/api/v1/banner"))

	router.Run(":8000")
}

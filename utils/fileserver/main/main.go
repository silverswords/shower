package main

import (
	"os"

	"github.com/abserari/shower/utils/fileserver"
)

func main() {
	wdir, _ := os.Getwd()
	fileserver.StartFileServer(":9573", wdir)
}

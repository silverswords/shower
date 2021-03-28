package fileserver

import (
	"log"
	"net/http"
	"os"
)

// StartFileServer start file server with specified path to file.
func StartFileServer(host, path string) {
	// http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(path+"images/"))))
	var restart = 1
restart:
	if path == "" {
		path, _ = os.Getwd()
	}
	log.Println("Starting FileServer in ", host, "with", path)
	if err := http.ListenAndServe(host, http.FileServer(http.Dir(path))); err != nil {
		log.Println(err)
	}
	restart++
	if restart <= 10 {
		goto restart
	}
	log.Println("[File Server] : Starting 10 times... but fialed.")
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

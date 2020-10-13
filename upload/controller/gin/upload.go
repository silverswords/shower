/*
 * Revision History:
 *     Initial: 2020/1016       Abserari
 */

package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	mysql "github.com/abserari/shower/upload/model/mysql"
	md "github.com/abserari/shower/utils/file"
	"github.com/gin-gonic/gin"
)

var (
	errRequest = errors.New("Request is not post method")
	erruserID  = errors.New("userID invalid")
)

const (
	// InvalidUID - userID invalid
	InvalidUID = 0
	// FileKey - key of the file
	FileKey = "file"
	// FileUploadDir - the root directory of the upload files
	FileUploadDir = "files"
	// PictureDir - save pictures file
	PictureDir = "picture"
	// VideoDir - save videos file
	VideoDir = "video"
	// OtherDir - files other than video and picture
	OtherDir = "other"
)

// UploadController -
type UploadController struct {
	db      *sql.DB
	BaseURL string
	getUID  func(c *gin.Context) (uint32, error)
}

// New -
func New(db *sql.DB, baseURL string, getUID func(c *gin.Context) (uint32, error)) *UploadController {
	return &UploadController{
		db:      db,
		BaseURL: "http://" + baseURL + "/",
		getUID:  getUID,
	}
}

// RegisterRouter -
func (u *UploadController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := mysql.CreateTable(u.db)
	if err != nil {
		log.Fatal(err)
	}

	err = md.CheckDir(PictureDir, VideoDir, OtherDir)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/upload", u.upload)
	r.POST("/delete", u.deleteByID)
}

func (u *UploadController) upload(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Error(errRequest)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	userID, err := u.getUID(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized})
		return
	}

	if userID == InvalidUID {
		c.Error(erruserID)
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		return
	}

	file, header, err := c.Request.FormFile(FileKey)
	defer func() {
		file.Close()
		c.Request.MultipartForm.RemoveAll()
	}()

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound})
		return
	}

	newfile, _ := ioutil.ReadAll(file)

	MD5Str, err := md.MD5(newfile)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": http.StatusMethodNotAllowed})
		return
	}

	filePath, err := mysql.QueryByMD5(u.db, MD5Str)
	// if the file exists, return it now.
	if err == nil {
		fmt.Println("The file already exists:", filePath)
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable, "URL": u.BaseURL + filePath})
		return
	}

	// check the error if is our expected - NoRows.
	if err != mysql.ErrNoRows {
		c.Error(err)
		c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict})
		return
	}

	fileSuffix := path.Ext(header.Filename)
	filePath = FileUploadDir + "/" + md.ClassifyBySuffix(fileSuffix) + "/" + MD5Str + fileSuffix

	err = md.CopyFile(filePath, newfile)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	err = mysql.Insert(u.db, userID, filePath, MD5Str)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"status": http.StatusUnsupportedMediaType})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "URL": u.BaseURL + filePath})
}

func (con *UploadController) deleteByID(c *gin.Context) {
	var (
		req struct {
			Path string `json:"path"    binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	req.Path = req.Path[strings.Index(req.Path, "files/"):]
	log.Println(req.Path, con.BaseURL)
	err = mysql.DeleteByPath(con.db, req.Path)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

package main

import (
	"crypto/sha256"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"

	orderHistoryDb "github.com/haruki/OrderHistory/db"

	webui "github.com/haruki/OrderHistory/ui_api"

	reqHandler "github.com/haruki/OrderHistory/reqHandler"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	cors "github.com/rs/cors/wrapper/gin"
)

//go reference time
//Mon Jan 2 15:04:05 -0700 MST 2006

var (
	//go:embed reactFrontend/dist
	webuifs embed.FS
)

func main() {
	pathstring, _ := filepath.Abs("/mnt/d/20230101_orderHistory-sqlite.db")
	//pathstring, _ := filepath.Abs("d:/20230101_orderHistory-sqlite.db")
	log.Printf("DB Pfad: %v", pathstring)
	db, err := sql.Open("sqlite3", pathstring)
	if err != nil {
		log.Fatal(err)
	}

	handler := reqHandler.NewHandler(db)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New() //Default ersetzt durch New. Default hat einen debug logger, der nicht mehr benötigt wird.
	r.Use(gin.Recovery())
	r.Use(cors.AllowAll())
	//Images:
	r.Static("/img", "./img")
	//WebUI:
	dist, _ := fs.Sub(webuifs, "reactFrontend/dist")
	r.StaticFS("/webui", http.FS(dist)) //package.json -> "build": "vite build --base=/webui/"
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
	r.GET("/allitems", func(c *gin.Context) {
		err, items := webui.LoadAllItems(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, items)
		}
	})
	r.POST("/order/:platform", func(c *gin.Context) {
		platform := c.Param("platform")
		if platform == "aliexpress" {
			handler.Aliexpress(c)
		} else if platform == "ebay" {
			handler.Ebay(c)
		} else if platform == "alternate" {
			handler.Alternate(c)
		} else {
			c.JSON(404, gin.H{
				"message": "Fail (Platform not supported yet)",
			})
		}
	})
	r.POST("/imageUpload", func(c *gin.Context) {
		id := c.PostForm("itemId")
		intId, err := strconv.Atoi(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No valid id",
			})
			return
		}
		vendor := c.PostForm("vendor")
		file, err := c.FormFile("file")
		// The file cannot be received.
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No file is received",
			})
			return
		}
		sha256Hash, err := hashImage(file)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Error while hashing file",
			})
			return
		}
		//get filename extension (returns string including the dot)
		ext := filepath.Ext(file.Filename)
		newDbFileName := fmt.Sprintf("./img/%s_%s_%s%s", vendor, id, sha256Hash[:5], ext)
		newFilePathName := fmt.Sprintf("./img/backup/%s_%s_%s%s", vendor, id, sha256Hash[:5], ext)
		err = c.SaveUploadedFile(file, newFilePathName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Error while saving file",
			})
			return
		}

		err = orderHistoryDb.UpdateImage(db, newDbFileName, sha256Hash, intId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Error while updating DB",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"imgFile": newDbFileName})
	})
	log.Printf("Starting OrderHistory-Server at Port: %d", 8081)
	r.Run(":8081")
}

func hashImage(file *multipart.FileHeader) (string, error) {
	realfile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer realfile.Close()
	fileBytes, err := ioutil.ReadAll(realfile)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	hash.Write(fileBytes)
	sha256Hash := fmt.Sprintf("%x", hash.Sum(nil))
	return sha256Hash, nil
}

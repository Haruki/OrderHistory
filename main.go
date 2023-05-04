package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

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
		vendor := c.PostForm("vendor")
		file, err := c.FormFile("file")
		// The file cannot be received.
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No file is received",
			})
			return
		}
		//get filename extension
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("./img/backup/%s_%s.%s", vendor, id, ext)
		err = c.SaveUploadedFile(file, newFileName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Error while saving file",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Your file has been successfully uploaded."})
	})
	log.Printf("Starting OrderHistory-Server at Port: %d", 8081)
	r.Run(":8081")
}

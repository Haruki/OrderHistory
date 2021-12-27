package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	cors "github.com/rs/cors/wrapper/gin"
)

type TestType struct {
	variable1 string
	variable2 int
}

type ebay struct {
	Id             int64
	PurchaseDate   string `jdon:"purchaseDate"`
	ItemName       string `json:"itemName"`
	VendorPlatform string
	Price          float32 `json:"price"`
	ImgUrl         string  `json:"imgUrl"`
	//special ebay variables:
	Vendor        string `json:"vendor"`
	Artikelnummer int    `json:"artikelnummer"`
}

func (t TestType) test(kalr string) {
	log.Printf("variable1: %v %v", t.variable1, kalr)
	log.Printf("variable2: %v %v", t.variable2, kalr)
}

func main() {
	r := gin.Default()
	// same as
	// config := cors.DefaultConfig()
	// config.AllowAllOrigins = true
	// router.Use(cors.New(config))
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
	r.POST("/order/:platform", func(c *gin.Context) {
		platform := c.Param("platform")
		if platform == "ebay" {
			var ebayOrder ebay
			err := c.ShouldBindJSON(&ebayOrder)
			if err == nil {
				log.Println(fmt.Sprintf("itemName: %v", ebayOrder.ItemName))
				log.Println(fmt.Sprintf("ImgUrl: %v", ebayOrder.ImgUrl))
				downloadFile(ebayOrder.ImgUrl, fmt.Sprintf("%s%s%d%s", "./", "ebay_", ebayOrder.Artikelnummer, ".jpg"))
			} else {
				log.Println("hmm irgendwas is fishy ", err.Error())
			}
			c.String(200, "Success")
		} else {
			c.String(404, "platform not supported")
		}
	})
	r.Run(":8081")
}

func store(db *sql.DB, order *ebay) {

	pathstring, _ := filepath.Abs("/mnt/d/purchase_history_sqlite.db")
	log.Printf("DB Pfad: %v", pathstring)
	db, err := sql.Open("sqlite3", pathstring)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}

func doDbStuff() {
	now := time.Now().UTC()
	log.Printf("Zeit: %v", now.Format("2006 01 02"))
	pathstring, _ := filepath.Abs("/mnt/d/purchase_history_sqlite.db")
	log.Printf("DB Pfad: %v", pathstring)
	db, err := sql.Open("sqlite3", pathstring)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var teststring string = "hallo"
	teststring2 := os.Args[1]
	fmt.Println(teststring)
	log.Printf("ergebnis: %v", teststring)
	log.Printf("ergebnis: %v", teststring2)
	var tsst TestType = TestType{"arsch", 2}
	tsst.test("wurst")

	rows, err := db.Query("select id, item_name, div from t_purchase")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name, div string
		err = rows.Scan(&id, &name, &div)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, ",DIV: ", div)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

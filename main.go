package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	cors "github.com/rs/cors/wrapper/gin"
)

//go reference time
//Mon Jan 2 15:04:05 -0700 MST 2006

type TestType struct {
	variable1 string
	variable2 int
}

type ebay struct {
	Id             int64
	PurchaseDate   string `jdon:"purchaseDate"`
	ItemName       string `json:"itemName"`
	VendorPlatform string
	Price          int    `json:"price"`
	Currency       string `json:"currency"`
	ImgUrl         string `json:"imgUrl"`
	imgFile        string
	//special ebay variables:
	Vendor        string `json:"vendor"`
	Artikelnummer int    `json:"artikelnummer"`
	imgHash       string
}

func (t TestType) test(kalr string) {
	log.Printf("variable1: %v %v", t.variable1, kalr)
	log.Printf("variable2: %v %v", t.variable2, kalr)
}

func main() {
	pathstring, _ := filepath.Abs("/mnt/d/orderHistory-sqlite.db")
	log.Printf("DB Pfad: %v", pathstring)
	db, err := sql.Open("sqlite3", pathstring)
	if err != nil {
		log.Fatal(err)
	}

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
			ebayOrder.VendorPlatform = "ebay"
			err := c.ShouldBindJSON(&ebayOrder)
			if err == nil {
				log.Println(fmt.Sprintf("itemName: %v", ebayOrder.ItemName))
				log.Println(fmt.Sprintf("ImgUrl: %v", ebayOrder.ImgUrl))
				fixedDate, err := time.Parse("02. Jan. 2006", ebayOrder.PurchaseDate)
				ebayOrder.PurchaseDate = fixedDate.Format("2006-01-02")
				if err != nil {
					log.Fatal(err)
				}
				ebayOrder.imgFile = fmt.Sprintf("%s%d", "ebay_", ebayOrder.Artikelnummer)
				var hash string
				ebayOrder.imgFile, hash, err = downloadFile(ebayOrder.ImgUrl, fmt.Sprintf("%s%s", "./img/", ebayOrder.imgFile))
				if err != nil {
					log.Printf("WARNUNG: Downloadfehler! %v", err)
				}
				log.Printf("image hash: %v\n", hash)
				ebayOrder.imgHash = hash
				store(db, &ebayOrder)
			} else {
				log.Println("hmm irgendwas is fishy ", err.Error())
			}
			c.JSON(200, gin.H{
				"message": "Success",
			})
		} else {
			c.JSON(404, gin.H{
				"message": "Fail (Platform not supported yet)",
			})
		}
	})
	r.Run(":8081")
}

func store(db *sql.DB, order *ebay) {

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into t_purchase(item_name, purchase_date, vendor_platform, price, img_url, div, img_file, currency, img_hash) values(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var divList []string
	divList = append(divList, strconv.Itoa(order.Artikelnummer))
	divList = append(divList, order.Vendor)
	jsondiv, err := json.Marshal(divList)
	div := string(jsondiv)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(order.ItemName, order.PurchaseDate, order.VendorPlatform, order.Price, order.ImgUrl, div, order.imgFile, order.Currency, order.imgHash)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func doDbStuff() {
	now := time.Now().UTC()
	log.Printf("Zeit: %v", now.Format("2006 01 02"))
	pathstring, _ := filepath.Abs("/mnt/d/orderHistory-sqlite.db")
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

func downloadFile(URL, fileName string) (string, string, error) {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", "", errors.New("Received non 200 response code")
	}
	//Create a empty file
	/*
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
	*/

	//Write the bytes to the file
	/*
		_, err = io.Copy(file, response.Body)
	*/
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}
	hsha2 := fmt.Sprintf("%x", sha256.Sum256(b))
	fmt.Println("SHA256: ", hsha2)
	fileName = fileName + "_" + hsha2[0:5] + ".jpg"
	if hsha2 != "a567462f4edd496bdf5cd00da5bbde64131c283e3cf396bfd58c0fac26b13d9a" {
		err = os.WriteFile(fileName, b, 0777)
	}
	return fileName, hsha2, nil
}

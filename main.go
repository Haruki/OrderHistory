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
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	webui "github.com/haruki/OrderHistory/ui_api"

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
	PurchaseDate   string `json:"purchaseDate"`
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

type aliexpress struct {
	Id            int64
	PurchaseDate  string
	ItemName      string
	Price         int
	SinglePrice   int
	ImgUrl        string
	imgFile       string
	imgHash       string
	Anzahl        int
	Vendor        string
	Bestellnummer int64
	ItemOption    string
	Currency      string
}

type alternate struct {
	Id           int64
	PurchaseDate string `json:"purchaseDate"`
	ItemName     string `json:"itemName"`
	Price        int    `json:"price"`
	ImgUrl       string `json:"imgUrl"`
	imgFile      string
	imgHash      string
	Anzahl       int `json:"anzahl"`
}

func (t TestType) test(kalr string) {
	log.Printf("variable1: %v %v", t.variable1, kalr)
	log.Printf("variable2: %v %v", t.variable2, kalr)
}

func main() {
	//pathstring, _ := filepath.Abs("/mnt/d/orderHistory-alternate-test.db")
	pathstring, _ := filepath.Abs("/mnt/d/20230101_orderHistory-sqlite.db")
	log.Printf("DB Pfad: %v", pathstring)
	db, err := sql.Open("sqlite3", pathstring)
	if err != nil {
		log.Fatal(err)
	}

	//r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New() //Default ersetzt durch New. Default hat einen debug logger, der nicht mehr benötigt wird.
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
			var aliOrder aliexpress
			err := c.ShouldBindJSON(&aliOrder)
			if err == nil {
				log.Println(fmt.Sprintf("itemName: %v", aliOrder.ItemName))
				log.Println(fmt.Sprintf("ImgUrl: %v", aliOrder.ImgUrl))
				fixedDate, err := time.Parse("02. Jan 2006", aliOrder.PurchaseDate)
				aliOrder.PurchaseDate = fixedDate.Format("2006-01-02")
				if err != nil {
					log.Fatal(err)
				}
				aliOrder.imgFile = fmt.Sprintf("%s%d", "aliexpress_", aliOrder.Bestellnummer)
				var hash string
				imgUrl, err := url.QueryUnescape(aliOrder.ImgUrl)
				if err != nil {
					log.Println("Fehler: URL kann nicht dekodiert werden!" + err.Error())
				}
				log.Println("unescaped: " + imgUrl)
				aliOrder.imgFile, hash, err = downloadFile(imgUrl, fmt.Sprintf("%s%s", "./img/", aliOrder.imgFile))
				if err != nil {
					log.Printf("WARNUNG: Downloadfehler! %v", err)
				}
				log.Printf("image hash: %v\n", hash)
				aliOrder.imgHash = hash
				storeAliexpress(db, &aliOrder)
			} else {
				log.Println("hmm irgendwas is fishy ", err.Error())
			}
			c.JSON(200, gin.H{
				"message": "Success",
			})
		} else if platform == "ebay" {
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
				storeEbay(db, &ebayOrder)
			} else {
				log.Println("hmm irgendwas is fishy ", err.Error())
			}
			c.JSON(200, gin.H{
				"message": "Success",
			})
		} else if platform == "alternate" {
			var alternateOrder alternate
			err := c.ShouldBindJSON(&alternateOrder)
			if err == nil {
				log.Println(fmt.Sprintf("itemName: %v", alternateOrder.ItemName))
				log.Println(fmt.Sprintf("ImgUrl: %v", alternateOrder.ImgUrl))
				fixedDate, err := time.Parse("02.01.2006", alternateOrder.PurchaseDate)
				alternateOrder.PurchaseDate = fixedDate.Format("2006-01-02")
				if err != nil {
					log.Fatal(err)
				}
				var hash string
				alternateOrder.imgFile = fmt.Sprintf("%s%s%s%s", "alternate_", alternateOrder.PurchaseDate, "_", SpaceMap(alternateOrder.ItemName)[:5])
				alternateOrder.imgFile, hash, err = downloadFile(alternateOrder.ImgUrl, fmt.Sprintf("%s%s", "./img/", alternateOrder.imgFile))
				if err != nil {
					log.Printf("WARNUNG: Downloadfehler! %v", err)
				}
				log.Printf("image hash: %v\n", hash)
				alternateOrder.imgHash = hash
				storeAlternate(db, &alternateOrder)
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
	log.Printf("Started OrderHistory-Server at Port: %d", 8081)
}

func storeAlternate(db *sql.DB, order *alternate) {

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
	divList = append(divList, strconv.Itoa(order.Anzahl))
	jsondiv, err := json.Marshal(divList)
	div := string(jsondiv)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(order.ItemName, order.PurchaseDate, "alternate", order.Price, order.ImgUrl, div, order.imgFile, "€", order.imgHash)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func storeEbay(db *sql.DB, order *ebay) {
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

func storeAliexpress(db *sql.DB, order *aliexpress) {
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
	divList = append(divList, order.ItemOption)
	divList = append(divList, order.Vendor)
	divList = append(divList, strconv.Itoa(order.SinglePrice))
	divList = append(divList, strconv.Itoa(order.Anzahl))
	jsondiv, err := json.Marshal(divList)
	div := string(jsondiv)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(order.ItemName, order.PurchaseDate, "aliexpress", order.Price, order.ImgUrl, div, order.imgFile, order.Currency, order.imgHash)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
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
	if hsha2 != "a567462f4edd496bdf5cd00da5bbde64131c283e3cf396bfd58c0fac26b13d9a" && hsha2 != "c041d4387a7d60b3d31d7f9c39e8ac531d8a342e24e695c739718a388f914f93" {
		err = os.WriteFile(fileName, b, 0777)
	}
	return fileName, hsha2, nil
}

func SpaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

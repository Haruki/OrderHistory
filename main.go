package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type TestType struct {
	variable1 string
	variable2 int
}

type ebay struct {
	purchaseDate   string
	ItemName       string `json:"itemName"`
	vendorPlatform string
	price          float32
	imgUrl         string
	//special ebay variables:
	vendor        string
	artikelnummer int
}

func (t TestType) test(kalr string) {
	log.Printf("variable1: %v %v", t.variable1, kalr)
	log.Printf("variable2: %v %v", t.variable2, kalr)
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/order", func(c *gin.Context) {
		var ebayOrder ebay
		err := c.ShouldBindJSON(&ebayOrder)
		if err == nil {
			log.Println(fmt.Sprintf("hmmm itemName: %v", ebayOrder.ItemName))
		} else {
			log.Println("hmm irgendwas is fishy ", err.Error())
		}
		c.String(200, "Success")
	})
	r.Run()
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

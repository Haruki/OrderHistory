package handler

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haruki/OrderHistory/db"
	structs "github.com/haruki/OrderHistory/struct"
)

func (h *Handler) Aliexpress(c *gin.Context) {

	var aliOrder structs.Aliexpress
	err := c.ShouldBindJSON(&aliOrder)
	if err == nil {
		log.Println(fmt.Sprintf("itemName: %v", aliOrder.ItemName))
		log.Println(fmt.Sprintf("ImgUrl: %v", aliOrder.ImgUrl))
		fixedDate, err := time.Parse("02. Jan 2006", aliOrder.PurchaseDate)
		aliOrder.PurchaseDate = fixedDate.Format("2006-01-02")
		if err != nil {
			log.Fatal(err)
		}
		aliOrder.ImgFile = fmt.Sprintf("%s%d", "aliexpress_", aliOrder.Bestellnummer)
		var hash string
		imgUrl, err := url.QueryUnescape(aliOrder.ImgUrl)
		if err != nil {
			log.Println("Fehler: URL kann nicht dekodiert werden!" + err.Error())
		}
		log.Println("unescaped: " + imgUrl)
		aliOrder.ImgFile, hash, err = downloadFile(imgUrl, fmt.Sprintf("%s%s", "./img/", aliOrder.ImgFile))
		if err != nil {
			log.Printf("WARNUNG: Downloadfehler! %v", err)
		}
		log.Printf("image hash: %v\n", hash)
		aliOrder.ImgHash = hash
		db.StoreAliexpress(h.db, &aliOrder)
	} else {
		log.Println("hmm irgendwas is fishy ", err.Error())
	}
	c.JSON(200, gin.H{
		"message": "Success",
	})
}

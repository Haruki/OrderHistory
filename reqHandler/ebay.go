package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haruki/OrderHistory/db"
	structs "github.com/haruki/OrderHistory/struct"
)

func (h *Handler) Ebay(c *gin.Context) {
	var ebayOrder structs.Ebay
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
		ebayOrder.ImgFile = fmt.Sprintf("%s%d", "ebay_", ebayOrder.Artikelnummer)
		var hash string
		ebayOrder.ImgFile, hash, err = downloadFile(ebayOrder.ImgUrl, fmt.Sprintf("%s%s", "./img/", ebayOrder.ImgFile))
		if err != nil {
			log.Printf("WARNUNG: Downloadfehler! %v", err)
		}
		log.Printf("image hash: %v\n", hash)
		ebayOrder.ImgHash = hash
		db.StoreEbay(h.db, &ebayOrder)
	} else {
		log.Println("hmm irgendwas is fishy ", err.Error())
	}
	c.JSON(200, gin.H{
		"message": "Success",
	})
}

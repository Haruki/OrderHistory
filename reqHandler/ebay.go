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
	ebayOrder.VendorPlatform = structs.String("ebay")
	err := c.ShouldBindJSON(&ebayOrder)
	if err == nil {
		//log.Println(fmt.Sprintf("itemName: %v", *ebayOrder.ItemName))
		// log.Println(fmt.Sprintf("ImgUrl: %v", *ebayOrder.ImgUrl))
		fixedDate, err := time.Parse("02. Jan. 2006", nvlString(ebayOrder.PurchaseDate))
		ebayOrder.PurchaseDate = structs.String(fixedDate.Format("2006-01-02"))
		if err != nil {
			log.Fatal(err)
		}
		ebayOrder.ImgFile = structs.String(fmt.Sprintf("%s%d", "ebay_", nvlInt(ebayOrder.EbaySpecial.Artikelnummer)))
		var hash string
		*ebayOrder.ImgFile, hash, err = downloadFile(*ebayOrder.ImgUrl, *ebayOrder.ImgFile)
		if err != nil {
			log.Printf("WARNUNG: Downloadfehler! %v", err)
		}
		// log.Printf("image hash: %v\n", hash)
		ebayOrder.ImgHash = structs.String(hash)
		db.StoreEbay(h.db, &ebayOrder)
	} else {
		log.Println("hmm irgendwas is fishy ", err.Error())
	}
	c.JSON(200, gin.H{
		"message": "Success",
	})
}

// nvl function
func nvlString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func nvlInt(s *int) int {
	if s == nil {
		return 0
	}
	return *s
}

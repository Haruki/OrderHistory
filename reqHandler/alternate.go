package handler

import (
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/haruki/OrderHistory/db"
	structs "github.com/haruki/OrderHistory/struct"
)

func (h *Handler) Alternate(c *gin.Context) {
	var alternateOrder structs.Alternate
	err := c.ShouldBindJSON(&alternateOrder)
	if err == nil {
		log.Printf("itemName: %v\n", alternateOrder.ItemName)
		log.Printf("ImgUrl: %v\n", alternateOrder.ImgUrl)
		fixedDate, err := time.Parse("02.01.2006", alternateOrder.PurchaseDate)
		alternateOrder.PurchaseDate = fixedDate.Format("2006-01-02")
		if err != nil {
			log.Fatal(err)
		}
		var hash string
		alternateOrder.ImgFile = fmt.Sprintf("%s%s%s%s", "alternate_", alternateOrder.PurchaseDate, "_", SpaceMap(alternateOrder.ItemName)[:5])
		alternateOrder.ImgFile, hash, err = downloadFile(alternateOrder.ImgUrl, alternateOrder.ImgFile)
		if err != nil {
			log.Printf("WARNUNG: Downloadfehler! %v", err)
		}
		log.Printf("image hash: %v\n", hash)
		alternateOrder.ImgHash = hash
		db.StoreAlternate(h.db, &alternateOrder)
	} else {
		log.Println("hmm irgendwas is fishy ", err.Error())
	}
	c.JSON(200, gin.H{
		"message": "Success",
	})

}

func SpaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"

	structs "github.com/haruki/OrderHistory/struct"
)

func StoreAlternate(db *sql.DB, order *structs.Alternate) {

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
	_, err = stmt.Exec(order.ItemName, order.PurchaseDate, "alternate", order.Price, order.ImgUrl, div, order.ImgFile, "€", order.ImgHash)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

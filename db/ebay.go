package db

import (
	"database/sql"
	"encoding/json"
	"log"

	structs "github.com/haruki/OrderHistory/struct"
)

func StoreEbay(db *sql.DB, order *structs.Ebay) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into t_purchase(item_name, purchase_date, vendor_platform, price, img_url, div, img_file, currency, img_hash) values(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	jsondiv, err := json.Marshal(order.EbaySpecial)
	div := string(jsondiv)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(order.ItemName, order.PurchaseDate, order.VendorPlatform, order.Price, order.ImgUrl, div, order.ImgFile, order.Currency, order.ImgHash)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

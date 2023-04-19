package ui

import (
	"database/sql"
	"log"
)

type Item struct {
	Name         string
	Vendor       string
	Price        int
	PurchaseDate string
	Id           int
}

func LoadAllItems(db *sql.DB) (error, []Item) {
	var result []Item
	var item Item
	rows, err := db.Query("select item_name, vendor_platform, id, purchase_date, price from t_purchase", 1)
	if err != nil {
		log.Fatal(err)
		return err, nil
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&item.Name, &item.Vendor, &item.Id, &item.PurchaseDate, &item.Price)
		if err != nil {
			log.Fatal(err)
			return err, nil
		}
		log.Println(item)
		result = append(result, item)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return err, nil
	}
	return nil, result
}

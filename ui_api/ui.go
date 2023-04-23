package ui

import (
	"database/sql"
	"log"
)

type Item struct {
	Name         string
	Vendor       string
	Price        *int `json:"Price,omitempty"`
	PurchaseDate string
	Id           int
	Currency     string
	Anzahl       *int `json:"Anzahl,omitempty"`
}

func LoadAllItems(db *sql.DB) (error, []Item) {
	var result []Item
	var item Item
	rows, err := db.Query("select item_name, vendor_platform, id, purchase_date, price, 1, currency from t_purchase", 1)
	if err != nil {
		log.Fatal(err)
		return err, nil
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&item.Name, &item.Vendor, &item.Id, &item.PurchaseDate, &item.Price, &item.Anzahl, &item.Currency)
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

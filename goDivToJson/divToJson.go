package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Div struct {
	Vendor string
	Div    string
}

func main() {
	divs, err := loadDivFromSqlLiteDB()
	if err != nil {
		log.Fatal(err)
	}
	for index, div := range divs {
		fmt.Println(div.Vendor+" "+div.Div, index)
	}
}

// Funktion zum Laden aller Zeilen mit dem Feld "div" in der Tabelle T_PURCHASE aus der SqlLite DB
func loadDivFromSqlLiteDB() ([]Div, error) {
	pathstring, _ := filepath.Abs("/mnt/d/20230101_orderHistory-sqlite.db")
	//pathstring, _ := filepath.Abs("d:/20230101_orderHistory-sqlite.db")
	log.Printf("DB Pfad: %v", pathstring)
	db, err := sql.Open("sqlite3", pathstring)
	if err != nil {
		log.Fatal(err)
	}

	var div Div
	var divArray []Div = make([]Div, 50)

	var query string = "select vendor_platform,div from t_purchase"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&div.Vendor, &div.Div)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		divArray = append(divArray, div)
		log.Println(div)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return divArray, nil
}

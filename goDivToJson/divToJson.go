package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Alternate struct {
	Anzahl int
}

type Ebay struct {
	Artikelnummer string
	Haendler      string
}

type Aliexpress struct {
	Option      *string `json:"Option,omitempty"`
	Haendler    *string `json:"Haendler,omitempty"`
	Einzelpreis *int    `json:"Einzelpreis,omitempty"`
	Anzahl      *int    `json:"Anzahl,omitempty"`
}

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
		fmt.Printf("current values: Vendor: %s, div: %s, index: %d\n", div.Vendor, div.Div, index)
		update(&divs[index].Div, div.Vendor)
	}
}

func update(div *string, vendor string) {
	tD := *div
	var trimmedDivArray []string
	fmt.Println("test: " + tD)
	err := json.Unmarshal([]byte(tD), &trimmedDivArray)
	if err != nil {
		log.Fatalf("Error unmarshalling div: %v", err)
		log.Fatal(err)
	}
	var divJsonString string
	switch vendor {
	case "alternate":
		divJsonString = updateAlternate(trimmedDivArray)
	case "ebay":
		divJsonString = updateEbay(trimmedDivArray)
	case "aliexpress":
		divJsonString = updateAliexpress(trimmedDivArray)
	}
	div = &divJsonString
}

func updateAliexpress(trimmedDivArray []string) string {
	var aliexpress Aliexpress
	var err error
	*aliexpress.Option = trimmedDivArray[0][1 : len(trimmedDivArray[0])-1]
	*aliexpress.Haendler = trimmedDivArray[1][1 : len(trimmedDivArray[1])-1]
	*aliexpress.Einzelpreis, err = strconv.Atoi(trimmedDivArray[2][1 : len(trimmedDivArray[2])-1])
	if err != nil {
		log.Fatal(err)
	}
	*aliexpress.Anzahl, err = strconv.Atoi(trimmedDivArray[3][1 : len(trimmedDivArray[3])-1])
	if err != nil {
		log.Fatal(err)
	}
	divJsonBytes, err := json.Marshal(aliexpress)
	if err != nil {
		log.Fatal(err)
	}
	var divJsonString string = string(divJsonBytes)
	return divJsonString
}

func updateEbay(trimmedDivArray []string) string {
	var ebay Ebay
	var err error
	ebay.Artikelnummer = trimmedDivArray[0][1 : len(trimmedDivArray[0])-1]
	ebay.Haendler = trimmedDivArray[1][1 : len(trimmedDivArray[1])-1]
	divJsonBytes, err := json.Marshal(ebay)
	if err != nil {
		log.Fatal(err)
	}
	divJsonString := string(divJsonBytes)
	return divJsonString
}

func updateAlternate(trimmedDivArray []string) string {
	var alternate Alternate
	var err error
	anzahl, err := strconv.Atoi(trimmedDivArray[0][1 : len(trimmedDivArray[1])-1])
	if err != nil {
		log.Fatal(err)
	}
	alternate.Anzahl = anzahl
	divJsonBytes, err := json.Marshal(alternate)
	divJsonString := string(divJsonBytes)
	return divJsonString
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
	var divArray []Div = make([]Div, 0)

	var query string = "select vendor_platform,div from t_purchase order by id desc"
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
		//log.Println(div)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Printf("length of array: %d\n", len(divArray))
	return divArray, nil
}

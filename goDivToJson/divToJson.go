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
	Anzahl *int `json:"Anzahl,omitempty"`
}

type Ebay struct {
	Artikelnummer *string `json:"Artikelnummer,omitempty"`
	Haendler      *string `json:"Haendler,omitempty"`
}

type Aliexpress struct {
	Option      *string `json:",omitempty"`
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
	fmt.Printf("Result: %s\n", divJsonString)
	div = &divJsonString
}

func updateAliexpress(trimmedDivArray []string) string {
	var aliexpress Aliexpress
	var err error
	//Option
	if &trimmedDivArray[0] != nil && len(trimmedDivArray[0]) > 2 {
		aliexpress.Option = &trimmedDivArray[0]
	} else {
		aliexpress.Option = nil
	}
	//Haendler
	if &trimmedDivArray[1] != nil && len(trimmedDivArray[1]) > 2 {
		aliexpress.Haendler = &trimmedDivArray[1]
	} else {
		aliexpress.Haendler = nil
	}
	//Einzelpreis
	var einzelpreisInt int
	if &trimmedDivArray[2] != nil && len(trimmedDivArray[2]) > 0 {
		einzelpreisInt, err = strconv.Atoi(trimmedDivArray[2])
		if err != nil {
			log.Fatal(err)
		}
		if einzelpreisInt > 0 {
			aliexpress.Einzelpreis = &einzelpreisInt
		} else {
			aliexpress.Einzelpreis = nil
		}
	} else {
		aliexpress.Einzelpreis = nil
	}
	//Anzahl
	var anzahlInt int
	if &trimmedDivArray[3] != nil && len(trimmedDivArray[3]) > 0 {
		anzahlInt, err = strconv.Atoi(trimmedDivArray[3])
		if err != nil {
			log.Fatal(err)
		}
		if anzahlInt > 0 {
			aliexpress.Anzahl = &anzahlInt
		} else {
			aliexpress.Anzahl = nil
		}
	} else {
		aliexpress.Anzahl = nil
	}
	//Write to json
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
	//Artikelnummer
	if &trimmedDivArray[0] != nil && len(trimmedDivArray[0]) > 2 {
		ebay.Artikelnummer = &trimmedDivArray[0]
	} else {
		ebay.Artikelnummer = nil
	}
	//Haendler
	if &trimmedDivArray[1] != nil && len(trimmedDivArray[1]) > 2 {
		ebay.Haendler = &trimmedDivArray[1]
	} else {
		ebay.Haendler = nil
	}
	//Write to json
	divJsonBytes, err := json.Marshal(ebay)
	if err != nil {
		log.Fatal(err)
	}
	divJsonString := string(divJsonBytes)
	return divJsonString
}

func updateAlternate(trimmedDivArray []string) string {
	var alternate Alternate
	//Anzahl
	if &trimmedDivArray[0] != nil && len(trimmedDivArray[0]) > 0 {
		anzahl, err := strconv.Atoi(trimmedDivArray[0])
		if err != nil {
			log.Fatal(err)
		}
		if anzahl > 0 {
			alternate.Anzahl = &anzahl
		} else {
			alternate.Anzahl = nil
		}
	} else {
		alternate.Anzahl = nil
	}
	//Write to json
	divJsonBytes, err := json.Marshal(alternate)
	if err != nil {
		log.Fatal(err)
	}
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

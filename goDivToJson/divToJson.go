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
	Id     int
}

func main() {
	pathstring, _ := filepath.Abs("/mnt/d/20230101_orderHistory-sqlite.db")
	//pathstring, _ := filepath.Abs("d:/20230101_orderHistory-sqlite.db")
	log.Printf("DB Pfad: %v", pathstring)
	db, err := sql.Open("sqlite3", pathstring)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	divs, err := loadDivFromSqlLiteDB(db)
	if err != nil {
		log.Fatal(err)
	}
	for index, div := range divs {
		//fmt.Printf("current values: Id: %d, Vendor: %s, div: %s, index: %d\n", div.Id, div.Vendor, div.Div, index)
		divs[index].Div = update(divs[index].Div, div.Vendor)
	}
	// for _, div := range divs {
	// 	fmt.Printf("Id: %d, Vendor: %s, div: %s\n", div.Id, div.Vendor, div.Div)
	// }
	updateDivInSqlLiteDB(divs, db)
	fmt.Println("Update complete")
}

func update(div string, vendor string) string {
	tD := div
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
	return divJsonString
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

// Funktion zum Updaten aller Zeilen mit dem Feld "div" in der Tablle T_PURCHASE mit einem neuen Wert für "div"
func updateDivInSqlLiteDB(divs []Div, db *sql.DB) error {
	for _, div := range divs {
		statement := "update t_purchase set div = ? where id = ?"
		stmt, err := db.Prepare(statement)
		if err != nil {
			log.Fatal(err)
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(div.Div, div.Id)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

// Funktion zum Laden aller Zeilen mit dem Feld "div" in der Tabelle T_PURCHASE aus der SqlLite DB
func loadDivFromSqlLiteDB(db *sql.DB) ([]Div, error) {

	var div Div
	var divArray []Div = make([]Div, 0)

	var query string = "select id, vendor_platform,div from t_purchase order by id desc"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&div.Id, &div.Vendor, &div.Div)
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

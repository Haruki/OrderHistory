package db

import (
	"database/sql"
	"log"
	"time"
)

func UpdateImage(db *sql.DB, fileName string, sha2 string, id int) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare("update t_purchase set img_file=?, img_hash=? , img_url=? where id=?")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	//set update query parameters and execute
	_, err = stmt.Exec(fileName, sha2, fileName, id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	tx.Commit()
	return nil
}

// Checks if an item is already in the database.
func ItemExists(db *sql.DB, itemName string, date string, vendor string) (bool, error) {
	var count int
	date = convertDate(date, vendor)
	err := db.QueryRow("select count(*) from t_purchase where item_name=? and purchase_date=? and vendor_platform=?", itemName, date, vendor).Scan(&count)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// Converts date formats depending on vendor.
func convertDate(date string, vendor string) string {
	var result string
	if vendor == "ebay" {
		fixedDate, err := time.Parse("02. Jan. 2006", date)
		result = fixedDate.Format("2006-01-02")
		if err != nil {
			log.Fatal(err)
		}
	}
	return result
}

func InsertNewItemManual(db *sql.DB, itemName string, date string, price int, currency string, vendor string, div string) error {
	date = convertDate(date, vendor)
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare("insert into t_purchase(item_name, purchase_date, price, currency, vendor_platform) values(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	//set insert query parameters and execute
	_, err = stmt.Exec(itemName, date, price, currency, vendor)
	if err != nil {
		log.Fatal(err)
		return err
	}
	tx.Commit()
	return nil
}

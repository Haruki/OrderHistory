package db

import (
	"database/sql"
	"log"
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

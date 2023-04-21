package handler

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func downloadFile(URL, fileName string) (string, string, error) {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", "", errors.New("Received non 200 response code")
	}
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}
	hsha2 := fmt.Sprintf("%x", sha256.Sum256(b))
	fmt.Println("SHA256: ", hsha2)
	fileName = fileName + "_" + hsha2[0:5] + ".jpg"
	if hsha2 != "a567462f4edd496bdf5cd00da5bbde64131c283e3cf396bfd58c0fac26b13d9a" && hsha2 != "c041d4387a7d60b3d31d7f9c39e8ac531d8a342e24e695c739718a388f914f93" {
		err = os.WriteFile(fileName, b, 0777)
	}
	return fileName, hsha2, nil
}

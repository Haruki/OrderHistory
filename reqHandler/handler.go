package handler

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	orderHistoryDb "github.com/haruki/OrderHistory/db"

	"github.com/gin-gonic/gin"
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

func (h *Handler) CheckItemExists(c *gin.Context) {
	itemName := c.Query("itemName")
	purchaseDate := c.Query("purchaseDate")
	vendor := c.Query("vendor")
	exists, err := orderHistoryDb.ItemExists(h.db, itemName, purchaseDate, vendor)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error while checking DB",
		})
		return
	}
	if exists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"exists": exists})
	} else {
		c.JSON(http.StatusOK, gin.H{"exists": exists})
	}
}

func (h *Handler) UpdateImage(c *gin.Context) {
	id := c.PostForm("itemId")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No valid id",
		})
		return
	}
	vendor := c.PostForm("vendor")
	file, err := c.FormFile("file")
	// The file cannot be received.
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}
	sha256Hash, err := hashImage(file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error while hashing file",
		})
		return
	}
	//get filename extension (returns string including the dot)
	ext := filepath.Ext(file.Filename)
	newDbFileName := fmt.Sprintf("./img/%s_%s_%s%s", vendor, id, sha256Hash[:5], ext)
	newFilePathName := fmt.Sprintf("./img/backup/%s_%s_%s%s", vendor, id, sha256Hash[:5], ext)
	err = c.SaveUploadedFile(file, newFilePathName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error while saving file",
		})
		return
	}
	err = orderHistoryDb.UpdateImage(h.db, newDbFileName, sha256Hash, intId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error while updating DB",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"imgFile": newDbFileName})
}

func hashImage(file *multipart.FileHeader) (string, error) {
	realfile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer realfile.Close()
	fileBytes, err := ioutil.ReadAll(realfile)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	hash.Write(fileBytes)
	sha256Hash := fmt.Sprintf("%x", hash.Sum(nil))
	return sha256Hash, nil
}

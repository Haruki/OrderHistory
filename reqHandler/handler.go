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
	"regexp"
	"strconv"
	"strings"
	"time"

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
		return "", "", errors.New("received non 200 response code")
	}
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}
	contentType := response.Header.Get("Content-type")
	//letzte drei Zeichen des 'Content-Type' für Ermittlung des Typs (jpg/png/gif/...)
	if len(contentType) >= 3 {
		contentType = contentType[len(contentType)-3:]
	} else {
		contentType = "jpg"
	}
	hsha2 := fmt.Sprintf("%x", sha256.Sum256(b))
	fmt.Println("SHA256: ", hsha2)
	fileName = fmt.Sprintf("./img/%s_%s.%s", fileName, hsha2[0:5], contentType)
	if hsha2 != "a567462f4edd496bdf5cd00da5bbde64131c283e3cf396bfd58c0fac26b13d9a" && hsha2 != "c041d4387a7d60b3d31d7f9c39e8ac531d8a342e24e695c739718a388f914f93" {
		err = os.WriteFile(fileName, b, 0777)
		if err != nil {
			return "", "", err
		}
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
	if !exists {
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
	// newFilePathName := fmt.Sprintf("./img/backup/%s_%s_%s%s", vendor, id, sha256Hash[:5], ext)
	err = c.SaveUploadedFile(file, newDbFileName)
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

func (h *Handler) AddNewItemManual(c *gin.Context) {
	itemName := c.PostForm("itemName")
	date := c.PostForm("date")
	if date == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No valid date",
		})
	}
	fixedDate, err := time.Parse("02.01.2006", date)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No valid date",
		})
	}
	date = fixedDate.Format("2006-01-02")
	price := c.PostForm("price")
	price = strings.ReplaceAll(price, ",", "")
	price = strings.ReplaceAll(price, ".", "")
	if match, err := regexp.MatchString("^[0-9]+$", price); err != nil || !match {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid price",
		})
	}

	imgUrl := c.PostForm("imgUrl")
	vendor := c.PostForm("platform")
	var imgFileName, sha256Hash string
	if imgUrl != "" {
		imgFileName, sha256Hash, err = downloadFile(imgUrl, vendor)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Error while downloading file",
			})
			return
		}
	}
	intPrice, err := strconv.Atoi(price)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No valid price",
		})
		return
	}
	currency := c.PostForm("currency")
	div := c.PostForm("div")
	err = orderHistoryDb.InsertNewItemManual(h.db, itemName, date, intPrice, currency, vendor, div, imgFileName, sha256Hash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error while inserting into DB",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item added"})
}

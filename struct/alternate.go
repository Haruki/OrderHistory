package structs

type Alternate struct {
	Id           int64
	PurchaseDate string `json:"purchaseDate"`
	ItemName     string `json:"itemName"`
	Price        int    `json:"price"`
	ImgUrl       string `json:"imgUrl"`
	ImgFile      string
	ImgHash      string
	Anzahl       int `json:"anzahl"`
}

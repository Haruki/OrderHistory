package structs

type Ebay struct {
	Id             int64
	PurchaseDate   string `json:"purchaseDate"`
	ItemName       string `json:"itemName"`
	VendorPlatform string
	Price          int    `json:"price"`
	Currency       string `json:"currency"`
	ImgUrl         string `json:"imgUrl"`
	ImgFile        string
	//special ebay variables:
	Vendor        string `json:"vendor"`
	Artikelnummer int    `json:"artikelnummer"`
	ImgHash       string
}

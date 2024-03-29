package structs

type Ebay struct {
	Id             *int64  `json:"id,omitempty"`
	PurchaseDate   *string `json:"purchaseDate,omitempty"`
	ItemName       *string `json:"itemName,omitempty"`
	VendorPlatform *string
	Price          *int        `json:"price,omitempty"`
	Currency       *string     `json:"currency,omitempty"`
	ImgUrl         *string     `json:"imgUrl,omitempty"`
	ImgFile        *string     `json:"imgFile,omitempty"`
	ImgHash        *string     `json:"imgHash,omitempty"`
	EbaySpecial    EbaySpecial `json:"ebaySpecial,omitempty"`
}

type EbaySpecial struct {
	//special ebay variables:
	Haendler      *string `json:"Haendler,omitempty"`
	Artikelnummer *int    `json:"Artikelnummer,omitempty"`
	ImgHash       *string `json:"imgHash,omitempty"`
}

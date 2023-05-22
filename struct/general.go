package structs

func String(s string) *string {
	return &s
}

func Int(i int) *int {
	return &i
}

type GenericItem struct {
	Name         string
	Vendor       string
	Price        *int `json:"Price,omitempty"`
	PurchaseDate string
	Currency     string
	ImgFile      string
}

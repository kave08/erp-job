package models

import "time"

type Fararavand struct {
	BaseData      BaseData
	Customers     Customers
	Products      Products
	Invoices      Invoices
	Treasuries    Treasuries
	InvoiceReturn InvoiceReturn
}

type BaseData struct {
	PaymentTypes []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"paymentTypes"`
	CustomerTypes []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"customerTypes"`
	GiuldTypes []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"giuldTypes"`
	SanadTypes []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"sanadTypes"`
	Units []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"units"`
	Branches []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"branches"`
	Brands []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"brands"`
	Areas []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"areas"`
	Districts []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"districts"`
	States []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"states"`
	Cities []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"cities"`
	WareHouses []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"wareHouses"`
	FirstProductGroup []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"firstProductGroup"`
	SecondProductGroup []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
}

type Customers struct {
	ID             int    `json:"id"`
	BranchID       int    `json:"branchId"`
	BranchName     string `json:"branchName"`
	Code           int    `json:"code"`
	Name           string `json:"name"`
	NameTablo      string `json:"nameTablo"`
	CodeEghtesady  string `json:"codeEghtesady"`
	StateID        int    `json:"stateId"`
	CityID         int    `json:"cityId"`
	AreaID         int    `json:"areaId"`
	DistrictID     int    `json:"districtId"`
	CustomerTypeID int    `json:"customerTypeId"`
	GuildTypeID    int    `json:"guildTypeId"`
	PaymentTypeID  int    `json:"paymentTypeId"`
	Status         int    `json:"status"`
	CodeTafsily1   int    `json:"codeTafsily1"`
	Address        string `json:"address"`
	Telephone      string `json:"telephone"`
	Mobile         string `json:"mobile"`
}

type Products struct {
	ID                   int     `json:"id"`
	Name                 string  `json:"name"`
	Code                 string  `json:"code"`
	UnitID               int     `json:"unitId"`
	Tol                  float64 `json:"tol"`
	Arz                  float64 `json:"arz"`
	Ertefa               float64 `json:"ertefa"`
	VaznKhales           float64 `json:"vaznKhales"`
	VaznNaKhales         float64 `json:"vaznNaKhales"`
	VaznKarton           float64 `json:"vaznKarton"`
	FirstProductGroupID  int     `json:"firstProductGroupId"`
	SecondProductGroupID int     `json:"secondProductGroupId"`
	Status               int     `json:"status"`
	BrandID              int     `json:"brandId"`
	TedadDarBasteh       int     `json:"tedadDarBasteh"`
	TedadDarKarton       int     `json:"tedadDarKarton"`
	HajmKala             float64 `json:"hajmKala"`
	SupplierID           int     `json:"supplierId"`
	HasAvarez            bool    `json:"hasAvarez"`
	HasMalyat            bool    `json:"hasMalyat"`
}

type Invoices struct {
	BranchID                  int       `json:"branchId"`
	CodeDoreh                 int       `json:"codeDoreh"`
	InvoiceId                 int       `json:"invoiceId"`
	InvoiceDate               string    `json:"invoiceDate"` //"14021027"
	InvoiceNumber             int       `json:"invoiceNumber"`
	CustomerID                int       `json:"customerId"`
	PaymentTypeID             int       `json:"paymentTypeId"`
	ModatCheck                int       `json:"modatCheck"`
	InvoicePriceKhales        int       `json:"invoicePriceKhales"`
	SumJayezeh                int       `json:"sumJayezeh"`
	SumMalyat                 int       `json:"sumMalyat"`
	SumAvarez                 int       `json:"sumAvarez"`
	InvoicePrice              int       `json:"invoicePrice"`
	Date                      time.Time `json:"date"`
	WareHouseID               int       `json:"wareHouseId"`
	ProductID                 int       `json:"productId"`
	ProductCount              int       `json:"productCount"`
	ProductFee                int       `json:"productFee"`
	ProductPrice              int       `json:"productPrice"`
	DiscountPercentage        int       `json:"discountPercentage"`
	ProductDiscount           int       `json:"productDiscount"`
	IsJayezeh                 int       `json:"isJayezeh"`
	VisitorName               string    `json:"visitorName"`
	VisitorAddress            string    `json:"visitorAddress"`
	VisitorTelephone          string    `json:"visitorTelephone"`
	VisitorCodeMely           string    `json:"visitorCodeMely"`
	VisitorCodePosty          string    `json:"visitorCodeposty"`
	VisitorCode               string    `json:"visitorCode"`
	CustomerName              string    `json:"customerName"`
	CustomerAddress           string    `json:"customerAddress"`
	ShenasehMeli              int       `json:"shenasehMeli"`
	ShomarehSabt              int       `json:"shomarehSabt"`
	CustomerCodePosti         string    `json:"customerCodePosti"`
	CustomerCodeMeli          string    `json:"customerCodeMeli"`
	Codekala                  string    `json:"codekala"`
	BarCode                   string    `json:"barCode"`
	NameKalaFaktor            string    `json:"nameKalaFaktor"`
	ProductPriceAfterDiscount int       `json:"productPriceAfterDiscount"`
	MaliatAvarez              int       `json:"malyatAvarez"`
	ProductPriceNet           int       `json:"productPriceNet"`
	SNoePardakht              int       `json:"sNoePardakht"`
	CCForoshandeh             int       `json:"CCForoshandeh"`
	CodeForoshandeh           int       `json:"codeForoshandeh"`
	CodeMahal                 int       `json:"codeMahal"`
	TozihatFaktor             string    `json:"tozihatFaktor"`
	NameAnbar                 string    `json:"nameAnbar"`
	TxtNoePardakht            string    `json:"txtNoePardakht"`
}

type Treasuries struct {
	CodeDoreh                 int    `json:"codeDoreh"`
	BranchID                  int    `json:"branchId"`
	SanadTypeID               int    `json:"sanadTypeId"`
	InvoiceID                 int    `json:"invoiceId"`
	MablaghKolDaryafti        int    `json:"mablaghKolDaryafti"`
	MablaghBasteShodeBeFaktor int    `json:"mablaghBasteShodeBeFaktor"`
	ShomarehHesabID           int    `json:"shomarehHesabId"`
	CodeGoroh                 string `json:"codeGoroh"`
	CodeKol                   string `json:"codeKol"`
	CodeMoeen                 string `json:"codeMoeen"`
	Tafsily1                  int    `json:"tafsily1"`
	TarikhDaryaft             string `json:"tarikhDaryaft"`  //"14021025"
	TarikhSarResid            string `json:"tarikhSarResid"` //"14021025"
	ShomarehSanad             string `json:"shomarehSanad"`
}

type InvoiceReturn struct {
	BranchID                 int    `json:"branchId"`
	CodeDoreh                int    `json:"codeDoreh"`
	InvoiceID                int    `json:"invoiceId"`
	ReturnDate               string `json:"returnDate"` //"14021027"
	ReturnNumber             int    `json:"returnNumber"`
	CustomerID               int    `json:"customerId"`
	PriceKhales              int    `json:"priceKhales"`
	ReturnMaliatAvarezPrice  int    `json:"returnMaliatAvarezPrice"`
	Price                    int    `json:"price"`
	ProductID                int    `json:"productId"`
	ProductCount             int    `json:"productCount"`
	ProductFee               int    `json:"productFee"`
	ProductPrice             int    `json:"productPrice"`
	IsReturnJayezeh          bool   `json:"isReturnJayezeh"`
	ProductMaliatAvarezPrice int    `json:"productMaliatAvarezPrice"`
}

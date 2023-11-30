package models

type Fararavand struct {
	BaseData   BaseData
	Customers  Customers
	Products   Products
	Invoices   Invoices
	Treasuries Treasuries
	Reverted   Reverted
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
	ID                int    `json:"id"`
	BranchID          int    `json:"branchId"`
	BranchName        string `json:"branchName"`
	Code              int    `json:"code"`
	Name              string `json:"name"`
	NameTablo         string `json:"nameTablo"`
	CodeEghtesady     string `json:"codeEghtesady"`
	StateID           int    `json:"stateId"`
	CityID            int    `json:"cityId"`
	AreaID            int    `json:"areaId"`
	DistrictID        int    `json:"districtId"`
	CustomerTypeID    int    `json:"customerTypeId"`
	GuildTypeID       int    `json:"guildTypeId"`
	PaymentTypeID     int    `json:"paymentTypeId"`
	Status            int    `json:"status"`
	CodeTafsily1      int    `json:"codeTafsily1"`
	Address           string `json:"address"`
	Telephone         string `json:"telephone"`
	Mobile            string `json:"mobile"`
	CustomerId        int    `json:"customer_id"`
	CustomerName      int    `json:"customer_name"`
	CustomerAddress   int    `json:"customer_address"`
	ShenasehMeli      int    `json:"shenaseh_meli"`
	ShomarehSabt      int    `json:"shomareh_sabt"`
	CustomerCodePosti int    `json:"customer_code_posti"`
	CustomerCodeMeli  int    `json:"customer_code_meli"`
}

type Products struct {
	ID                        int     `json:"id"`
	Name                      string  `json:"name"`
	Code                      string  `json:"code"`
	UnitID                    int     `json:"unit_id"`
	Tol                       float64 `json:"tol"`
	Arz                       float64 `json:"arz"`
	Ertefa                    float64 `json:"ertefa"`
	VaznKhales                float64 `json:"vazn_khales"`
	VaznNaKhales              float64 `json:"vazn_na_khales"`
	VaznKarton                float64 `json:"vazn_karton"`
	FirstProductGroupID       int     `json:"first_product_group_id"`
	SecondProductGroupID      int     `json:"secondProductGroupId"`
	Status                    int     `json:"status"`
	BrandID                   int     `json:"brandId"`
	TedadDarBasteh            int     `json:"tedadDarBasteh"`
	TedadDarKarton            int     `json:"tedadDarKarton"`
	HajmKala                  float64 `json:"hajmKala"`
	SupplierID                int     `json:"supplierId"`
	HasAvarez                 bool    `json:"hasAvarez"`
	HasMalyat                 bool    `json:"hasMalyat"`
	WareHouseId               int     `json:"wareHouse_id"`
	ProductId                 int     `json:"product_id"`
	ProductCount              int     `json:"product_count"`
	ProductFee                int     `json:"product_fee"`
	ProductPrice              int     `json:"product_price"`
	ProductDiscount           int     `json:"product_discount"`
	IsJayezeh                 int     `json:"is_jayezeh"`
	Codekala                  int     `json:"code_kala"`
	BarCode                   int     `json:"bar_code"`
	NameKalaFaktor            int     `json:"name_kala_faktor"`
	ProductPriceAfterDiscount int     `json:"product_price_after_discount"`
	MalyatAvarez              int     `json:"malyat_avarez"`
	ProductPriceNetz          int     `json:"product_price_netz"`
}

type Invoices struct {
	BranchID           int    `json:"branch_id"`
	CodeDoreh          int    `json:"code_doreh"`
	InvoiceId          int    `json:"invoice_id"`
	InvoiceDate        string `json:"invoice_date"`
	InvoiceNumber      int    `json:"invoice_number"`
	CustomerID         int    `json:"customer_id"`
	PaymentTypeID      int    `json:"payment_type_id"`
	ModatCheck         int    `json:"modat_check"`
	InvoicePriceKhales int    `json:"invoice_price_khales"`
	SumJayezeh         int    `json:"sum_jayezeh"`
	SumMalyat          int    `json:"sum_malyat"`
	SumAvarez          int    `json:"sum_avarez"`
	InvoicePrice       int    `json:"invoice_price"`
	Date               string `json:"date"`
	WareHouseID        int    `json:"ware_houseId"`
	ProductID          int    `json:"product_id"`
	ProductCount       int    `json:"product_count"`
	ProductFee         int    `json:"product_fee"`
	ProductPrice       int    `json:"product_price"`
	DiscountPercentage int    `json:"discount_percentage"`
	ProductDiscount    int    `json:"product_discount"`
	IsJayezeh          int    `json:"is_jayezeh"`
	VisitorAddress     int    `json:"visitor_address"`
	VisitorCode        int    `json:"visitor_code"`
	VisitorCodeMely    int    `json:"visitor_code_mely"`
	VisitorCodePosty   int    `json:"visitor_code_posty"`
	VisitorName        int    `json:"visitor_name"`
	VisitorTelephone   int    `json:"visitor_telephone"`
	CustomerCodePosti  int    `json:"customer_code_posti"`
}

type Treasuries struct {
	CodeDoreh                 int         `json:"codeDoreh"`
	BranchID                  int         `json:"branchId"`
	SanadTypeID               int         `json:"sanadTypeId"`
	InvoiceID                 int         `json:"invoiceId"`
	MablaghKolDaryafti        int         `json:"mablaghKolDaryafti"`
	MablaghBasteShodeBeFaktor int         `json:"mablaghBasteShodeBeFaktor"`
	ShomarehHesabID           interface{} `json:"shomarehHesabId"`
	CodeGoroh                 string      `json:"codeGoroh"`
	CodeKol                   string      `json:"codeKol"`
	CodeMoeen                 string      `json:"codeMoeen"`
	Tafsily1                  int         `json:"tafsily1"`
	TarikhDaryaft             string      `json:"tarikhDaryaft"`
	TarikhSarResid            string      `json:"tarikhSarResid"`
	ShomarehSanad             string      `json:"shomarehSanad"`
}

type Reverted struct {
	BranchID                 int    `json:"branchId"`
	CodeDoreh                int    `json:"codeDoreh"`
	InvoiceID                int    `json:"invoiceId"`
	ReturnDate               string `json:"returnDate"`
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

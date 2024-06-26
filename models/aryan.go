package models

import "time"

type SaleOrder struct {
	CustomerId       int     `json:"customer_id"`
	VoucherDate      string  `json:"voucher_date"`
	SecondNumber     int     `json:"second_number"`
	VoucherDesc      string  `json:"voucher_desc"`
	StockID          int     `json:"stock_id"`
	SaleTypeId       int     `json:"sale_type_id"`
	DeliveryCenterID int     `json:"delivery_center_id"`
	SaleCenterID     int     `json:"sale_center_id"`
	PaymentWayID     int     `json:"payment_way_id"`
	SellerVisitorID  int     `json:"seller_visitor_id"`
	ServiceGoodsID   int     `json:"service_goods_id"`
	Quantity         float64 `json:"quantity"`
	Fee              float64 `json:"fee"`
	DetailDesc       string  `json:"detail_desc"`
}

type SaleCustomer struct {
	CustomerID   int    `json:"customer_id"`
	CustomerCode string `json:"customer_code"`
}

type SaleTypeSelect struct {
	BuySaleTypeID   int    `json:"buy_sale_type_id"`
	BuySaleTypeCode string `json:"buy_sale_type_code"`
	BuySaleTypeDesc string `json:"buy_sale_type_desc"`
}

type SaleCenter4SaleSelect struct {
	StockID   int    `json:"stock_id"`
	StockCode string `json:"stock_code"`
	StockDesc string `json:"stock_desc"`
}
type SalePaymentSelect struct {
	PaymentWayID   int    `json:"payment_way_id"`
	PaymentwayDesc string `json:"payment_way_desc"`
}

type SaleCenterSelect struct {
	CentersID   int    `json:"centers_id"`
	CentersCode string `json:"centers_code"`
	CenterDesc  string `json:"center_desc"`
}

type DeliverCenter_SaleSelect struct {
	CentersID   int    `json:"centers_id"`
	CentersCode string `json:"centers_code"`
}

type SalerSelect struct {
	SaleVisitorID   int    `json:"sale_visitor_id"`
	SaleVisitorDesc string `json:"sale_visitor_desc"`
}

type SaleSellerVisitor struct {
	CentersID   int    `json:"centers_id"`
	CentersCode string `json:"centers_code"`
}

type Goods struct {
	ServiceGoodsID   int    `json:"service_goods_id"`
	ServiceGoodsCode string `json:"service_goods_code"`
	ServiceGoodsDesc string `json:"service_goods_desc"`
	GroupId          int    `json:"group_id"`
	TypeID           int    `json:"type_id"`
	SecUnitType      int    `json:"sec_unit_type"`
	Level1           int    `json:"level1"`
	Level2           int    `json:"level2"`
	Level3           int    `json:"level3"`
}

type SaleProforma struct {
	CustomerId       int       `json:"customer_id"`
	VoucherDate      time.Time `json:"voucher_date"`
	SecondNumber     string    `json:"second_number"`
	VoucherDesc      string    `json:"voucher_desc"`
	StockID          int       `json:"stock_id"`
	SaleTypeId       int       `json:"sale_type_id"`
	DeliveryCenterID int       `json:"delivery_center_id"`
	SaleCenterID     int       `json:"sale_center_id"`
	PaymentWayID     int       `json:"payment_way_id"`
	SellerVisitorID  int       `json:"seller_visitor_id"`
	ServiceGoodsID   int       `json:"service_goods_id"`
	Quantity         float64   `json:"quantity"`
	Fee              float64   `json:"fee"`
	DetailDesc       string    `json:"detail_desc"`
}

type SaleFactor struct {
	CustomerID       int     `json:"CustomerId"`
	VoucherDate      string  `json:"VoucherDate"`
	StockID          int     `json:"StockID"`
	VoucherDesc      string  `json:"VoucherDesc"`
	SaleTypeID       int     `json:"SaleTypeId"`
	DeliveryCenterID int     `json:"DeliveryCenterID"`
	SaleCenterID     int     `json:"SaleCenterID"`
	PaymentWayID     int     `json:"PaymentWayID"`
	SellerID         int     `json:"SellerID"`
	SaleManID        int     `json:"SaleManID"`
	DistributerID    int     `json:"DistributerId"`
	SecondNumber     string  `json:"second_number"`
	ServiceGoodsID   int     `json:"service_goods_id"`
	Quantity         float64 `json:"quantity"`
	Fee              float64 `json:"fee"`
	DetailDesc       string  `json:"detail_desc"`
	Element          float64 `json:"element"`
}

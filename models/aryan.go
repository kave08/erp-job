package models

import "time"

type Aryan struct {
	SaleOrder struct {
		CustomerId       int           `json:"customer_id"`
		VoucherDate      time.Duration `json:"voucher_date"`
		SecondNumber     int           `json:"second_number"`
		VoucherDesc      string        `json:"voucher_desc"`
		StockID          int           `json:"stock_id"`
		SaleTypeId       int           `json:"sale_type_id"`
		DeliveryCenterID int           `json:"delivery_center_id"`
		SaleCenterID     int           `json:"sale_center_id"`
		PaymentWayID     int           `json:"payment_way_id"`
		SellerVisitorID  int           `json:"seller_visitor_id"`
		ServiceGoodsID   int           `json:"service_goods_id"`
		Quantity         float64       `json:"quantity"`
		Fee              float64       `json:"fee"`
		DetailDesc       string        `json:"detail_desc"`
	} `json:"sale_order"`
	SaleCustomer struct {
		CustomerID   int    `json:"customer_id"`
		CustomerCode string `json:"customer_code"`
	} `json:"sale_customer"`
	SaleTypeSelect struct {
		BuySaleTypeID   int    `json:"buy_sale_type_id"`
		BuySaleTypeCode string `json:"buy_sale_type_code"`
		BuySaleTypeDesc string `json:"buy_sale_type_desc"`
	} `json:"sale_type_select"`
	SaleCenter4SaleSelect struct {
		StockID   int    `json:"stock_id"`
		StockCode string `json:"stock_code"`
		StockDesc string `json:"stock_desc"`
	} `json:"sale_center_4_sale_select"`
	SalePaymentSelect struct {
		PaymentWayID   int    `json:"payment_way_id"`
		PaymentwayDesc string `json:"payment_way_desc"`
	} `json:"sale_payment_select"`
	SaleCenterSelect struct {
		CentersID   int    `json:"centers_id"`
		CentersCode string `json:"centers_code"`
		CenterDesc  string `json:"center_desc"`
	} `json:"sale_center_select"`
	DeliverCenter_SaleSelect struct {
		CentersID   int    `json:"centers_id"`
		CentersCode string `json:"centers_code"`
	} `json:"deliver_center_sale_select"`
	SalerSelect struct {
		SaleVisitorID   int    `json:"sale_visitor_id"`
		SaleVisitorDesc string `json:"sale_visitor_desc"`
	} `json:"saler_select"`
	SaleSellerVisitor struct {
		CentersID   int    `json:"centers_id"`
		CentersCode string `json:"centers_code"`
	} `json:"sale_seller_visitor"`
	Goods struct {
		ServiceGoodsID   int    `json:"service_goods_id"`
		ServiceGoodsCode string `json:"service_goods_code"`
		ServiceGoodsDesc string `json:"service_goods_desc"`
		GroupId          int    `json:"group_id"`
		TypeID           int    `json:"type_id"`
		SecUnitType      int    `json:"sec_unit_type"`
		Level1           int    `json:"level1"`
		Level2           int    `json:"level2"`
		Level3           int    `json:"level3"`
	} `json:"goods"`
	SaleProforma struct {
		CustomerId       int           `json:"customer_id"`
		VoucherDate      time.Duration `json:"voucher_date"`
		SecondNumber     string        `json:"second_number"`
		VoucherDesc      string        `json:"voucher_desc"`
		StockID          int           `json:"stock_id"`
		SaleTypeId       int           `json:"sale_type_id"`
		DeliveryCenterID int           `json:"delivery_center_id"`
		SaleCenterID     int           `json:"sale_center_id"`
		PaymentWayID     int           `json:"payment_way_id"`
		SellerVisitorID  int           `json:"seller_visitor_id"`
		ServiceGoodsID   int           `json:"service_goods_id"`
		Quantity         float64       `json:"quantity"`
		Fee              float64       `json:"fee"`
		DetailDesc       string        `json:"detail_desc"`
	} `json:"sale_proforma"`
	SaleFactor struct {
		CustomerId       int           `json:"customer_id"`
		VoucherDate      time.Duration `json:"voucher_date"`
		SecondNumber     string        `json:"second_number"`
		VoucherDesc      string        `json:"voucher_desc"`
		StockID          int           `json:"stock_id"`
		SaleTypeId       int           `json:"sale_type_id"`
		DeliveryCenterID int           `json:"delivery_center_id"`
		SaleCenterID     int           `json:"sale_center_id"`
		PaymentWayID     int           `json:"payment_way_id"`
		SellerID         int           `json:"seller_id"`
		SaleManID        int           `json:"sale_man_id"`
		DistributerId    int           `json:"distributer_id"`
		ServiceGoodsID   int           `json:"service_goods_id"`
		Quantity         float64       `json:"quantity"`
		Fee              float64       `json:"fee"`
		DetailDesc       string        `json:"detail_desc"`
		Element          float64       `json:"element"`
	} `json:"sale_factor"`
}

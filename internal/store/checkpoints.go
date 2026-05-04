package store

import (
	"context"
	"time"
)

type Entity string

const (
	EntityInvoice  Entity = "invoice"
	EntityCustomer Entity = "customer"
	EntityProduct  Entity = "product"
	EntityBaseData Entity = "base_data"
)

type Operation string

const (
	OperationInvoiceSaleFactor        Operation = "invoice_sale_factor"
	OperationInvoiceSaleOrder         Operation = "invoice_sale_order"
	OperationInvoiceSalePayment       Operation = "invoice_sale_payment"
	OperationInvoiceSalerSelect       Operation = "invoice_saler_select"
	OperationInvoiceSaleSellerVisitor Operation = "invoice_sale_seller_visitor"
	OperationInvoiceSaleProforma      Operation = "invoice_sale_proforma"
	OperationInvoiceSaleCenter        Operation = "invoice_sale_center"
	OperationInvoiceSaleTypeSelect    Operation = "invoice_sale_type_select"
	OperationCustomerSaleCustomer     Operation = "customer_sale_customer"
	OperationProductsGoods            Operation = "products_goods"
	OperationBaseDataDeliverCenter    Operation = "base_data_deliver_center"
	OperationBaseDataSaleCenterSelect Operation = "base_data_sale_center_select"
)

type DeliveryStatus string

const (
	DeliveryStatusFailed           DeliveryStatus = "failed"
	DeliveryStatusSucceeded        DeliveryStatus = "succeeded"
	DeliveryStatusDelivered        DeliveryStatus = "delivered"
	DeliveryStatusPermanentFailure DeliveryStatus = "permanent_failure"
)

type DeliveryAttempt struct {
	Operation    Operation
	SourceCursor int
	EntityKey    string
	Status       DeliveryStatus
	HTTPStatus   int
	Error        string
	AttemptedAt  time.Time
}

type DeliveredRecord struct {
	SourceCursor int
	EntityKey    string
	HTTPStatus   int
	DeliveredAt  time.Time
}

type EntityAttemptCount struct {
	EntityKey    string
	AttemptCount int
}

type CheckpointStore interface {
	GetSourceProgress(ctx context.Context, entity Entity) (int, error)
	AdvanceSourceProgress(ctx context.Context, entity Entity, lastSourceID int) error
	GetOperationCheckpoint(ctx context.Context, operation Operation) (int, error)
	GetDeliveredEntityKeys(ctx context.Context, operation Operation, entityKeys []string) (map[string]struct{}, error)
	GetAttemptCounts(ctx context.Context, operation Operation, entityKeys []string) (map[string]int, error)
	RecordDeliveryAttempt(ctx context.Context, attempt DeliveryAttempt) error
	MarkBatchDelivered(ctx context.Context, operation Operation, lastSourceID int, records []DeliveredRecord) error
	MarkPermanentFailures(ctx context.Context, operation Operation, entityKeys []string) error
}

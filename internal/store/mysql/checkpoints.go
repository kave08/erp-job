package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"erp-job/internal/store"
)

const (
	insertInvoiceProgressQuery  = "INSERT INTO invoice_progress_info (last_id, page_number, created_at) VALUES(?, ?, ?);"
	getInvoiceProgressQuery     = "SELECT last_id, page_number FROM invoice_progress_info ORDER BY id DESC LIMIT 1;"
	insertCustomerProgressQuery = "INSERT INTO customer_progress_info (last_id, page_number, created_at) VALUES(?, ?, ?);"
	getCustomerProgressQuery    = "SELECT last_id, page_number FROM customer_progress_info ORDER BY id DESC LIMIT 1;"
	insertProductProgressQuery  = "INSERT INTO product_progress_info (last_id, page_number, created_at) VALUES(?, ?, ?);"
	getProductProgressQuery     = "SELECT last_id, page_number FROM product_progress_info ORDER BY id DESC LIMIT 1;"
	insertBaseDataProgressQuery = "INSERT INTO base_data_progress_info (last_id, page_number, created_at) VALUES(?, ?, ?);"
	getBaseDataProgressQuery    = "SELECT last_id, page_number FROM base_data_progress_info ORDER BY id DESC LIMIT 1;"

	insertProductsToGoodsQuery         = "INSERT INTO products_to_goods (p_id, created_at) VALUES(?, ?);"
	getProductsToGoodsQuery            = "SELECT Max(p_id) FROM products_to_goods;"
	insertCustomerToSaleCustomerQuery  = "INSERT INTO customer_to_sale_customer (c_id, created_at) VALUES(?, ?);"
	getCustomerToSaleCustomerQuery     = "SELECT Max(c_id) FROM customer_to_sale_customer;"
	insertInvoiceToSaleFactorQuery     = "INSERT INTO invoice_to_sale_factor (i_id, created_at) VALUES(?, ?);"
	getInvoiceToSaleFactorQuery        = "SELECT Max(i_id) FROM invoice_to_sale_factor;"
	insertInvoiceToSaleOrderQuery      = "INSERT INTO invoice_to_sale_order (i_id, created_at) VALUES(?, ?);"
	getInvoiceToSaleOrderQuery         = "SELECT Max(i_id) FROM invoice_to_sale_order;"
	insertInvoiceToSalePaymentQuery    = "INSERT INTO invoice_to_sale_payment (i_id, created_at) VALUES(?, ?);"
	getInvoiceToSalePaymentQuery       = "SELECT Max(i_id) FROM invoice_to_sale_payment;"
	insertInvoiceToSalerSelectQuery    = "INSERT INTO invoice_to_saler_select (i_id, created_at) VALUES(?, ?);"
	getInvoiceToSalerSelectQuery       = "SELECT Max(i_id) FROM invoice_to_saler_select;"
	insertInvoiceToSaleProformaQuery   = "INSERT INTO invoice_to_sale_proforma (i_id, created_at) VALUES(?, ?);"
	getInvoiceToSaleProformaQuery      = "SELECT Max(i_id) FROM invoice_to_sale_proforma;"
	insertInvoiceToSaleCenterQuery     = "INSERT INTO invoice_to_sale_center (i_id, created_at) VALUES(?, ?);"
	getInvoiceToSaleCenterQuery        = "SELECT Max(i_id) FROM invoice_to_sale_center;"
	insertInvoiceToSaleTypeSelectQuery = "INSERT INTO invoice_to_sale_type_select (i_id, created_at) VALUES(?, ?);"
	getInvoiceToSaleTypeSelectQuery    = "SELECT Max(i_id) FROM invoice_to_sale_type_select;"
	insertBaseDataToDeliverCenterQuery = "INSERT INTO base_data_to_deliver_center (i_id, created_at) VALUES(?, ?);"
	getBaseDataToDeliverCenterQuery    = "SELECT Max(i_id) FROM base_data_to_deliver_center;"
)

type Checkpoints struct {
	db *sql.DB
}

func New(db *sql.DB) *Checkpoints {
	return &Checkpoints{db: db}
}

func scanProgress(row *sql.Row, context string) (store.Progress, error) {
	var progress store.Progress

	err := row.Scan(&progress.LastID, &progress.PageNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return progress, nil
		}

		switch err {
		case sql.ErrConnDone, sql.ErrTxDone:
			log.Printf("database connection or transaction error: %v", err)
		}

		return progress, fmt.Errorf("%s: %w", context, err)
	}

	return progress, nil
}

func scanNullableMaxInt(row *sql.Row, context string) (int, error) {
	var maxID sql.NullInt64

	if err := row.Scan(&maxID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, nil
		case sql.ErrConnDone, sql.ErrTxDone:
			log.Printf("database connection or transaction error: %v", err)
		}

		return 0, fmt.Errorf("%s: %w", context, err)
	}

	if !maxID.Valid {
		return 0, nil
	}

	return int(maxID.Int64), nil
}

func insertTimestampedID(db *sql.DB, query string, id int) error {
	_, err := db.Exec(query, id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (c *Checkpoints) GetInvoiceProgress() (store.Progress, error) {
	return scanProgress(c.db.QueryRow(getInvoiceProgressQuery), "get invoice progress")
}

func (c *Checkpoints) SaveInvoiceProgress(progress store.Progress) error {
	_, err := c.db.Exec(insertInvoiceProgressQuery, progress.LastID, progress.PageNumber, time.Now())
	return err
}

func (c *Checkpoints) GetCustomerProgress() (store.Progress, error) {
	return scanProgress(c.db.QueryRow(getCustomerProgressQuery), "get customer progress")
}

func (c *Checkpoints) SaveCustomerProgress(progress store.Progress) error {
	_, err := c.db.Exec(insertCustomerProgressQuery, progress.LastID, progress.PageNumber, time.Now())
	return err
}

func (c *Checkpoints) GetProductProgress() (store.Progress, error) {
	return scanProgress(c.db.QueryRow(getProductProgressQuery), "get product progress")
}

func (c *Checkpoints) SaveProductProgress(progress store.Progress) error {
	_, err := c.db.Exec(insertProductProgressQuery, progress.LastID, progress.PageNumber, time.Now())
	return err
}

func (c *Checkpoints) GetBaseDataProgress() (store.Progress, error) {
	return scanProgress(c.db.QueryRow(getBaseDataProgressQuery), "get base-data progress")
}

func (c *Checkpoints) SaveBaseDataProgress(progress store.Progress) error {
	_, err := c.db.Exec(insertBaseDataProgressQuery, progress.LastID, progress.PageNumber, time.Now())
	return err
}

func (c *Checkpoints) GetProductsToGoods() (int, error) {
	return scanNullableMaxInt(c.db.QueryRow(getProductsToGoodsQuery), "get products-to-goods checkpoint")
}

func (c *Checkpoints) SaveProductsToGoods(id int) error {
	return insertTimestampedID(c.db, insertProductsToGoodsQuery, id)
}

func (c *Checkpoints) GetCustomerToSaleCustomer() (int, error) {
	return scanNullableMaxInt(c.db.QueryRow(getCustomerToSaleCustomerQuery), "get customer-to-sale-customer checkpoint")
}

func (c *Checkpoints) SaveCustomerToSaleCustomer(id int) error {
	return insertTimestampedID(c.db, insertCustomerToSaleCustomerQuery, id)
}

func (c *Checkpoints) GetInvoiceToSaleFactor() (int, error) {
	return scanNullableMaxInt(c.db.QueryRow(getInvoiceToSaleFactorQuery), "get invoice-to-sale-factor checkpoint")
}

func (c *Checkpoints) SaveInvoiceToSaleFactor(id int) error {
	return insertTimestampedID(c.db, insertInvoiceToSaleFactorQuery, id)
}

func (c *Checkpoints) GetInvoiceToSaleOrder() (int, error) {
	return scanNullableMaxInt(c.db.QueryRow(getInvoiceToSaleOrderQuery), "get invoice-to-sale-order checkpoint")
}

func (c *Checkpoints) SaveInvoiceToSaleOrder(id int) error {
	return insertTimestampedID(c.db, insertInvoiceToSaleOrderQuery, id)
}

func (c *Checkpoints) GetInvoiceToSalePayment() (int, error) {
	return scanNullableMaxInt(c.db.QueryRow(getInvoiceToSalePaymentQuery), "get invoice-to-sale-payment checkpoint")
}

func (c *Checkpoints) SaveInvoiceToSalePayment(id int) error {
	return insertTimestampedID(c.db, insertInvoiceToSalePaymentQuery, id)
}

func (c *Checkpoints) GetInvoiceToSalerSelect() (int, error) {
	return scanNullableMaxInt(c.db.QueryRow(getInvoiceToSalerSelectQuery), "get invoice-to-saler-select checkpoint")
}

func (c *Checkpoints) SaveInvoiceToSalerSelect(id int) error {
	return insertTimestampedID(c.db, insertInvoiceToSalerSelectQuery, id)
}

func (c *Checkpoints) GetInvoiceToSaleProforma() (int, error) {
	return scanNullableMaxInt(c.db.QueryRow(getInvoiceToSaleProformaQuery), "get invoice-to-sale-proforma checkpoint")
}

func (c *Checkpoints) SaveInvoiceToSaleProforma(id int) error {
	return insertTimestampedID(c.db, insertInvoiceToSaleProformaQuery, id)
}

func (c *Checkpoints) GetInvoiceToSaleCenter() (int, error) {
	return scanNullableMaxInt(c.db.QueryRow(getInvoiceToSaleCenterQuery), "get invoice-to-sale-center checkpoint")
}

func (c *Checkpoints) SaveInvoiceToSaleCenter(id int) error {
	return insertTimestampedID(c.db, insertInvoiceToSaleCenterQuery, id)
}

func (c *Checkpoints) GetInvoiceToSaleTypeSelect() (int, error) {
	return scanNullableMaxInt(c.db.QueryRow(getInvoiceToSaleTypeSelectQuery), "get invoice-to-sale-type-select checkpoint")
}

func (c *Checkpoints) SaveInvoiceToSaleTypeSelect(id int) error {
	return insertTimestampedID(c.db, insertInvoiceToSaleTypeSelectQuery, id)
}

func (c *Checkpoints) GetBaseDataToDeliverCenter() (int, error) {
	return scanNullableMaxInt(c.db.QueryRow(getBaseDataToDeliverCenterQuery), "get base-data-to-deliver-center checkpoint")
}

func (c *Checkpoints) SaveBaseDataToDeliverCenter(id int) error {
	return insertTimestampedID(c.db, insertBaseDataToDeliverCenterQuery, id)
}

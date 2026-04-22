package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

const (

	//progress details will store and retrieve pagination's data
	InsertInvoiceProgressQuery  = "INSERT INTO invoice_progress_info (last_id, page_number, created_at) VALUES(?, ?, ?);"
	GetInvoiceProgressQuery     = "SELECT last_id, page_number FROM invoice_progress_info ORDER BY id DESC LIMIT 1;"
	InsertCustomerProgressQuery = "INSERT INTO customer_progress_info (last_id, page_number, created_at) VALUES(?, ?, ?);"
	GetCustomerProgressQuery    = "SELECT last_id, page_number FROM customer_progress_info ORDER BY id DESC LIMIT 1;"
	InsertProductProgressQuery  = "INSERT INTO product_progress_info (last_id, page_number, created_at) VALUES(?, ?, ?);"
	GetProductProgressQuery     = "SELECT last_id, page_number FROM product_progress_info ORDER BY id DESC LIMIT 1;"
	InsertBaseDataProgressQuery = "INSERT INTO base_data_progress_info (last_id, page_number, created_at) VALUES(?, ?, ?);"
	GetBaseDataProgressQuery    = "SELECT last_id, page_number FROM base_data_progress_info ORDER BY id DESC LIMIT 1;"

	//id query will store and retrieve max id's
	InsertProductsToGoodsMaxIdQuery         = "INSERT INTO products_to_goods (p_id, created_at) VALUES(?, ?);"
	GetProductsToGoodsMaxIdQuery            = "SELECT Max(p_id) FROM products_to_goods;"
	InsertCustomerToSaleCustomerMaxIdQuery  = "INSERT INTO customer_to_sale_customer (c_id, created_at) VALUES(?, ?);"
	GetCustomerToSaleCustomerMaxIdQuery     = "SELECT Max(c_id) FROM customer_to_sale_customer;"
	InsertInvoiceToSaleFactorMaxIdQuery     = "INSERT INTO invoice_to_sale_factor (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSaleFactorMaxIdQuery        = "SELECT Max(i_id) FROM invoice_to_sale_factor;"
	InsertInvoiceToSaleOrderMaxIdQuery      = "INSERT INTO invoice_to_sale_order (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSaleOrderMaxIdQuery         = "SELECT Max(i_id) FROM invoice_to_sale_order;"
	InsertInvoiceToSalePaymentMaxIdQuery    = "INSERT INTO invoice_to_sale_payment (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSalePaymentMaxIdQuery       = "SELECT Max(i_id) FROM invoice_to_sale_payment;"
	InsertInvoiceToSalerSelectMaxIdQuery    = "INSERT INTO invoice_to_saler_select (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSalerSelectMaxIdQuery       = "SELECT Max(i_id) FROM invoice_to_saler_select;"
	InsertInvoiceToSaleProformaMaxIdQuery   = "INSERT INTO invoice_to_sale_proforma (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSaleProformaMaxIdQuery      = "SELECT Max(i_id) FROM invoice_to_sale_proforma;"
	InsertInvoiceToSaleCenterMaxIdQuery     = "INSERT INTO invoice_to_sale_center (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSaleCenterMaxIdQuery        = "SELECT Max(i_id) FROM invoice_to_sale_center;"
	InsertInvoiceToSaleTypeSelectMaxIdQuery = "INSERT INTO invoice_to_sale_type_select (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSaleTypeSelectMaxIdQuery    = "SELECT Max(i_id) FROM invoice_to_sale_type_select;"
	InsertBaseDataToDeliverCenterMaxIdQuery = "INSERT INTO base_data_to_deliver_center (i_id, created_at) VALUES(?, ?);"
	GetBaseDataToDeliverCentertMaxIdQuery   = "SELECT Max(i_id) FROM base_data_to_deliver_center;"
	InsertTreasuriesMaxIdQuery              = "INSERT INTO treasuries (t_id, created_at) VALUES(?, ?);"
	GetTreasuriesMaxIdQuery                 = "SELECT Max(t_id) FROM treasuries;"
	InsertInvoiceReturnMaxIdQuery           = "INSERT INTO invoice_return (r_id, created_at) VALUES(?, ?);"
	GetInvoiceReturnMaxIdQuery              = "SELECT Max(r_id) FROM invoice_return;"
)

type DatabaseInterface interface {
	GetInvoiceProgress() (int, int, error)
	InsertInvoiceProgress(laseId, pageNumber int) error
	GetCustomerProgress() (int, int, error)
	InsertCustomerProgress(laseId, pageNumber int) error
	GetProductProgress() (int, int, error)
	InsertProductProgress(laseId, pageNumber int) error
	GetBaseDataProgress() (int, int, error)
	InsertBaseDataProgress(laseId, pageNumber int) error

	InsertProductsToGoods(p_id int) error
	GetProductsToGoods() (int, error)

	InsertCustomerToSaleCustomer(c_id int) error
	GetCustomerToSaleCustomer() (int, error)

	InsertInvoiceToSaleFactor(i_id int) error
	GetInvoiceToSaleFactor() (int, error)

	InsertInvoiceToSaleOrder(i_id int) error
	GetInvoiceToSaleOrder() (int, error)

	InsertInvoiceToSalePayment(i_id int) error
	GetInvoiceToSalePayment() (int, error)

	InsertInvoiceToSalerSelect(i_id int) error
	GetInvoiceToSalerSelect() (int, error)

	InsertInvoiceToSaleProforma(i_id int) error
	GetInvoiceToSaleProforma() (int, error)

	InsertInvoiceToSaleCenter(i_id int) error
	GetInvoiceToSaleCenter() (int, error)

	InsertInvoiceToSaleTypeSelect(i_id int) error
	GetInvoiceToSaleTypeSelect() (int, error)

	InsertBaseDataToDeliverCenter(i_id int) error
	GetBaseDataToDeliverCenter() (int, error)

	InsertTreasuries(t_id int) error
	GetTreasuries() (int, error)

	InsertInvoiceReturn(r_id int) error
	GetInvoiceReturn() (int, error)
}

type Database struct {
	sdb *sql.DB
}

func NewDatabase(sdb *sql.DB) DatabaseInterface {
	return &Database{
		sdb: sdb,
	}
}

func scanNullableMaxInt(row *sql.Row, context string) (int, error) {
	var maxID sql.NullInt64

	err := row.Scan(&maxID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, nil
		case sql.ErrConnDone, sql.ErrTxDone:
			log.Printf("Database connection or transaction error: %v", err)
			return 0, fmt.Errorf("%s %w", context, err)
		default:
			return 0, fmt.Errorf("%s %w", context, err)
		}
	}

	if !maxID.Valid {
		return 0, nil
	}

	return int(maxID.Int64), nil
}

// GetInvoiceProgress implements DatabaseInterface
func (d *Database) GetInvoiceProgress() (int, int, error) {
	var laseId int
	var pageNumber int

	err := d.sdb.QueryRow(GetInvoiceProgressQuery).Scan(&laseId, &pageNumber)
	if err != nil {
		switch err {
		case sql.ErrConnDone:
			log.Printf("Database connection error: %v", err)
			return laseId, pageNumber, fmt.Errorf("@ERP.repository.database.database.GetInvoiceProgress %w %v ", err, laseId)
		case sql.ErrTxDone:
			log.Printf("Database transaction error: %v", err)
			return laseId, pageNumber, fmt.Errorf("@ERP.repository.database.database.GetInvoiceProgress %w %v ", err, laseId)
		default:
			return laseId, pageNumber, fmt.Errorf("@ERP.repository.database.database.GetInvoiceProgress %w %v ", err, laseId)
		}
	}

	return laseId, pageNumber, err
}

// InsertInvoiceProgress implements DatabaseInterface
func (d *Database) InsertInvoiceProgress(laseId, pageNumber int) error {
	_, err := d.sdb.Exec(InsertInvoiceProgressQuery, laseId, pageNumber, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetCustomerProgress implements DatabaseInterface
func (d *Database) GetCustomerProgress() (int, int, error) {
	var laseId int
	var pageNumber int

	err := d.sdb.QueryRow(GetCustomerProgressQuery).Scan(&laseId, &pageNumber)
	if err != nil {
		switch err {
		case sql.ErrConnDone:
			log.Printf("Database connection error: %v", err)
			return laseId, pageNumber, fmt.Errorf("@ERP.repository.database.database.GetCustomerProgress %w %v ", err, laseId)
		case sql.ErrTxDone:
			log.Printf("Database transaction error: %v", err)
			return laseId, pageNumber, fmt.Errorf("@ERP.repository.database.database.GetCustomerProgress %w %v ", err, laseId)
		default:
			return laseId, pageNumber, fmt.Errorf("@ERP.repository.database.database.GetCustomerProgress %w %v ", err, laseId)
		}
	}

	return laseId, pageNumber, err
}

// InsertCustomerProgress implements DatabaseInterface
func (d *Database) InsertCustomerProgress(laseId, pageNumber int) error {
	_, err := d.sdb.Exec(InsertCustomerProgressQuery, laseId, pageNumber, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetProductProgress implements DatabaseInterface
func (d *Database) GetProductProgress() (int, int, error) {
	var laseId int
	var pageNumber int

	err := d.sdb.QueryRow(GetProductProgressQuery).Scan(&laseId, &pageNumber)
	if err != nil {
		switch err {
		case sql.ErrConnDone:
			log.Printf("Database connection error: %v", err)
			return laseId, pageNumber, fmt.Errorf("@ERP.repository.database.database.GetProductProgress %w %v ", err, laseId)
		case sql.ErrTxDone:
			log.Printf("Database transaction error: %v", err)
			return laseId, pageNumber, fmt.Errorf("@ERP.repository.database.database.GetProductProgress %w %v ", err, laseId)
		default:
			return laseId, pageNumber, fmt.Errorf("@ERP.repository.database.database.GetProductProgress %w %v ", err, laseId)
		}
	}

	return laseId, pageNumber, err
}

// InsertProductProgress implements DatabaseInterface
func (d *Database) InsertProductProgress(laseId, pageNumber int) error {
	_, err := d.sdb.Exec(InsertProductProgressQuery, laseId, pageNumber, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetBaseDataProgress implements DatabaseInterface
func (d *Database) GetBaseDataProgress() (int, int, error) {
	var laseId int
	var pageNumber int

	err := d.sdb.QueryRow(GetBaseDataProgressQuery).Scan(&laseId, &pageNumber)
	if err != nil {
		switch err {
		case sql.ErrConnDone:
			log.Printf("Database connection error: %v", err)
			return laseId, pageNumber, fmt.Errorf("@ERP.repository.database.database.GetBaseDataProgress %w %v ", err, laseId)
		case sql.ErrTxDone:
			log.Printf("Database transaction error: %v", err)
			return laseId, pageNumber, fmt.Errorf("@ERP.repository.database.database.GetBaseDataProgress %w %v ", err, laseId)
		default:
			return laseId, pageNumber, fmt.Errorf("@ERP.repository.database.database.GetBaseDataProgress %w %v ", err, laseId)
		}
	}

	return laseId, pageNumber, err
}

// InsertBaseDataProgress implements DatabaseInterface
func (d *Database) InsertBaseDataProgress(laseId, pageNumber int) error {
	_, err := d.sdb.Exec(InsertBaseDataProgressQuery, laseId, pageNumber, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetCustomerToSaleCustomer implements DatabaseInterface
func (d *Database) GetCustomerToSaleCustomer() (int, error) {
	return scanNullableMaxInt(d.sdb.QueryRow(GetCustomerToSaleCustomerMaxIdQuery), "@ERP.repository.databese.databese.GetCustomerToSaleCustomer")
}

// InsertCustomerToSaleCustomer implements DatabaseInterface
func (d *Database) InsertCustomerToSaleCustomer(c_id int) error {
	_, err := d.sdb.Exec(InsertCustomerToSaleCustomerMaxIdQuery, c_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// InsertInvoiceToSaleFactor implements DatabaseInterface
func (d *Database) InsertInvoiceToSaleFactor(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToSaleFactorMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetInvoiceToSaleFactor implements DatabaseInterface
func (d *Database) GetInvoiceToSaleFactor() (int, error) {
	return scanNullableMaxInt(d.sdb.QueryRow(GetInvoiceToSaleFactorMaxIdQuery), "@ERP.repository.databese.databese.GetInvoiceToSaleFactor")
}

// InsertInvoiceToSaleOrder implements DatabaseInterface
func (d *Database) InsertInvoiceToSaleOrder(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToSaleOrderMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetInvoiceToSaleOrder implements DatabaseInterface
func (d *Database) GetInvoiceToSaleOrder() (int, error) {
	return scanNullableMaxInt(d.sdb.QueryRow(GetInvoiceToSaleOrderMaxIdQuery), "@ERP.repository.databese.databese.GetInvoiceToSaleOrder")
}

// InsertInvoiceToSalePayment implements DatabaseInterface
func (d *Database) InsertInvoiceToSalePayment(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToSalePaymentMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetInvoiceToSalePayment implements DatabaseInterface
func (d *Database) GetInvoiceToSalePayment() (int, error) {
	return scanNullableMaxInt(d.sdb.QueryRow(GetInvoiceToSalePaymentMaxIdQuery), "@ERP.repository.databese.databese.GetInvoiceToSalePayment")
}

// InsertInvoiceToSalerSelect implements DatabaseInterface
func (d *Database) InsertInvoiceToSalerSelect(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToSalerSelectMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetInvoiceToSalerSelect implements DatabaseInterface
func (d *Database) GetInvoiceToSalerSelect() (int, error) {
	return scanNullableMaxInt(d.sdb.QueryRow(GetInvoiceToSalerSelectMaxIdQuery), "@ERP.repository.databese.databese.GetInvoiceToSalerSelect")
}

// InsertInvoiceToSaleProforma implements DatabaseInterface
func (d *Database) InsertInvoiceToSaleProforma(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToSaleProformaMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetInvoiceToSaleProforma implements DatabaseInterface
func (d *Database) GetInvoiceToSaleProforma() (int, error) {
	return scanNullableMaxInt(d.sdb.QueryRow(GetInvoiceToSaleProformaMaxIdQuery), "@ERP.repository.databese.databese.GetInvoiceToSaleProforma")
}

// InsertInvoiceToSaleTypeSelect implements DatabaseInterface
func (d *Database) InsertInvoiceToSaleTypeSelect(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToSaleTypeSelectMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetInvoiceToSaleTypeSelect implements DatabaseInterface
func (d *Database) GetInvoiceToSaleTypeSelect() (int, error) {
	return scanNullableMaxInt(d.sdb.QueryRow(GetInvoiceToSaleTypeSelectMaxIdQuery), "@ERP.repository.databese.databese.GetInvoiceToSaleTypeSelect")
}

// InsertInvoiceToSaleCenter implements DatabaseInterface
func (d *Database) InsertInvoiceToSaleCenter(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToSaleCenterMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetInvoiceToSaleCenter implements DatabaseInterface
func (d *Database) GetInvoiceToSaleCenter() (int, error) {
	return scanNullableMaxInt(d.sdb.QueryRow(GetInvoiceToSaleCenterMaxIdQuery), "@ERP.repository.databese.databese.GetInvoiceToSaleCenter")
}

// InsertBaseDataToDeliverCenter implements DatabaseInterface
func (d *Database) InsertBaseDataToDeliverCenter(i_id int) error {
	_, err := d.sdb.Exec(InsertBaseDataToDeliverCenterMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetBaseDataToDeliverCenter implements DatabaseInterface
func (d *Database) GetBaseDataToDeliverCenter() (int, error) {
	return scanNullableMaxInt(d.sdb.QueryRow(GetBaseDataToDeliverCentertMaxIdQuery), "@ERP.repository.databese.databese.GetBaseDataToDeliverCenter")
}

// GetProductsToGoods implements DatabaseInterface
func (d *Database) GetProductsToGoods() (int, error) {
	return scanNullableMaxInt(d.sdb.QueryRow(GetProductsToGoodsMaxIdQuery), "@ERP.repository.databese.databese.GetProductsToGoods")
}

// GetInvoiceReturn implements DatabaseInterface
func (d *Database) GetInvoiceReturn() (int, error) {
	return scanNullableMaxInt(d.sdb.QueryRow(GetInvoiceReturnMaxIdQuery), "@ERP.repository.databese.databese.GetInvoiceReturn")
}

// GetTreasuries implements DatabaseInterface
func (d *Database) GetTreasuries() (int, error) {
	return scanNullableMaxInt(d.sdb.QueryRow(GetTreasuriesMaxIdQuery), "@ERP.repository.databese.databese.GetTreasuries")
}

// InsertProductsToGoods implements DatabaseInterface
func (d *Database) InsertProductsToGoods(p_id int) error {
	_, err := d.sdb.Exec(InsertProductsToGoodsMaxIdQuery, p_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// InsertInvoiceReturn implements DatabaseInterface
func (d *Database) InsertInvoiceReturn(r_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceReturnMaxIdQuery, r_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// InsertTreasuries implements DatabaseInterface
func (d *Database) InsertTreasuries(t_id int) error {
	_, err := d.sdb.Exec(InsertTreasuriesMaxIdQuery, t_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

const (
	//TODO: basedata and details

	InsertProductsToGoodsMaxIdQuery = "INSERT INTO products_to_goods (p_id, created_at) VALUES(?, ?);"
	GetProductsToGoodsMaxIdQuery    = "SELECT Max(p_id) FROM products_to_goods;"

	InsertCustomerToSaleCustomerMaxIdQuery = "INSERT INTO customer_to_sale_customer (c_id, created_at) VALUES(?, ?);"
	GetCustomerToSaleCustomerMaxIdQuery    = "SELECT Max(c_id) FROM customer_to_sale_customer;"

	InsertInvoiceToSaleFactorMaxIdQuery = "INSERT INTO invoice_to_sale_factor (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSaleFactorMaxIdQuery    = "SELECT Max(i_id) FROM invoice_to_sale_factor;"

	InsertInvoiceToGoodsMaxIdQuery = "INSERT INTO invoice_to_goods (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToGoodsMaxIdQuery    = "SELECT Max(i_id) FROM invoice_to_goods;"

	InsertInvoiceToSaleOrderMaxIdQuery = "INSERT INTO invoice_to_sale_order (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSaleOrderMaxIdQuery    = "SELECT Max(i_id) FROM invoice_to_sale_order;"

	InsertInvoiceToSalePaymentMaxIdQuery = "INSERT INTO invoice_to_sale_payment (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSalePaymentMaxIdQuery    = "SELECT Max(i_id) FROM invoice_to_sale_payment;"

	InsertInvoiceToSalerSelectMaxIdQuery = "INSERT INTO invoice_to_saler_select (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSalerSelectMaxIdQuery    = "SELECT Max(i_id) FROM invoice_to_saler_select;"

	InsertInvoiceToSaleProformaMaxIdQuery = "INSERT INTO invoice_to_sale_proforma (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSaleProformaMaxIdQuery    = "SELECT Max(i_id) FROM invoice_to_sale_proforma;"

	InsertTreasuriesMaxIdQuery = "INSERT INTO treasuries (t_id, created_at) VALUES(?, ?);"
	GetTreasuriesMaxIdQuery    = "SELECT Max(t_id) FROM treasuries;"

	InsertRevertedMaxIdQuery = "INSERT INTO reverted (r_id, created_at) VALUES(?, ?);"
	GetRevertedMaxIdQuery    = "SELECT Max(r_id) FROM reverted;"
)

type DatabaseInterface interface {
	InsertProductsToGoods(p_id int) error
	GetProductsToGoods() (int, error)

	InsertCustomerToSaleCustomer(c_id int) error
	GetCustomerToSaleCustomer() (int, error)

	InsertInvoiceToSaleFactor(i_id int) error
	GetInvoiceToSaleFactor() (int, error)

	InsertInvoiceToGoods(i_id int) error
	GetInvoiceToGoods() (int, error)

	InsertInvoiceToSaleOrder(i_id int) error
	GetInvoiceToSaleOrder() (int, error)

	InsertInvoiceToSalePayment(i_id int) error
	GetInvoiceToSalePayment() (int, error)

	InsertInvoiceToSalerSelect(i_id int) error
	GetInvoiceToSalerSelect() (int, error)

	InsertInvoiceToSaleProforma(i_id int) error
	GetInvoiceToSaleProforma() (int, error)

	InsertTreasuries(t_id int) error
	GetTreasuries() (int, error)

	InsertReverted(r_id int) error
	GetReverted() (int, error)
}

type Database struct {
	sdb *sql.DB
}

func NewDatabase(sdb *sql.DB) DatabaseInterface {
	return &Database{
		sdb: sdb,
	}
}

// GetCustomerToSaleCustomer implements DatabaseInterface
func (d *Database) GetCustomerToSaleCustomer() (int, error) {
	var id int
	var maxId sql.NullInt64
	err := d.sdb.QueryRow(GetCustomerToSaleCustomerMaxIdQuery).Scan(&maxId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			if maxId.Valid {
				id = int(maxId.Int64)
			} else {
				id = 0
			}
		case sql.ErrConnDone, sql.ErrTxDone:
			log.Printf("Database connection or transaction error: %v", err)
			return id, fmt.Errorf("@ERP.repository.databese.databese.GetCustomerToSaleCustomer %w %v ", err, maxId)
		default:
			return id, fmt.Errorf("@ERP.repository.databese.databese.GetCustomerToSaleCustomer %w %v ", err, maxId)
		}
	}

	return id, err
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
	var id int
	var maxId sql.NullInt64
	err := d.sdb.QueryRow(GetInvoiceToSaleFactorMaxIdQuery).Scan(&maxId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			if maxId.Valid {
				id = int(maxId.Int64)
			} else {
				id = 0
			}
		case sql.ErrConnDone, sql.ErrTxDone:
			log.Printf("Database connection or transaction error: %v", err)
			return id, fmt.Errorf("@ERP.repository.databese.databese.GetInvoiceToSaleFactor %w %v ", err, maxId)
		default:
			return id, fmt.Errorf("@ERP.repository.databese.databese.GetInvoiceToSaleFactor %w %v ", err, maxId)
		}
	}

	return id, err
}

// InsertInvoiceToGoods implements DatabaseInterface
func (d *Database) InsertInvoiceToGoods(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToGoodsMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetInvoiceToGoods implements DatabaseInterface
func (d *Database) GetInvoiceToGoods() (int, error) {
	var id int
	var maxId sql.NullInt64
	err := d.sdb.QueryRow(GetInvoiceToGoodsMaxIdQuery).Scan(&maxId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			if maxId.Valid {
				id = int(maxId.Int64)
			} else {
				id = 0
			}
		case sql.ErrConnDone, sql.ErrTxDone:
			log.Printf("Database connection or transaction error: %v", err)
			return id, fmt.Errorf("@ERP.repository.databese.databese.GetInvoiceToGoods %w %v ", err, maxId)
		default:
			return id, fmt.Errorf("@ERP.repository.databese.databese.GetInvoiceToGoods %w %v ", err, maxId)
		}
	}

	return id, err
}

func (d *Database) InsertInvoiceToSaleOrder(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToSaleOrderMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetInvoiceToSaleOrder() (int, error) {
	var maxId int
	err := d.sdb.QueryRow(GetInvoiceToSaleOrderMaxIdQuery).Scan(&maxId)
	if err != nil {
		return 0, err
	}

	return maxId, err
}

func (d *Database) InsertInvoiceToSalePayment(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToSalePaymentMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetInvoiceToSalePayment() (int, error) {
	var maxId int
	err := d.sdb.QueryRow(GetInvoiceToSalePaymentMaxIdQuery).Scan(&maxId)
	if err != nil {
		return 0, err
	}

	return maxId, err
}

func (d *Database) InsertInvoiceToSalerSelect(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToSalerSelectMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetInvoiceToSalerSelect() (int, error) {
	var maxId int
	err := d.sdb.QueryRow(GetInvoiceToSalerSelectMaxIdQuery).Scan(&maxId)
	if err != nil {
		return 0, err
	}

	return maxId, err
}

func (d *Database) InsertInvoiceToSaleProforma(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToSaleProformaMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetInvoiceToSaleProforma() (int, error) {
	var maxId int
	err := d.sdb.QueryRow(GetInvoiceToSaleProformaMaxIdQuery).Scan(&maxId)
	if err != nil {
		return 0, err
	}

	return maxId, err
}

// GetProduct implements DatabaseInterface
func (d *Database) GetProductsToGoods() (int, error) {
	var maxId int
	err := d.sdb.QueryRow(GetProductsToGoodsMaxIdQuery).Scan(&maxId)
	if err != nil {
		return 0, err
	}

	return maxId, err
}

// GetReverted implements DatabaseInterface
func (d *Database) GetReverted() (int, error) {
	var maxId int
	err := d.sdb.QueryRow(GetRevertedMaxIdQuery).Scan(&maxId)
	if err != nil {
		return 0, err
	}

	return maxId, err
}

// GetTreasuries implements DatabaseInterface
func (d *Database) GetTreasuries() (int, error) {
	var maxId int
	err := d.sdb.QueryRow(GetTreasuriesMaxIdQuery).Scan(&maxId)
	if err != nil {
		return 0, err
	}

	return maxId, err
}

// InsertProduct implements DatabaseInterface
func (d *Database) InsertProductsToGoods(p_id int) error {
	_, err := d.sdb.Exec(InsertProductsToGoodsMaxIdQuery, p_id)
	if err != nil {
		return err
	}

	return nil
}

// InsertReverted implements DatabaseInterface
func (d *Database) InsertReverted(r_id int) error {
	_, err := d.sdb.Exec(InsertRevertedMaxIdQuery, r_id)
	if err != nil {
		return err
	}

	return nil
}

// InsertTreasuries implements DatabaseInterface
func (d *Database) InsertTreasuries(t_id int) error {
	_, err := d.sdb.Exec(InsertTreasuriesMaxIdQuery, t_id)
	if err != nil {
		return err
	}

	return nil
}

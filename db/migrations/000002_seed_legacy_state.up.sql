SET @now := NOW(6);

SET @seed_invoice_sale_factor = IF(
	(SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'invoice_to_sale_factor') > 0,
	"INSERT INTO operation_checkpoint (operation_name, last_source_id, created_at, updated_at)
	 SELECT 'invoice_sale_factor', COALESCE(MAX(i_id), 0), NOW(6), NOW(6) FROM invoice_to_sale_factor
	 ON DUPLICATE KEY UPDATE
	  last_source_id = GREATEST(last_source_id, VALUES(last_source_id)),
	  updated_at = VALUES(updated_at)",
	"SELECT 1"
);
PREPARE stmt FROM @seed_invoice_sale_factor;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @seed_invoice_sale_order = IF(
	(SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'invoice_to_sale_order') > 0,
	"INSERT INTO operation_checkpoint (operation_name, last_source_id, created_at, updated_at)
	 SELECT 'invoice_sale_order', COALESCE(MAX(i_id), 0), NOW(6), NOW(6) FROM invoice_to_sale_order
	 ON DUPLICATE KEY UPDATE
	  last_source_id = GREATEST(last_source_id, VALUES(last_source_id)),
	  updated_at = VALUES(updated_at)",
	"SELECT 1"
);
PREPARE stmt FROM @seed_invoice_sale_order;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @seed_invoice_sale_payment = IF(
	(SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'invoice_to_sale_payment') > 0,
	"INSERT INTO operation_checkpoint (operation_name, last_source_id, created_at, updated_at)
	 SELECT 'invoice_sale_payment', COALESCE(MAX(i_id), 0), NOW(6), NOW(6) FROM invoice_to_sale_payment
	 ON DUPLICATE KEY UPDATE
	  last_source_id = GREATEST(last_source_id, VALUES(last_source_id)),
	  updated_at = VALUES(updated_at)",
	"SELECT 1"
);
PREPARE stmt FROM @seed_invoice_sale_payment;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @seed_invoice_saler_select = IF(
	(SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'invoice_to_saler_select') > 0,
	"INSERT INTO operation_checkpoint (operation_name, last_source_id, created_at, updated_at)
	 SELECT 'invoice_saler_select', COALESCE(MAX(i_id), 0), NOW(6), NOW(6) FROM invoice_to_saler_select
	 ON DUPLICATE KEY UPDATE
	  last_source_id = GREATEST(last_source_id, VALUES(last_source_id)),
	  updated_at = VALUES(updated_at)",
	"SELECT 1"
);
PREPARE stmt FROM @seed_invoice_saler_select;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @seed_invoice_sale_proforma = IF(
	(SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'invoice_to_sale_proforma') > 0,
	"INSERT INTO operation_checkpoint (operation_name, last_source_id, created_at, updated_at)
	 SELECT 'invoice_sale_proforma', COALESCE(MAX(i_id), 0), NOW(6), NOW(6) FROM invoice_to_sale_proforma
	 ON DUPLICATE KEY UPDATE
	  last_source_id = GREATEST(last_source_id, VALUES(last_source_id)),
	  updated_at = VALUES(updated_at)",
	"SELECT 1"
);
PREPARE stmt FROM @seed_invoice_sale_proforma;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @seed_invoice_sale_center = IF(
	(SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'invoice_to_sale_center') > 0,
	"INSERT INTO operation_checkpoint (operation_name, last_source_id, created_at, updated_at)
	 SELECT 'invoice_sale_center', COALESCE(MAX(i_id), 0), NOW(6), NOW(6) FROM invoice_to_sale_center
	 ON DUPLICATE KEY UPDATE
	  last_source_id = GREATEST(last_source_id, VALUES(last_source_id)),
	  updated_at = VALUES(updated_at)",
	"SELECT 1"
);
PREPARE stmt FROM @seed_invoice_sale_center;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @seed_invoice_sale_type_select = IF(
	(SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'invoice_to_sale_type_select') > 0,
	"INSERT INTO operation_checkpoint (operation_name, last_source_id, created_at, updated_at)
	 SELECT 'invoice_sale_type_select', COALESCE(MAX(i_id), 0), NOW(6), NOW(6) FROM invoice_to_sale_type_select
	 ON DUPLICATE KEY UPDATE
	  last_source_id = GREATEST(last_source_id, VALUES(last_source_id)),
	  updated_at = VALUES(updated_at)",
	"SELECT 1"
);
PREPARE stmt FROM @seed_invoice_sale_type_select;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @seed_customer_sale_customer = IF(
	(SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'customer_to_sale_customer') > 0,
	"INSERT INTO operation_checkpoint (operation_name, last_source_id, created_at, updated_at)
	 SELECT 'customer_sale_customer', COALESCE(MAX(c_id), 0), NOW(6), NOW(6) FROM customer_to_sale_customer
	 ON DUPLICATE KEY UPDATE
	  last_source_id = GREATEST(last_source_id, VALUES(last_source_id)),
	  updated_at = VALUES(updated_at)",
	"SELECT 1"
);
PREPARE stmt FROM @seed_customer_sale_customer;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @seed_products_goods = IF(
	(SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'products_to_goods') > 0,
	"INSERT INTO operation_checkpoint (operation_name, last_source_id, created_at, updated_at)
	 SELECT 'products_goods', COALESCE(MAX(p_id), 0), NOW(6), NOW(6) FROM products_to_goods
	 ON DUPLICATE KEY UPDATE
	  last_source_id = GREATEST(last_source_id, VALUES(last_source_id)),
	  updated_at = VALUES(updated_at)",
	"SELECT 1"
);
PREPARE stmt FROM @seed_products_goods;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @seed_base_data_deliver_center = IF(
	(SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'base_data_to_deliver_center') > 0,
	"INSERT INTO operation_checkpoint (operation_name, last_source_id, created_at, updated_at)
	 SELECT 'base_data_deliver_center', COALESCE(MAX(i_id), 0), NOW(6), NOW(6) FROM base_data_to_deliver_center
	 ON DUPLICATE KEY UPDATE
	  last_source_id = GREATEST(last_source_id, VALUES(last_source_id)),
	  updated_at = VALUES(updated_at)",
	"SELECT 1"
);
PREPARE stmt FROM @seed_base_data_deliver_center;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

INSERT INTO source_progress (entity, last_source_id, created_at, updated_at)
VALUES
	(
		'invoice',
		LEAST(
			COALESCE((SELECT last_source_id FROM operation_checkpoint WHERE operation_name = 'invoice_sale_factor'), 0),
			COALESCE((SELECT last_source_id FROM operation_checkpoint WHERE operation_name = 'invoice_sale_order'), 0),
			COALESCE((SELECT last_source_id FROM operation_checkpoint WHERE operation_name = 'invoice_sale_payment'), 0),
			COALESCE((SELECT last_source_id FROM operation_checkpoint WHERE operation_name = 'invoice_saler_select'), 0),
			COALESCE((SELECT last_source_id FROM operation_checkpoint WHERE operation_name = 'invoice_sale_proforma'), 0),
			COALESCE((SELECT last_source_id FROM operation_checkpoint WHERE operation_name = 'invoice_sale_center'), 0),
			COALESCE((SELECT last_source_id FROM operation_checkpoint WHERE operation_name = 'invoice_sale_type_select'), 0)
		),
		@now,
		@now
	),
	(
		'customer',
		COALESCE((SELECT last_source_id FROM operation_checkpoint WHERE operation_name = 'customer_sale_customer'), 0),
		@now,
		@now
	),
	(
		'product',
		COALESCE((SELECT last_source_id FROM operation_checkpoint WHERE operation_name = 'products_goods'), 0),
		@now,
		@now
	),
	(
		'base_data',
		COALESCE((SELECT last_source_id FROM operation_checkpoint WHERE operation_name = 'base_data_deliver_center'), 0),
		@now,
		@now
	)
ON DUPLICATE KEY UPDATE
	last_source_id = GREATEST(source_progress.last_source_id, VALUES(last_source_id)),
	updated_at = VALUES(updated_at);

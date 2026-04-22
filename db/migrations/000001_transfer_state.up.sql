CREATE TABLE IF NOT EXISTS source_progress (
	entity VARCHAR(64) NOT NULL PRIMARY KEY,
	last_source_id BIGINT NOT NULL DEFAULT 0,
	created_at DATETIME(6) NOT NULL,
	updated_at DATETIME(6) NOT NULL
);

CREATE TABLE IF NOT EXISTS operation_checkpoint (
	operation_name VARCHAR(64) NOT NULL PRIMARY KEY,
	last_source_id BIGINT NOT NULL DEFAULT 0,
	created_at DATETIME(6) NOT NULL,
	updated_at DATETIME(6) NOT NULL
);

CREATE TABLE IF NOT EXISTS delivery_state (
	operation_name VARCHAR(64) NOT NULL,
	source_id BIGINT NOT NULL,
	dedupe_key VARCHAR(255) NOT NULL,
	status VARCHAR(32) NOT NULL,
	attempt_count INT NOT NULL DEFAULT 0,
	last_http_status INT NULL,
	last_error TEXT NULL,
	delivered_at DATETIME(6) NULL,
	created_at DATETIME(6) NOT NULL,
	updated_at DATETIME(6) NOT NULL,
	PRIMARY KEY (operation_name, source_id),
	KEY idx_delivery_state_status (status),
	KEY idx_delivery_state_dedupe_key (dedupe_key)
);

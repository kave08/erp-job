SET @delivery_state_exists := (
	SELECT COUNT(*)
	FROM information_schema.TABLES
	WHERE TABLE_SCHEMA = DATABASE()
	  AND TABLE_NAME = 'delivery_state'
);

SET @legacy_exists := (
	SELECT COUNT(*)
	FROM information_schema.TABLES
	WHERE TABLE_SCHEMA = DATABASE()
	  AND TABLE_NAME = 'delivery_state_legacy_000003'
);

SET @delivery_state_final := IF(
	@delivery_state_exists > 0
	AND (
		SELECT COUNT(*)
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'delivery_state'
		  AND COLUMN_NAME = 'source_cursor'
	) > 0
	AND (
		SELECT COUNT(*)
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'delivery_state'
		  AND COLUMN_NAME = 'entity_key'
	) > 0
	AND (
		SELECT COUNT(*)
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'delivery_state'
		  AND INDEX_NAME = 'PRIMARY'
		  AND COLUMN_NAME = 'entity_key'
		  AND SEQ_IN_INDEX = 2
	) > 0,
	1,
	0
);

SET @source_table := CASE
	WHEN @delivery_state_exists > 0 THEN 'delivery_state'
	WHEN @legacy_exists > 0 THEN 'delivery_state_legacy_000003'
	ELSE ''
END;

SET @source_cursor_col := CASE
	WHEN @source_table = 'delivery_state' AND (
		SELECT COUNT(*)
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'delivery_state'
		  AND COLUMN_NAME = 'source_cursor'
	) > 0 THEN 'source_cursor'
	WHEN @source_table = 'delivery_state_legacy_000003' AND (
		SELECT COUNT(*)
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'delivery_state_legacy_000003'
		  AND COLUMN_NAME = 'source_cursor'
	) > 0 THEN 'source_cursor'
	ELSE 'source_id'
END;

SET @entity_key_col := CASE
	WHEN @source_table = 'delivery_state' AND (
		SELECT COUNT(*)
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'delivery_state'
		  AND COLUMN_NAME = 'entity_key'
	) > 0 THEN 'entity_key'
	WHEN @source_table = 'delivery_state_legacy_000003' AND (
		SELECT COUNT(*)
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'delivery_state_legacy_000003'
		  AND COLUMN_NAME = 'entity_key'
	) > 0 THEN 'entity_key'
	ELSE 'dedupe_key'
END;

SET @needs_rebuild := IF(@delivery_state_final = 1 OR @source_table = '', 0, 1);

SET @drop_rebuild_sql := IF(
	@needs_rebuild = 1,
	'DROP TABLE IF EXISTS delivery_state_rebuild_000003',
	'SELECT 1'
);
PREPARE stmt FROM @drop_rebuild_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @create_rebuild_sql := IF(
	@needs_rebuild = 1,
	"CREATE TABLE delivery_state_rebuild_000003 (
		operation_name VARCHAR(64) NOT NULL,
		source_cursor BIGINT NOT NULL,
		entity_key VARCHAR(255) NOT NULL,
		status VARCHAR(32) NOT NULL,
		attempt_count INT NOT NULL DEFAULT 0,
		last_http_status INT NULL,
		last_error TEXT NULL,
		delivered_at DATETIME(6) NULL,
		created_at DATETIME(6) NOT NULL,
		updated_at DATETIME(6) NOT NULL,
		PRIMARY KEY (operation_name, entity_key),
		KEY idx_delivery_state_status (status),
		KEY idx_delivery_state_entity_key (entity_key),
		KEY idx_delivery_state_operation_cursor (operation_name, source_cursor)
	)",
	'SELECT 1'
);
PREPARE stmt FROM @create_rebuild_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @insert_rebuild_sql := IF(
	@needs_rebuild = 1,
	CONCAT(
		"INSERT INTO delivery_state_rebuild_000003 (
			operation_name,
			source_cursor,
			entity_key,
			status,
			attempt_count,
			last_http_status,
			last_error,
			delivered_at,
			created_at,
			updated_at
		)
		SELECT
			agg.operation_name,
			agg.source_cursor,
			agg.entity_key,
			agg.status,
			agg.attempt_count,
			(
				SELECT src.last_http_status
				FROM ", @source_table, " src
				WHERE src.operation_name = agg.operation_name
				  AND src.", @entity_key_col, " = agg.entity_key
				ORDER BY src.updated_at DESC, src.created_at DESC
				LIMIT 1
			) AS last_http_status,
			(
				SELECT src.last_error
				FROM ", @source_table, " src
				WHERE src.operation_name = agg.operation_name
				  AND src.", @entity_key_col, " = agg.entity_key
				ORDER BY src.updated_at DESC, src.created_at DESC
				LIMIT 1
			) AS last_error,
			agg.delivered_at,
			agg.created_at,
			agg.updated_at
		FROM (
			SELECT
				operation_name,
				MAX(", @source_cursor_col, ") AS source_cursor,
				", @entity_key_col, " AS entity_key,
				CASE
					WHEN SUM(status = 'delivered') > 0 THEN 'delivered'
					WHEN SUM(status = 'succeeded') > 0 THEN 'succeeded'
					ELSE 'failed'
				END AS status,
				SUM(attempt_count) AS attempt_count,
				MAX(delivered_at) AS delivered_at,
				MIN(created_at) AS created_at,
				MAX(updated_at) AS updated_at
			FROM ", @source_table, "
			GROUP BY operation_name, ", @entity_key_col, "
		) agg"
	),
	'SELECT 1'
);
PREPARE stmt FROM @insert_rebuild_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_current_sql := CASE
	WHEN @needs_rebuild = 1 AND @source_table = 'delivery_state' AND @legacy_exists > 0 THEN 'DROP TABLE delivery_state'
	WHEN @needs_rebuild = 1 AND @source_table = 'delivery_state_legacy_000003' THEN 'DROP TABLE IF EXISTS delivery_state'
	ELSE 'SELECT 1'
END;
PREPARE stmt FROM @drop_current_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @rename_sql := CASE
	WHEN @needs_rebuild = 1 AND @source_table = 'delivery_state' AND @legacy_exists = 0 THEN
		'RENAME TABLE delivery_state TO delivery_state_legacy_000003, delivery_state_rebuild_000003 TO delivery_state'
	WHEN @needs_rebuild = 1 AND (
		@source_table = 'delivery_state_legacy_000003'
		OR (@source_table = 'delivery_state' AND @legacy_exists > 0)
	) THEN
		'RENAME TABLE delivery_state_rebuild_000003 TO delivery_state'
	ELSE 'SELECT 1'
END;
PREPARE stmt FROM @rename_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

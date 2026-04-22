ALTER TABLE delivery_state
	DROP PRIMARY KEY,
	DROP INDEX idx_delivery_state_operation_cursor,
	DROP INDEX idx_delivery_state_entity_key;

ALTER TABLE delivery_state
	CHANGE COLUMN source_cursor source_id BIGINT NOT NULL,
	CHANGE COLUMN entity_key dedupe_key VARCHAR(255) NOT NULL;

ALTER TABLE delivery_state
	ADD PRIMARY KEY (operation_name, source_id),
	ADD KEY idx_delivery_state_dedupe_key (dedupe_key);

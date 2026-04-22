ALTER TABLE delivery_state
	CHANGE COLUMN source_id source_cursor BIGINT NOT NULL,
	CHANGE COLUMN dedupe_key entity_key VARCHAR(255) NOT NULL;

ALTER TABLE delivery_state
	DROP PRIMARY KEY,
	ADD PRIMARY KEY (operation_name, entity_key);

ALTER TABLE delivery_state
	DROP INDEX idx_delivery_state_dedupe_key,
	ADD KEY idx_delivery_state_entity_key (entity_key),
	ADD KEY idx_delivery_state_operation_cursor (operation_name, source_cursor);

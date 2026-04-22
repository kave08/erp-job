package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"erp-job/internal/store"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	getSourceProgressQuery = `
SELECT last_source_id
FROM source_progress
WHERE entity = ?;
`
	upsertSourceProgressQuery = `
INSERT INTO source_progress (entity, last_source_id, created_at, updated_at)
VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	last_source_id = GREATEST(last_source_id, VALUES(last_source_id)),
	updated_at = VALUES(updated_at);
`
	getOperationCheckpointQuery = `
SELECT last_source_id
FROM operation_checkpoint
WHERE operation_name = ?;
`
	upsertOperationCheckpointQuery = `
INSERT INTO operation_checkpoint (operation_name, last_source_id, created_at, updated_at)
VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	last_source_id = GREATEST(last_source_id, VALUES(last_source_id)),
	updated_at = VALUES(updated_at);
`
	upsertDeliveryAttemptQuery = `
INSERT INTO delivery_state (
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
VALUES (?, ?, ?, ?, 1, ?, ?, NULL, ?, ?)
ON DUPLICATE KEY UPDATE
	source_cursor = GREATEST(delivery_state.source_cursor, VALUES(source_cursor)),
	status = VALUES(status),
	attempt_count = delivery_state.attempt_count + 1,
	last_http_status = VALUES(last_http_status),
	last_error = VALUES(last_error),
	updated_at = VALUES(updated_at);
`
	upsertDeliveredRecordQuery = `
INSERT INTO delivery_state (
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
VALUES (?, ?, ?, ?, 1, ?, NULL, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	source_cursor = GREATEST(delivery_state.source_cursor, VALUES(source_cursor)),
	status = VALUES(status),
	attempt_count = delivery_state.attempt_count,
	last_http_status = VALUES(last_http_status),
	last_error = NULL,
	delivered_at = VALUES(delivered_at),
	updated_at = VALUES(updated_at);
`
	getAttemptCountsQuery = `
SELECT entity_key, attempt_count
FROM delivery_state
WHERE operation_name = ?
  AND entity_key IN (%s);
`
	markPermanentFailuresQuery = `
UPDATE delivery_state
SET status = ?,
    updated_at = ?
WHERE operation_name = ?
  AND entity_key IN (%s);
`
)

var tracer = otel.Tracer("erp-job/store/mysql")

type Checkpoints struct {
	db *sql.DB
}

func New(db *sql.DB) *Checkpoints {
	return &Checkpoints{db: db}
}

func (c *Checkpoints) GetSourceProgress(ctx context.Context, entity store.Entity) (int, error) {
	ctx, span := tracer.Start(ctx, "store.get_source_progress")
	defer span.End()

	span.SetAttributes(attribute.String("entity", string(entity)))

	var lastSourceID sql.NullInt64
	err := c.db.QueryRowContext(ctx, getSourceProgressQuery, string(entity)).Scan(&lastSourceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		recordSpanError(span, err)
		return 0, fmt.Errorf("get source progress for %s: %w", entity, err)
	}

	if !lastSourceID.Valid {
		return 0, nil
	}

	return int(lastSourceID.Int64), nil
}

func (c *Checkpoints) AdvanceSourceProgress(ctx context.Context, entity store.Entity, lastSourceID int) error {
	ctx, span := tracer.Start(ctx, "store.advance_source_progress")
	defer span.End()

	span.SetAttributes(
		attribute.String("entity", string(entity)),
		attribute.Int("last_source_id", lastSourceID),
	)

	now := time.Now().UTC()
	if _, err := c.db.ExecContext(ctx, upsertSourceProgressQuery, string(entity), lastSourceID, now, now); err != nil {
		recordSpanError(span, err)
		return fmt.Errorf("advance source progress for %s: %w", entity, err)
	}

	return nil
}

func (c *Checkpoints) GetOperationCheckpoint(ctx context.Context, operation store.Operation) (int, error) {
	ctx, span := tracer.Start(ctx, "store.get_operation_checkpoint")
	defer span.End()

	span.SetAttributes(attribute.String("operation", string(operation)))

	var lastSourceID sql.NullInt64
	err := c.db.QueryRowContext(ctx, getOperationCheckpointQuery, string(operation)).Scan(&lastSourceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		recordSpanError(span, err)
		return 0, fmt.Errorf("get operation checkpoint for %s: %w", operation, err)
	}

	if !lastSourceID.Valid {
		return 0, nil
	}

	return int(lastSourceID.Int64), nil
}

func (c *Checkpoints) RecordDeliveryAttempt(ctx context.Context, attempt store.DeliveryAttempt) error {
	ctx, span := tracer.Start(ctx, "store.record_delivery_attempt")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", string(attempt.Operation)),
		attribute.Int("source_cursor", attempt.SourceCursor),
		attribute.String("entity_key", attempt.EntityKey),
		attribute.String("status", string(attempt.Status)),
		attribute.Int("http_status", attempt.HTTPStatus),
	)

	now := attempt.AttemptedAt.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}

	var lastHTTPStatus interface{}
	if attempt.HTTPStatus > 0 {
		lastHTTPStatus = attempt.HTTPStatus
	}

	if _, err := c.db.ExecContext(
		ctx,
		upsertDeliveryAttemptQuery,
		string(attempt.Operation),
		attempt.SourceCursor,
		attempt.EntityKey,
		string(attempt.Status),
		lastHTTPStatus,
		nullableString(attempt.Error),
		now,
		now,
	); err != nil {
		recordSpanError(span, err)
		return fmt.Errorf("record delivery attempt for %s:%s: %w", attempt.Operation, attempt.EntityKey, err)
	}

	return nil
}

func (c *Checkpoints) GetDeliveredEntityKeys(ctx context.Context, operation store.Operation, entityKeys []string) (map[string]struct{}, error) {
	ctx, span := tracer.Start(ctx, "store.get_delivered_entity_keys")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", string(operation)),
		attribute.Int("entity_key_count", len(entityKeys)),
	)

	delivered := make(map[string]struct{}, len(entityKeys))
	if len(entityKeys) == 0 {
		return delivered, nil
	}

	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(entityKeys)), ",")
	query := fmt.Sprintf(`
SELECT entity_key
FROM delivery_state
WHERE operation_name = ?
  AND status = ?
  AND entity_key IN (%s);
`, placeholders)

	args := make([]interface{}, 0, len(entityKeys)+2)
	args = append(args, string(operation), string(store.DeliveryStatusDelivered))
	for _, entityKey := range entityKeys {
		args = append(args, entityKey)
	}

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		recordSpanError(span, err)
		return nil, fmt.Errorf("get delivered entity keys for %s: %w", operation, err)
	}
	defer rows.Close()

	for rows.Next() {
		var entityKey string
		if err := rows.Scan(&entityKey); err != nil {
			recordSpanError(span, err)
			return nil, fmt.Errorf("scan delivered entity key for %s: %w", operation, err)
		}
		delivered[entityKey] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		recordSpanError(span, err)
		return nil, fmt.Errorf("iterate delivered entity keys for %s: %w", operation, err)
	}

	return delivered, nil
}

func (c *Checkpoints) MarkBatchDelivered(ctx context.Context, operation store.Operation, lastSourceID int, records []store.DeliveredRecord) error {
	ctx, span := tracer.Start(ctx, "store.mark_batch_delivered")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", string(operation)),
		attribute.Int("last_source_id", lastSourceID),
		attribute.Int("record_count", len(records)),
	)

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		recordSpanError(span, err)
		return fmt.Errorf("begin delivery transaction for %s: %w", operation, err)
	}

	if err := c.markBatchDelivered(ctx, tx, operation, lastSourceID, records); err != nil {
		_ = tx.Rollback()
		recordSpanError(span, err)
		return err
	}

	if err := tx.Commit(); err != nil {
		recordSpanError(span, err)
		return fmt.Errorf("commit delivery transaction for %s: %w", operation, err)
	}

	return nil
}

func (c *Checkpoints) markBatchDelivered(ctx context.Context, tx *sql.Tx, operation store.Operation, lastSourceID int, records []store.DeliveredRecord) error {
	now := time.Now().UTC()

	for _, record := range records {
		deliveredAt := record.DeliveredAt.UTC()
		if deliveredAt.IsZero() {
			deliveredAt = now
		}

		var lastHTTPStatus interface{}
		if record.HTTPStatus > 0 {
			lastHTTPStatus = record.HTTPStatus
		}

		if _, err := tx.ExecContext(
			ctx,
			upsertDeliveredRecordQuery,
			string(operation),
			record.SourceCursor,
			record.EntityKey,
			string(store.DeliveryStatusDelivered),
			lastHTTPStatus,
			deliveredAt,
			now,
			now,
		); err != nil {
			return fmt.Errorf("upsert delivered state for %s:%s: %w", operation, record.EntityKey, err)
		}
	}

	if _, err := tx.ExecContext(ctx, upsertOperationCheckpointQuery, string(operation), lastSourceID, now, now); err != nil {
		return fmt.Errorf("advance operation checkpoint for %s: %w", operation, err)
	}

	return nil
}

func nullableString(value string) interface{} {
	if value == "" {
		return nil
	}

	return value
}

func (c *Checkpoints) GetAttemptCounts(ctx context.Context, operation store.Operation, entityKeys []string) (map[string]int, error) {
	ctx, span := tracer.Start(ctx, "store.get_attempt_counts")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", string(operation)),
		attribute.Int("entity_key_count", len(entityKeys)),
	)

	counts := make(map[string]int, len(entityKeys))
	if len(entityKeys) == 0 {
		return counts, nil
	}

	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(entityKeys)), ",")
	query := fmt.Sprintf(getAttemptCountsQuery, placeholders)

	args := make([]interface{}, 0, len(entityKeys)+1)
	args = append(args, string(operation))
	for _, entityKey := range entityKeys {
		args = append(args, entityKey)
	}

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		recordSpanError(span, err)
		return nil, fmt.Errorf("get attempt counts for %s: %w", operation, err)
	}
	defer rows.Close()

	for rows.Next() {
		var entityKey string
		var attemptCount int
		if err := rows.Scan(&entityKey, &attemptCount); err != nil {
			recordSpanError(span, err)
			return nil, fmt.Errorf("scan attempt count for %s: %w", operation, err)
		}
		counts[entityKey] = attemptCount
	}

	if err := rows.Err(); err != nil {
		recordSpanError(span, err)
		return nil, fmt.Errorf("iterate attempt counts for %s: %w", operation, err)
	}

	return counts, nil
}

func (c *Checkpoints) MarkPermanentFailures(ctx context.Context, operation store.Operation, entityKeys []string) error {
	ctx, span := tracer.Start(ctx, "store.mark_permanent_failures")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", string(operation)),
		attribute.Int("entity_key_count", len(entityKeys)),
	)

	if len(entityKeys) == 0 {
		return nil
	}

	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(entityKeys)), ",")
	query := fmt.Sprintf(markPermanentFailuresQuery, placeholders)

	args := make([]interface{}, 0, len(entityKeys)+3)
	args = append(args, string(store.DeliveryStatusPermanentFailure), time.Now().UTC(), string(operation))
	for _, entityKey := range entityKeys {
		args = append(args, entityKey)
	}

	if _, err := c.db.ExecContext(ctx, query, args...); err != nil {
		recordSpanError(span, err)
		return fmt.Errorf("mark permanent failures for %s: %w", operation, err)
	}

	return nil
}

func recordSpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

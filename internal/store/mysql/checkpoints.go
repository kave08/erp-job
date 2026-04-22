package mysql

import (
	"context"
	"database/sql"
	"fmt"
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
	source_id,
	dedupe_key,
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
	dedupe_key = VALUES(dedupe_key),
	status = VALUES(status),
	attempt_count = delivery_state.attempt_count + 1,
	last_http_status = VALUES(last_http_status),
	last_error = VALUES(last_error),
	updated_at = VALUES(updated_at);
`
	upsertDeliveredRecordQuery = `
INSERT INTO delivery_state (
	operation_name,
	source_id,
	dedupe_key,
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
	dedupe_key = VALUES(dedupe_key),
	status = VALUES(status),
	attempt_count = GREATEST(delivery_state.attempt_count, VALUES(attempt_count)),
	last_http_status = VALUES(last_http_status),
	last_error = NULL,
	delivered_at = VALUES(delivered_at),
	updated_at = VALUES(updated_at);
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
		attribute.Int("source_id", attempt.SourceID),
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
		attempt.SourceID,
		attempt.DedupeKey,
		string(attempt.Status),
		lastHTTPStatus,
		nullableString(attempt.Error),
		now,
		now,
	); err != nil {
		recordSpanError(span, err)
		return fmt.Errorf("record delivery attempt for %s:%d: %w", attempt.Operation, attempt.SourceID, err)
	}

	return nil
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
			record.SourceID,
			record.DedupeKey,
			string(store.DeliveryStatusDelivered),
			lastHTTPStatus,
			deliveredAt,
			now,
			now,
		); err != nil {
			return fmt.Errorf("upsert delivered state for %s:%d: %w", operation, record.SourceID, err)
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

func recordSpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

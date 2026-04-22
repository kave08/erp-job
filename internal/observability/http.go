package observability

import (
	"context"
	"net/http"

	"erp-job/internal/retry"

	"go.uber.org/zap"
)

func LogHTTPAttempt(log *zap.SugaredLogger, ctx context.Context, system, endpointGroup string, attempt retry.Attempt) {
	if log == nil {
		return
	}

	fields := []interface{}{
		"run_id", RunIDFromContext(ctx),
		"system", system,
		"endpoint_group", endpointGroup,
		"attempt", attempt.Attempt,
		"status_code", attempt.StatusCode,
		"duration_ms", attempt.Duration.Milliseconds(),
		"will_retry", attempt.WillRetry,
	}

	if attempt.Error != nil {
		fields = append(fields, "error", attempt.Error.Error(), "error_class", ClassifyHTTPError(attempt.StatusCode, attempt.Error))
		log.Warnw(system+" request attempt failed", fields...)
		return
	}

	log.Infow(system+" request succeeded", fields...)
}

func ClassifyHTTPError(statusCode int, err error) string {
	switch {
	case err == nil:
		return "none"
	case statusCode == http.StatusTooManyRequests:
		return "rate_limit"
	case statusCode >= 500:
		return "upstream_5xx"
	case statusCode >= 400:
		return "upstream_4xx"
	default:
		return "transport"
	}
}

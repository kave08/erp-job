# ERP Job

ERP Job is a one-shot sync CLI for moving operational data from Fararavand to Aryan. Each run fetches new records, sends mapped payloads to Aryan, and stores progress in MySQL so the next run can continue incrementally instead of replaying the full dataset.

- Syncs invoices, products, customers, and base data
- Persists cursor and delivery state in MySQL for resumable, duplicate-safe runs
- Runs locally or as a containerized job

## Quick Start

Create `env.yml` from `env.example.yml`, then fill in the Fararavand, Aryan, MySQL, and optional OTLP collector settings.

```sh
cp env.example.yml env.yml
go build .
./erp-job migrate --config-path env.yml
./erp-job transfer --config-path env.yml
```

## Docker

The container runs the same one-shot `transfer` command and exits when the sync completes. `compose.yaml` mounts `env.yml` into `/config/env.yml`.

```sh
docker compose run --rm erp-job
```

## Notes

- This project is a job runner, not an HTTP server.
- Sync progress and per-operation delivery attempts are checkpointed in MySQL.
- `transfer` applies pending migrations automatically; `migrate` is available for explicit DB bootstrap.
- When `OTEL.ENABLED=true`, the app exports OTLP traces and metrics to the configured collector endpoint.
- The application code lives under `internal/`, with `internal/app` as the wiring entrypoint.

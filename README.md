# ERP Job

ERP Job is a command-line tool that synchronizes data between Fararavand and Aryan ERP systems. It transfers invoices, products, customers, and base data while persisting sync checkpoints in MySQL so repeated runs can continue incrementally.

## Features

- Synchronize invoices, products, and customer data between ERP systems.
- Convert data models from one system's format to another.
- Post data to various endpoints, including sale orders, sale centers, and payment selections.
- Handle database operations for tracking synchronization status.

## Getting Started

### Prerequisites

- Go (version 1.x)
- MySQL (version 5.7 or higher)
- Access to ERP systems' APIs

### Build

```sh
go build .
```

### Configuration

Before running the application, create `env.yml` from `env.example.yml` and fill in the ERP and MySQL credentials.

### Usage

Run the sync job with:

```sh
./erp-job transfer --config-path env.yml
```

## Architecture

The project is organized as an application, not a reusable library, so implementation packages now live under `internal/`:

- `internal/app`: application wiring and startup flow
- `internal/config`: config loading from YAML and environment overrides
- `internal/database`: MySQL connection setup
- `internal/domain`: Fararavand source models and Aryan target payload models
- `internal/source/fararavand`: HTTP client for paged Fararavand reads
- `internal/target/aryan`: HTTP client for Aryan write operations
- `internal/store` and `internal/store/mysql`: checkpoint interfaces and MySQL persistence
- `internal/transfer`: sync orchestration and checkpoint-aware batch processing

## Docker

This application runs as a one-shot sync job, so the container starts `erp-job transfer` and exits when the sync finishes.

### Build the image

```sh
docker build -t erp-job:local .
```

### Run with a mounted config file

```sh
docker run --rm \
  --read-only \
  --tmpfs /tmp \
  -v "$(pwd)/env.yml:/config/env.yml:ro" \
  erp-job:local
```

### Run with Docker Compose

`compose.yaml` mounts `./env.yml` into `/config/env.yml` and keeps the container filesystem read-only except for `/tmp`.

```sh
docker compose run --rm erp-job
```

Use [`env.example.yml`](env.example.yml) as the template for the mounted configuration. In containers, `APP.LOG_PATH: "/tmp/erp-job-logs"` is a safe default if file logs are needed; stdout logging is enabled either way.


## Contributing

Contributions are welcome! Please feel free to submit pull requests or create issues for bugs and feature requests.

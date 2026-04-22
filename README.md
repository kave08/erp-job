# ERP Job

ERP Job is a command-line tool designed to facilitate the synchronization of data between different ERP (Enterprise Resource Planning) systems. It provides a set of functionalities to transfer and transform data such as invoices, products, and customer information from one system to another, ensuring data consistency and integrity.

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

### Installation

1. Clone the repository:
sh git clone https://github.com/kave08/erp-job.git

2. Navigate to the project directory:
sh cd erp-job

3. Build the project:

### Configuration

Before running the application, configure the necessary API keys and database connection strings in the `config` directory.

### Usage

Run the application with the following command:

sh ./erp-job [transfer]

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


## Modules

### Fararavand

The `fararavand.go` module contains the logic for interacting with the Fararavand ERP system. It defines an interface for the operations that can be performed and a struct that implements these operations.

### Aryan

The `aryan.go` module contains the logic for interacting with the Aryan ERP system. Similar to the Fararavand module, it defines an interface and a struct for the Aryan-specific operations.

### Database

The `database.go` module in the `repository/database` directory handles all database-related operations. It provides functions to insert and retrieve the maximum ID values for various synchronization operations, ensuring that only new data is transferred between systems.

## Contributing

Contributions are welcome! Please feel free to submit pull requests or create issues for bugs and feature requests.

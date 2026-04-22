# Code Review

## Critical Findings

1. `erp-job transfer` does not execute any sync work.
   `cmd/transfer.go:22-31` only loads config, builds dependencies, and calls `syncdata.NewSync(...)` without invoking any sync method. `sync_data/sync.go:31-39` is also inert: it only constructs handler objects and never calls `Invoices()`, `Customers()`, `Products()`, or `BaseData()`. As written, the CLI exits without syncing anything.

2. Sale-factor sync posts the wrong payload.
   `services/aryan/aryan.go:48-150` builds `req []Req`, but `json.Marshal` is called on `newSaleFactor`, which is never populated. The request body sent to the Aryan sale-factor endpoint is therefore `null`, not the invoice payload assembled in the loop.

## High Findings

3. Invoice and base-data deduplication is broken, so already-synced rows are reposted.
   In `services/fararavand/fararavand.go:144-151`, `192-199`, `241-247`, `289-295`, `339-345`, `388-394`, `436-442`, and `485-491`, each loop compares items against the current batch max (`lastInvoiceID`/`lastInvoiceID`) instead of the persisted checkpoint (`lastSaleFactorID`, `lastSaleOrderID`, and so on). Those inner conditions can never become true for items from the same batch, so old rows are never trimmed before reposting.

4. Resume checkpoints read the oldest saved row, not the latest one.
   `repository/database/database.go:13-20` uses `ORDER BY id LIMIT 1` for invoice, customer, product, and base-data progress tables. The insert methods append new checkpoint rows, so restarts resume from the first checkpoint ever written rather than the most recent one.

5. Several DB insert helpers will fail at runtime because they bind too few SQL parameters.
   `repository/database/database.go:23`, `43-45`, and `613-636` define SQL statements with two placeholders (`VALUES(?, ?)`), but `InsertProductsToGoods`, `InsertInvoiceReturn`, and `InsertTreasuries` each pass only one argument to `Exec`. `services/fararavand/fararavand.go:105-115` already routes product sync into `InsertProductsToGoods`, so the product path will fail as soon as it reaches the checkpoint write.

6. Error handling can panic because the shared logger is never initialized.
   `utility/logger/logger.go:15-42` defines `Initialize()`, but there is no call site in the repository. Constructors such as `services/fararavand/fararavand.go:22-28` and `sync_data/invoice.go:33-43` store `logger.Logger()` immediately, which leaves `log` nil. On the first error path, calls like `i.log.Errorw(...)` will panic instead of reporting the original failure.

7. Config loading leaks secrets to stdout and panic output.
   `config/config.go:77-81` prints the full config struct and includes `Cfg.Database` in a panic message. That exposes API keys, ERP credentials, and database passwords in logs or terminal output.

## Medium Findings

8. Paged Fararavand fetchers leak HTTP response bodies.
   `sync_data/invoice.go:75-97`, `sync_data/customer.go:75-91`, `sync_data/product.go:77-99`, and `sync_data/baseData.go:76-97` decode `res.Body` but never close it. These functions run in loops, so repeated pages can consume connections until sync calls start failing.

9. Invoice-return and treasury flows are still stubs.
   `sync_data/invoiceReturns.go:36-38`, `sync_data/treasuries.go:36-38`, and `services/fararavand/fararavand.go:521-535` all return `nil` without fetching, transforming, or posting any data. Even if the main sync orchestration is fixed, these two domains still do nothing.

## Testing Gaps

- `go test ./...` passes only because the repository has no `_test.go` files.
- There is no integration harness or fixture config in the repo to exercise the Cobra command, the MySQL checkpoint tables, or the Fararavand/Aryan HTTP contracts safely.

## Checks Run

- `go test ./...`
- `go vet ./...`

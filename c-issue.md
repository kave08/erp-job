Critical Issues & Edge Cases
1. All-or-Nothing Batch Failure (Major)
- Issue: If a batch of 1000 invoices is fetched and PostInvoiceToSaleFactor succeeds for all but one invoice due to a target-side validation error (e.g., invalid VisitorCode), the entire batch fails and will be retried indefinitely. The successful items in that batch will be re-sent on every retry.
- Impact: This can cause:
  - Duplicate data in Aryan for operations that did succeed before the failure.
  - Infinite loop if one record is permanently bad (poison pill).
- Evidence: deliverBatch calls post(observedCtx, itemsFromCandidates(prepared.pending)) with the full pending slice. If post returns an error, MarkBatchDelivered is never called, so the checkpoint doesn't advance.
2. No Per-Record Error Handling / Dead Letter Queue
- Issue: There is no mechanism to skip or quarantine permanently failing records. A single bad record blocks the entire pipeline for that entity.
- Recommendation: Introduce a DeliveryStatusSkipped or DeliveryStatusPermanentFailure state after N retries, with a dead-letter table for manual review.
3. Invoice Operation Dependency Risk
- Issue: Invoices have 7 sub-operations (sale_payment, sale_center, saler_select, sale_type_select, sale_factor, sale_order, sale_proforma). These are executed sequentially. If sale_factor succeeds but sale_order fails:
  - source_progress for invoice does NOT advance (correct).
  - But sale_factor checkpoint DOES advance (correct for idempotency).
  - However, on retry, sale_factor will re-fetch the same batch, deduplicate by entity_key, and skip already-delivered items. This is fine.
  - BUT: What if Aryan is not idempotent? If Aryan receives the same SaleFactor payload again with the same SecondNumber, does it create a duplicate or reject it? The code assumes idempotency, but this is not verified.
4. Entity Key Collisions in Invoices
- Issue: The entityKey for sale_factor, sale_order, and sale_proforma is:
    fmt.Sprintf("invoice:%d:product:%d", item.InvoiceId, item.ProductID)
    This implies an invoice can have multiple line items with the same InvoiceId but different ProductIDs. However, if the source API returns multiple rows for the same InvoiceId + ProductID combination (e.g., due to a data issue or partial shipment), the deduplication will silently drop duplicates, potentially losing data.
- Question: Does Fararavand.GetInvoices guarantee uniqueness on (InvoiceId, ProductID)? If not, this is a data loss risk.
5. BaseData Sync Only Syncs PaymentTypes
- Issue: FetchBaseData returns a rich BaseData struct with CustomerTypes, GiuldTypes, SanadTypes, Units, Branches, Brands, Areas, Districts, States, Cities, WareHouses, FirstProductGroup, SecondProductGroup.
- But runBaseData only syncs PaymentTypes to DeliverCenter_SaleSelect. All other base data is fetched but discarded.
- Impact: If Aryan needs other base data (e.g., WareHouses for SaleCenter), the mappings may fail or be incomplete.
6. context.WithoutCancel in Delivery Observer (Minor)
- Issue: withDeliveryObserver uses context.WithoutCancel(ctx) to persist attempts after the main context is cancelled.
- Risk: If the DB is slow or the process is shutting down, this can hang indefinitely because there's no timeout on recordCtx.
- Recommendation: Wrap in a context.WithTimeout(context.Background(), 5*time.Second).
7. InvoiceDate vs Date Field Confusion
- Issue: In PostInvoiceToSaleFactor, VoucherDate is mapped from item.InvoiceDate (string).
- In PostInvoiceToSaleOrder, VoucherDate is also mapped from item.InvoiceDate (string).
- In PostInvoiceToSaleProforma, VoucherDate is mapped from item.Date (time.Time).
- Question: Why the inconsistency? Is InvoiceDate a string format that Aryan expects, while Date is a time.Time? This could cause serialization issues or date format mismatches.
8. SaleProforma.SecondNumber is string, but mapped from strconv.Itoa(item.InvoiceId)
- Issue: In PostInvoiceToSaleProforma:
    SecondNumber: strconv.Itoa(item.InvoiceId),
    But in SaleOrder, SecondNumber is int and mapped directly from item.InvoiceId.
- Question: Does Aryan expect different types for the same semantic field across endpoints? If Aryan's API is strict, this will cause deserialization errors on their side.
9. No Handling of HTTP 2xx != 200
- Issue: Both source and target clients only accept StatusCode == http.StatusOK (200).
- Edge Case: If Aryan returns 201 Created or 204 NoContent for some endpoints, the client will treat it as a failure and retry.
- Recommendation: Accept statusCode >= 200 && statusCode < 300.
10. Missing Content-Type Check on Source Responses
- Issue: fararavand/client.go decodes JSON without checking Content-Type. If the API returns an HTML error page (e.g., 502 Bad Gateway from a proxy), json.NewDecoder will fail with an unhelpful error.
11. PageNumeber Typo in API Paths
- Issue: The query parameter is spelled PageNumeber (missing 'b') in all four paths:
    getProductsPath  = "/GetProducts?PageNumeber=%d&PageSize=%d&LastId=%d/"
  - Impact: If the Fararavand API is case-sensitive or strict about parameter names, pagination may not work correctly. The code always passes cursorPageNumber = 0, so this may be masked currently, but it's a latent bug if offset-based pagination is ever needed.
12. runEntitySync Assumes Monotonically Increasing IDs
- Issue: The sync logic relies on lastID pagination:
    items, err := fetch(ctx, cursorPageNumber, j.pageSize, lastSourceID)
    This assumes that InvoiceId, Customer.ID, Product.ID, and PaymentType.ID are strictly monotonically increasing and that new records always have higher IDs than old ones.
- Edge Case: If IDs are UUIDs, non-sequential, or if records can be backdated/inserted with lower IDs, this will miss data.
- Question: Is this assumption valid for Fararavand's data model?
13. source_progress Advancement is Coarse
- Issue: source_progress advances by the batchLastID (the highest ID in the fetched batch). If a batch of 1000 records is processed and 999 operations succeed but the last one fails, the entire batch is retried on the next run.
- This is by design (all-or-nothing per batch), but it means the job is not as incremental as it could be.
14. No Circuit Breaker
- Issue: If Aryan is down for an extended period, the job will keep retrying every run (with backoff per-request, but not across runs).
- Recommendation: Consider a circuit breaker pattern to fail fast and reduce load on a struggling target.
15. SaleFactor JSON Tags are PascalCase, Others are snake_case
- Issue: In domain/aryan.go:
  - SaleFactor uses PascalCase tags: CustomerId, VoucherDate, StockID, etc.
  - SaleOrder uses snake_case tags: customer_id, voucher_date, stock_id, etc.
- Impact: If Aryan's API is consistent in expecting one convention, half the endpoints will fail. This looks like a bug.
16. PostBaseDataToDeliverCenterSaleSelect Payload Construction
- Issue: The entity_key for base data is fmt.Sprintf("payment_type:%d", item.ID), but DeliverCenterSaleSelect only maps CentersID and CentersCode. The Name from PaymentTypes is lost.
- Question: Is CentersCode = strconv.Itoa(item.ID) correct? Should it be the Name?
---
Questions for You
Before I proceed with any fixes, I need clarification on a few points:
1. Aryan API Idempotency: Does Aryan handle duplicate POST requests gracefully (e.g., via SecondNumber uniqueness constraints), or will it create duplicates?
2. Date Fields: Why does SaleProforma use item.Date (time.Time) while SaleFactor and SaleOrder use item.InvoiceDate (string)? What format does Aryan expect?
3. Base Data Scope: Should base_data sync include WareHouses, CustomerTypes, etc., or is PaymentTypes sufficient for now?
4. Monotonic IDs: Can we rely on Fararavand IDs always increasing, or should we consider updated_at / cursor-based pagination?
5. Poison Pills: What should happen if a single record permanently fails (e.g., invalid data)? Should it block the pipeline, or should we skip it after N retries?
---
Recommended Priority of Fixes
Priority	Issue
P0	Fix SaleFactor JSON tag casing inconsistency
P0	Accept 2xx status codes, not just 200
P1	Add per-record dead-letter handling for permanent failures
P1	Clarify/fix InvoiceDate vs Date mapping
P1	Fix PageNumeber typo
P2	Add circuit breaker for target health
P2	Add Content-Type validation on responses
P2	Scope BaseData sync to all required fields
---
Would you like me to start implementing any of these fixes, or would you prefer to address the clarifying questions first?

-----------------------------------------------------------------------------

Implementation Plan
Phase 1: Critical Bugs & Data Integrity (P0)
1.1 Fix SaleFactor JSON Tag Casing Inconsistency
- File: internal/domain/aryan.go
- Problem: SaleFactor uses PascalCase tags (CustomerId, VoucherDate, StockID...) while every other Aryan model uses snake_case (customer_id, voucher_date, stock_id...). This will cause Aryan's API to miss fields or fail validation.
- Fix: Change all SaleFactor JSON tags to snake_case to match the convention used by SaleOrder, SaleProforma, etc.
1.2 Accept All 2xx HTTP Status Codes
- Files: internal/source/fararavand/client.go, internal/target/aryan/client.go
- Problem: Both clients only accept http.StatusOK (200). If Aryan returns 201 Created or Fararavand returns 204 NoContent, the job treats it as a failure and retries.
- Fix: Change res.StatusCode != http.StatusOK to res.StatusCode < 200 || res.StatusCode >= 300 in both getOnce and postJSONOnce.
1.3 Fix PageNumeber → PageNumber Typo
- File: internal/source/fararavand/client.go
- Problem: The query parameter is misspelled as PageNumeber in all four endpoint paths. If Fararavand's API is strict, pagination is broken.
- Fix: Correct to PageNumber in getProductsPath, getInvoicesPath, getCustomersPath, getBaseDataPath.
1.4 Fix SaleProforma Date Mapping Inconsistency
- File: internal/target/aryan/client.go
- Problem: PostInvoiceToSaleProforma maps VoucherDate from item.Date (time.Time), while SaleFactor and SaleOrder map from item.InvoiceDate (string). SaleProforma.VoucherDate is declared as time.Time in domain/aryan.go, but Aryan may expect a string format consistent with other endpoints.
- Fix: Change SaleProforma.VoucherDate to string and map from item.InvoiceDate for consistency. (Alternative: if Aryan truly needs time.Time for proforma, format it explicitly — but consistency is safer).
1.5 Fix SaleProforma.SecondNumber Type Inconsistency
- File: internal/domain/aryan.go
- Problem: SaleProforma.SecondNumber is string, but SaleOrder.SecondNumber is int. The client maps strconv.Itoa(item.InvoiceId) for proforma but item.InvoiceId directly for order.
- Fix: Change SaleProforma.SecondNumber to int for consistency, and map directly from item.InvoiceId.
---
Phase 2: Resilience & Edge Cases (P1)
2.1 Add Content-Type Validation on Source Responses
- File: internal/source/fararavand/client.go
- Problem: If a proxy returns HTML (e.g., 502 error page), json.NewDecoder fails with an obscure error.
- Fix: After checking status code, verify Content-Type header contains application/json. If not, return a clear error with the raw body snippet.
2.2 Add Timeout to context.WithoutCancel in Delivery Observer
- File: internal/transfer/job.go
- Problem: withDeliveryObserver uses context.WithoutCancel(ctx) which has no deadline. If the DB is slow during shutdown, this can hang.
- Fix: Wrap with context.WithTimeout(context.Background(), 5*time.Second).
2.3 Harden trimAfterCheckpoint for Unsorted Batches
- File: internal/transfer/job.go
- Problem: trimAfterCheckpoint assumes items are sorted by ID ascending. If Fararavand ever returns out-of-order data, the logic breaks.
- Fix: Add a defensive check / log warning if the batch is not monotonically increasing. (Or document the invariant explicitly).
2.4 Add DeliveryStatusPermanentFailure and Max Attempt Tracking
- Files: internal/store/checkpoints.go, internal/store/mysql/checkpoints.go, internal/transfer/job.go
- Problem: A permanently bad record (poison pill) blocks the entire pipeline indefinitely.
- Fix:
  - Add DeliveryStatusPermanentFailure to DeliveryStatus enum.
  - Add max_attempts config (default 3) to the job.
  - In deliverBatch, if post fails and the retry policy exhausts, check attempt_count in delivery_state. If >= max_attempts, mark as permanent_failure, log loudly, and advance the checkpoint for that entity key only.
  - Add a dead_letter table migration for manual review of permanent failures.
---
Phase 3: Observability & Operational Hardening (P2)
3.1 Add Circuit Breaker for Target Client
- File: internal/target/aryan/client.go (new file: internal/circuitbreaker/circuitbreaker.go)
- Problem: If Aryan is down, every run hammers it with retries.
- Fix: Implement a simple in-memory circuit breaker (closed → open → half-open) based on consecutive failure threshold. When open, fail fast without calling the API.
3.2 Add Last-Source-ID / Batch Metadata to Logs
- File: internal/transfer/job.go
- Problem: When debugging, it's hard to correlate a failed batch with the exact source cursor range.
- Fix: Add batch_first_id, batch_last_id, and checkpoint_before to all step-level log entries.
3.3 Export checkpointLag as Gauge (not just Histogram)
- File: internal/observability/otel.go
- Problem: RecordCheckpointLag uses a histogram, which is great for distribution but hard to alert on "current lag".
- Fix: Add an Int64ObservableGauge for erp_job_current_checkpoint_lag per entity, updated after each batch.
---
Phase 4: Base Data Completeness (P2)
4.1 Expand BaseData Sync or Document Scope
- File: internal/transfer/job.go
- Problem: runBaseData only syncs PaymentTypes. Other base data fields (WareHouses, CustomerTypes, Units, etc.) are fetched but discarded.
- Decision Point: Do we need to sync all base data types to Aryan? If yes, we need new Aryan endpoints and mappings. If no, we should add a comment documenting why PaymentTypes is sufficient.
- Fix (Minimal): Add a // TODO or explicit comment explaining that only PaymentTypes are mapped to DeliverCenter_SaleSelect and other base data is intentionally out of scope.
---
Files to Modify (Summary)
File
internal/domain/aryan.go
internal/target/aryan/client.go
internal/source/fararavand/client.go
internal/retry/retry.go
internal/transfer/job.go
internal/store/checkpoints.go
internal/store/mysql/checkpoints.go
internal/observability/otel.go
db/migrations/000004_dead_letter.up.sql
internal/circuitbreaker/circuitbreaker.go
---
Questions Before I Start
1. Aryan API Convention: Can you confirm Aryan expects snake_case JSON for ALL endpoints? The SaleFactor PascalCase tags look like a copy-paste error from a different API spec.
2. Date Format: What string format does Aryan expect for voucher_date? (e.g., "2006-01-02", "2006/01/02", ISO8601?) This affects whether time.Time or string is correct.
3. Poison Pills: Should permanently failing records be skipped (advance checkpoint) or should they block the pipeline until manually fixed?
4. Base Data Scope: Is syncing only PaymentTypes to DeliverCenter_SaleSelect intentional, or do we need to map WareHouses → SaleCenter4SaleSelect, etc.?
Please let me know your preferences, and I'll begin implementation immediately.
▣  Plan · Kimi K2.5 · 48.0s
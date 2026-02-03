# fx-rate-service

The service refreshes FX rates asynchronously and stores results in PostgreSQL.

Supported currencies: EUR, USD, MXN

Rate provider: `exchangeratesapi.io`

> Note: Plan features vary by token. On the free tier, the API only supports `EUR` as the base currency.  
> Requests for pairs where the base is not EUR (e.g. `USD/EUR`) are accepted and processed, but will end with `status='error'`.


## Run (Docker Compose)

```bash
# set RATE_PROVIDER_ACCESS_TOKEN
export RATE_PROVIDER_ACCESS_TOKEN=
# check if RATE_PROVIDER_ACCESS_TOKEN is set
echo "$RATE_PROVIDER_ACCESS_TOKEN"

docker compose up --build -d
```

API: `http://localhost:8080` (migrations run automatically).

## API

Rate values are represented as `value_e6`, which is an integer that results from multiplying the float rate by 1,000,000.

- `POST /refresh-requests` body `{"pair":"USD/EUR"}` → refresh request id
- `GET /refresh-requests?id=<uuid>` → refresh request
- `GET /rates/latest?pair=USD/EUR` → latest stored rate

Try:

```bash
curl -s -X POST http://localhost:8080/refresh-requests \
  -H 'Content-Type: application/json' \
  -d '{"pair":"EUR/USD"}'

curl -s 'http://localhost:8080/refresh-requests?id=<uuid>'

curl -s 'http://localhost:8080/rates/latest?pair=EUR/USD'
```

## Local dev

```bash
cp .env.example .env
# set RATE_PROVIDER_ACCESS_TOKEN in .env file

# start postgres
docker compose up postgres -d

go run ./cmd/api
```

Default .env values:

- `HTTP_ADDR` (default `:8080`)
- `MIGRATIONS_DIR` (default `./migrations`)
- `RATE_PROVIDER_BASE_URL` (default `https://api.exchangeratesapi.io`)
- `REFRESH_WORKER_NUM` (default `4`)
- `REFRESH_WORKER_QUEUE_SIZE` (default `50`)
- `REFRESH_SWEEP_INTERVAL` (default `30s`)
- `REFRESH_STALE_AFTER` (default `5m`)

## Endpoint behavior

### `POST /refresh-requests`
Creates (or reuses) a refresh request for a currency pair and returns a request id. Processing happens asynchronously.

- **Idempotency:** repeated calls for the same `pair` while an existing request is in `processing` return the same request id.
- **Backpressure:** if the worker queue is full, the request is marked `error` with `service overloaded`.

### `GET /refresh-requests?id=<uuid>`
Returns the refresh request object (including its `status`).

`status` values:
- `processing` — still running
- `success` — completed; contains `value_e6`
- `error` — failed; contains `error_message`

### `GET /rates/latest?pair=XXX/YYY`
Returns the latest stored rate for the pair (if present).  
This does **not** trigger a refresh.

## How the async refresh pipeline works

### 1) Refresh request creation (HTTP)
`POST /refresh-requests` parses the pair (see [`models.ParseCurrencyPair`](internal/models/currency_pair.go)) and calls [`app.RateService.RefreshRate`](internal/app/service.go).

- Creates or reuses a row in `refresh_requests` via [`postgres.RefreshRateRepository.GetOrCreate`](internal/infra/postgres/refresh_rate_repository.go).
- If newly created, enqueues it into an in-memory buffered queue (`refreshChan`).
- If the queue is full, marks the request as `error` with `service overloaded` (see [`models.ErrServiceOverloaded`](internal/models/errors.go)).

Outcome:
- A `refresh_requests` row exists for the `pair`.

### 2) Background processing (worker pool)
On startup, the service launches a worker pool via [`app.RateService.StartWorkerPool`](internal/app/refresh_worker.go).

Each worker:
- pulls tasks from `refreshChan`
- fetches the rate from the provider (see [`exchange_rates_api.Client.FetchRate`](internal/infra/exchange_rates_api/client.go))
- converts it into `value_e6` (see [`models.NewValueE6FromFloat`](internal/models/value.go), i.e. `value_e6 = round(rate * 1e6)`)
- updates storage in one transaction via [`app.RateService.ProcessRefresh`](internal/app/refresh_processor.go) and [`postgres.TxManager.WithTx`](internal/infra/postgres/tx_manager.go):
    - upsert into `rates` (see [`postgres.RateRepository.Upsert`](internal/infra/postgres/rate_repository.go))
    - mark the request `success` and store `value_e6` (see [`postgres.RefreshRateRepository.Update`](internal/infra/postgres/refresh_rate_repository.go))

Outcome:
- `rates` contains the latest rate for the pair.
- `refresh_requests` is updated to `success` (or `error` if provider fetch/conversion/storage fails).

### 3) Stale request handling (sweeper)
A periodic sweeper marks stuck `processing` requests as `error` after `REFRESH_STALE_AFTER`.  
See [`app.StaleRefreshSweeper`](internal/app/refresh_sweeper.go) and [`postgres.RefreshRateRepository.MarkStaleProcessingRequest`](internal/infra/postgres/refresh_rate_repository.go).

## Data model
- `rates` stores the latest known rate per pair.
- `refresh_requests` stores refresh request records (the requested `pair`, current `status`, and the result fields `value_e6` / `error_message`, plus metadata like ids/timestamps).
  Schema: [migrations/0001_init.sql](migrations/0001_init.sql)
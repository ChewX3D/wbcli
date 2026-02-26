# WhiteBIT Integration Notes

Source references:

- Collateral limit order docs: https://docs.whitebit.com/api-reference/collateral-trading/collateral-limit-order
- Auth docs: https://docs.whitebit.com/private/http-auth
- Trade API index (contains bulk order endpoint): https://docs.whitebit.com/private/http-trade-v4

## Relevant Endpoints

- `POST /api/v4/collateral-limit-order`
- `POST /api/v4/collateral-bulk-limit-order` (for batch/range placement)

## Collateral Limit Order Request Fields

Required:

- `market` (example: `BTC_PERP`)
- `side` (`buy` or `sell`)
- `amount`
- `price`
- `expiration`

Optional:

- `clientOrderId`
- `postOnly`
- `ioc`
- `rpi`

Important constraint:

- WhiteBIT rejects requests where both `rpi=true` and `ioc=true`.

## Authentication and Signing

WhiteBIT private requests use headers:

- `X-TXC-APIKEY`: API key
- `X-TXC-PAYLOAD`: base64-encoded JSON payload
- `X-TXC-SIGNATURE`: hex HMAC-SHA512 of the base64 payload, signed with API secret
- `Content-Type: application/json`

Payload should include:

- `request`: endpoint path (for example `/api/v4/collateral-limit-order`)
- `nonce`: increasing numeric value
- `nonceWindow`: optional boolean for tolerant nonce validation
- request body fields (market, side, amount, etc.)

Implementation rule:

- generate payload JSON
- base64 encode payload
- sign encoded payload with HMAC-SHA512(secret)
- send payload fields as JSON body and include headers above

## Nonce Strategy

Recommended for CLI:

- use millisecond timestamp as default nonce
- ensure strict monotonic increase per profile/process
- keep a local last-nonce cache to avoid duplicate nonce on rapid calls

## API Error Handling

Normalize errors into categories:

- auth/signature/nonce errors
- validation errors (price/amount precision, invalid market)
- risk/business rule rejections
- temporary transport/server failures (retryable)

## Range Order Mapping

For `order range`:

1. generate deterministic order plan locally
2. validate each item (precision + notional)
3. map to bulk limit order payload
4. submit in chunks if needed
5. store per-order submission results for retry and reconciliation

# CLI Design (MVP)

## Goal

Provide a safe, scriptable CLI for WhiteBIT collateral trading with enough structure to power a future UI.

## Product Constraints

- do not expose API secrets in command history, logs, or repo files
- all order placement supports preview/validation before execution
- default behavior is conservative (explicit confirmation for multi-order submissions)
- command output can be rendered by both humans and machines (`table` and `json`)

## Command Model

### `whitbit keys`

- `whitbit keys set --profile default --api-key <...> --api-secret <...>`
- `whitbit keys list`
- `whitbit keys remove --profile default`
- `whitbit keys test --profile default`

Implementation notes:

- store secrets using platform secret storage (Keychain/libsecret/Credential Manager)
- keep only non-sensitive metadata in local config (profile name, created_at, last_used_at)

### `whitbit order place`

Example:

```bash
whitbit order place \
  --profile default \
  --market BTC_PERP \
  --side buy \
  --amount 0.01 \
  --price 50000 \
  --expiration 0 \
  --client-order-id my-order-001
```

Flow:

1. validate args and market precision rules
2. sign WhiteBIT request
3. submit order
4. print normalized response and local audit record ID

### `whitbit order range`

Example:

```bash
whitbit order range \
  --profile default \
  --market BTC_PERP \
  --side buy \
  --start-price 49000 \
  --end-price 50000 \
  --step 50 \
  --amount-mode constant \
  --base-amount 0.005 \
  --dry-run
```

### Range Amount Modes

- `constant`: `amount_i = base_amount`
- `arithmetic`: `amount_i = base_amount + i * delta_amount`
- `geometric`: `amount_i = base_amount * ratio^i` (optional, useful for aggressive laddering)

`i` is zero-based step index.

### Range Safety Controls

- hard cap for number of generated orders
- min/max notional checks before submission
- require `--confirm` for live batch placement unless interactive confirmation succeeds
- support `--dry-run` to preview generated plan and estimated exposure

## Output Contract (UI-ready)

Commands should return a stable structure internally:

- `request_id`
- `mode` (`single` or `range`)
- `orders_planned`
- `orders_submitted`
- `orders_failed`
- `errors[]`

CLI can render this as table/json; UI can consume same object via app layer.

## MVP Delivery Order

1. key storage + auth signing
2. single order placement
3. range planning (`dry-run` only)
4. range live submission with safeguards
5. cancellation helpers (`cancel by client-order-id prefix`, `cancel batch file`)

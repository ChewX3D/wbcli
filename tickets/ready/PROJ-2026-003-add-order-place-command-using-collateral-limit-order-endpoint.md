# PROJ-2026-003: Add collateral order place command using collateral limit order endpoint

ID: PROJ-2026-003
Title: Add collateral order place command using collateral limit order endpoint
Priority: P1
Status: Ready
Owner: chewbaccalol
Due Date: 2026-03-05
Created: 2026-02-26
Updated: 2026-03-01
Links: [CLI Design](../../docs/cli-design.md), [WhiteBIT Integration](../../docs/whitebit-integration.md)

Problem:
MVP requires placing a single collateral limit order safely from CLI with validation and stable output.

Outcome:
`wbcli collateral order place` submits one collateral limit order, validates inputs, and returns consistent table/json output.

Acceptance Criteria:
- [ ] Command supports required fields (`market`, `side`, `amount`, `price`) and optional flags.
- [ ] `--side` accepts aliases: `buy|long` and `sell|short`.
- [ ] Side alias normalization is performed in CLI command adapter files (`cmd/*`) before calling services; services receive canonical values only and do not normalize aliases.
- [ ] `collateral order place` always submits with `postOnly=true` (no CLI flag to disable).
- [ ] Command reads credentials from single-session auth state and uses signed client adapter.
- [ ] Command path is `wbcli collateral order place` (not `wbcli order place`).
- [ ] `--help` output for `wbcli collateral order place` is exhaustive and includes concrete `BTC-PERP` examples for common flows.
- [ ] Output contract includes `request_id`, `mode`, `orders_planned`, `orders_submitted`, `orders_failed`, `errors[]`.
- [ ] Unit tests cover successful placement and representative validation/auth failures; no integration tests for order submission endpoints.

Risks:
- Precision/notional validation gaps can cause rejected orders.
- Output contract instability will break future UI or automation wrappers.

Rollout Plan:
1. Implement command parser and validator.
2. Connect use-case to WhiteBIT adapter.
3. Add table/json renderer support and tests.

Alignment With Existing Code (Reuse-First):

1. Command surface and validation baseline already exist:
   - `cmd/order/place.go`
   - `cmd/order/validation.go`
   Use existing required flag parsing and shared validation helpers as the single source of truth for CLI-level checks.
2. Runtime dependency wiring already exists through one app factory:
   - `cmd/root.go`
   - `cmd/application_runtime.go`
   - `internal/app/application/factory.go`
   Extend `Application` with `Order` use-cases instead of creating per-command wiring.
3. Credential loading and single-session auth storage are already implemented:
   - `internal/adapters/secretstore/os_keychain.go`
   - `internal/adapters/configstore/profile_store.go`
   Reuse these adapters; do not add profile flags/backends for `collateral order place`.
4. Signed WhiteBIT transport is already implemented and should be reused directly:
   - `internal/adapters/whitebit/client.go`
   - `internal/adapters/whitebit/collateral.go`
   `PlaceCollateralLimitOrder` is already present; map use-case request into this client request.
5. Existing enum and request validation in transport client should remain source-of-truth for API contracts:
   - `OrderSide` / `PositionSide` enum validation in `internal/adapters/whitebit/collateral.go`
   Keep business logic outside transport client; only map/validate documented payload constraints there.
6. Existing auth command error-mapping style should be reused for actionable CLI errors:
   - `cmd/auth/errors.go`
   Implement equivalent order error mapping in command adapter layer, not in transport client.

Implementation Notes For This Ticket:

1. Keep current `collateral order place` flags for MVP required fields and `client-order-id`.
2. Normalize side aliases in CLI command adapter only (`cmd/*`) before service call:
   - `buy` and `long` map to transport `side=buy`
   - `sell` and `short` map to transport `side=sell`
   Service layer must treat incoming side as canonical and must not contain alias normalization logic.
3. Force `postOnly=true` in use-case/adapters for every submission.
4. Keep this ticket scope simple: do not add `ioc`/`rpi` flags and do not implement related conflict handling in this ticket.
5. Remove legacy `--profile` behavior from order commands to align with current single-session auth model.
6. Rename command path to `wbcli collateral order place`.
7. `--help` text must be detailed and explicit:
   - required/optional flags with meaning and constraints
   - side alias mapping (`buy|long`, `sell|short`)
   - note that `postOnly=true` is enforced
   - include multiple concrete `BTC-PERP` examples (buy, sell, with `--client-order-id`, and `--output json`)
8. Add `--output table|json` and render normalized contract fields:
   - `request_id`
   - `mode`
   - `orders_planned`
   - `orders_submitted`
   - `orders_failed`
   - `errors[]`
9. Keep order endpoint tests unit-only using mocks/fakes; no live integration tests for order submission endpoint.

Rollback Plan:
1. Disable live submission behind `--dry-run-only` temporary mode.
2. Revert endpoint mapping if response parsing errors appear in production use.

Status Notes:
- 2026-02-26: Created in Backlog.
- 2026-02-26: Sequenced for MVP after secure key storage and signed client are complete.
- 2026-02-26: Promoted to Ready with explicit dependency ordering and acceptance checks.
- 2026-02-28: Updated for single-session auth and safety policy (no order integration tests).
- 2026-03-01: Added reuse-first alignment with current architecture, clarified `postOnly=true` requirement, and documented concrete implementation mapping to existing code.
- 2026-03-01: Added `--side` alias requirement (`buy|long`, `sell|short`) and explicit normalization mapping for implementation.
- 2026-03-01: Updated command target to `wbcli collateral order place`, required CLI-layer side normalization only, and added mandatory exhaustive `--help` documentation with `BTC-PERP` examples.
- 2026-03-01: Simplified scope by removing `expiration` and removing `ioc/rpi` validation from this ticket.

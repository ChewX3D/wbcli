# PROJ-2026-003: Add order place command using collateral limit order endpoint

ID: PROJ-2026-003
Title: Add order place command using collateral limit order endpoint
Priority: P1
Status: Ready
Owner: chewbaccalol
Due Date: 2026-03-05
Created: 2026-02-26
Updated: 2026-02-26
Links: [CLI Design](../../docs/cli-design.md), [WhiteBIT Integration](../../docs/whitebit-integration.md)

Problem:
MVP requires placing a single collateral limit order safely from CLI with validation and stable output.

Outcome:
`wbcli order place` submits one collateral limit order, validates inputs, and returns consistent table/json output.

Acceptance Criteria:
- [ ] Command supports required fields (`market`, `side`, `amount`, `price`, `expiration`) and optional flags.
- [ ] Input validation blocks invalid combinations such as `rpi=true` with `ioc=true`.
- [ ] Command reads credentials from single-session auth state and uses signed client adapter.
- [ ] Output contract includes `request_id`, `mode`, `orders_planned`, `orders_submitted`, `orders_failed`, `errors[]`.
- [ ] Unit tests cover successful placement and representative validation/auth failures; no integration tests for order submission endpoints.

Risks:
- Precision/notional validation gaps can cause rejected orders.
- Output contract instability will break future UI or automation wrappers.

Rollout Plan:
1. Implement command parser and validator.
2. Connect use-case to WhiteBIT adapter.
3. Add table/json renderer support and tests.

Rollback Plan:
1. Disable live submission behind `--dry-run-only` temporary mode.
2. Revert endpoint mapping if response parsing errors appear in production use.

Status Notes:
- 2026-02-26: Created in Backlog.
- 2026-02-26: Sequenced for MVP after secure key storage and signed client are complete.
- 2026-02-26: Promoted to Ready with explicit dependency ordering and acceptance checks.
- 2026-02-28: Updated for single-session auth and safety policy (no order integration tests).

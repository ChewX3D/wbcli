# PROJ-2026-002: Implement WhiteBIT signed HTTP client

ID: PROJ-2026-002
Title: Implement WhiteBIT signed HTTP client
Priority: P1
Status: Done
Owner: chewbaccalol
Due Date: 2026-03-03
Created: 2026-02-26
Updated: 2026-02-28
Links: [WhiteBIT Integration](../../docs/whitebit-integration.md)

Problem:
Order endpoints require strict WhiteBIT payload signing and nonce behavior; ad-hoc HTTP calls are error-prone.

Outcome:
A reusable adapter signs private WhiteBIT requests correctly, handles nonce monotonicity, and normalizes API errors.

Acceptance Criteria:
- [x] Client sets `X-TXC-APIKEY`, `X-TXC-PAYLOAD`, and `X-TXC-SIGNATURE` per WhiteBIT requirements.
- [x] Nonce generator guarantees strictly increasing nonce values per process.
- [x] Error responses map into normalized categories (auth, validation, business-rule, transport).
- [x] Unit tests verify signature generation and payload wiring.
- [x] Integration-style tests with mock server validate headers/body for:
  - `/api/v4/collateral-account/hedge-mode`
- [x] Safety policy: order-submission endpoints are not covered by integration tests.

Risks:
- Nonce collisions can cause hard-to-debug intermittent failures.
- Endpoint contract drift may break signature behavior over time.

Rollout Plan:
1. Implement signer and nonce interfaces.
2. Build authenticated HTTP adapter with retry policy for retryable failures.
3. Add fixture-based and mock-server tests.

Rollback Plan:
1. Feature-flag new adapter and fall back to previous request path during stabilization.
2. Disable retries if they amplify rate-limit failures.

Status Notes:
- 2026-02-26: Created in Ready.
- 2026-02-26: Updated due date to reflect sequencing after project scaffold.
- 2026-02-28: Started implementation.
- 2026-02-28: Implemented shared signed WhiteBIT client with hedge-mode, collateral limit, and collateral bulk endpoints; auth login now reuses this client.
- 2026-02-28: Removed integration-style tests for order submission endpoints due safety requirement.

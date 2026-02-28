# PROJ-2026-017: Finalize auth test and capture macOS/Linux auth security evidence

ID: PROJ-2026-017
Title: Finalize auth test and capture macOS/Linux auth security evidence
Priority: P1
Status: Ready
Owner: chewbaccalol
Due Date: 2026-03-10
Created: 2026-02-28
Updated: 2026-02-28
Links: [PROJ-2026-001](../done/PROJ-2026-001-implement-secure-api-key-storage-adapter.md), [PROJ-2026-002](./PROJ-2026-002-implement-whitebit-signed-http-client.md), [PROJ-2026-016](../done/PROJ-2026-016-remove-auth-profiles-and-switch-to-single-session.md)

Problem:
`PROJ-2026-001` was closed for single-session auth storage scope, but two required security-delivery items remain:

- `auth test` runtime implementation is still deferred
- required manual macOS/Linux security verification evidence is not yet captured

Outcome:
`auth test` is implemented and verified using WhiteBIT probe wiring (after `PROJ-2026-002` readiness), and manual platform security evidence is attached for macOS and Linux.

Acceptance Criteria:
- [ ] `auth test` command is wired to `AuthProbe` implementation (no longer returns not-implemented).
- [ ] `auth test` reads credential from secure store and executes authenticated connectivity check.
- [ ] `auth test` error mapping distinguishes:
  - auth failure
  - transport/network failure
  - unknown/internal failure
- [ ] `auth test` output and errors do not leak API secret, payload, signature, or full API key.
- [ ] manual macOS verification evidence is captured for auth flow:
  - [x] login
  - [ ] status
  - [x] logout
  - [ ] keychain unavailable and permission-denied scenarios
  - [ ] metadata-only + `0600` config assertions
- [ ] manual Linux verification evidence is captured for auth flow:
  - [ ] login
  - [ ] status
  - [ ] logout
  - [ ] keychain unavailable and permission-denied scenarios
  - [ ] metadata-only + `0600` config assertions
- [ ] README includes short section pointing to where auth security verification evidence is tracked.
- [ ] `go test ./...` and `go build .` pass after `auth test` implementation.

Risks:
- WhiteBIT probe behavior may vary between environments and require stable fixture strategy.
- Linux keychain backend differences (DBus/secret service availability) can cause false negatives in manual runs.

Rollout Plan:
1. Wait for `PROJ-2026-002` readiness (signed WhiteBIT client + probe adapter prerequisites).
2. Implement `AuthProbe` adapter wiring in composition root.
3. Add/adjust tests for `auth test` success and failure-class mapping with no-leak assertions.
4. Execute and capture manual macOS security verification evidence.
5. Execute and capture manual Linux security verification evidence.
6. Update docs with evidence links/summary.
7. Re-run build/tests and finalize ticket.

Rollback Plan:
1. Revert `auth test` wiring commit if probe behavior is unstable.
2. Keep `auth test` returning explicit not-implemented message until probe is fixed.
3. Preserve existing login/status/logout behavior as-is.

Status Notes:
- 2026-02-28: Created in Ready.
- 2026-02-28: Scope extracted from PROJ-2026-001 canonical checklist to allow closing PROJ-2026-001 while tracking remaining auth test and platform evidence work.
- 2026-02-28: Manual verification update from user: `auth login` and `auth logout` confirmed working on macOS; Linux not tested yet.

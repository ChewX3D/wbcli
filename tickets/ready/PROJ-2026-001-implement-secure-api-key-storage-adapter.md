# PROJ-2026-001: Implement secure API key storage adapter

ID: PROJ-2026-001
Title: Implement secure API key storage adapter
Priority: P1
Status: Ready
Owner: chewbaccalol
Due Date: 2026-03-02
Created: 2026-02-26
Updated: 2026-02-26
Links: [CLI Design](../../docs/cli-design.md), [PROJ-2026-002](./PROJ-2026-002-implement-whitebit-signed-http-client.md), [PROJ-2026-006](./PROJ-2026-006-add-profile-config-store-for-non-secret-cli-metadata.md), [PROJ-2026-014](./PROJ-2026-014-define-credential-encryption-policy-and-encrypted-file-fallback-backend.md), [PROJ-2026-015](../backlog/PROJ-2026-015-implement-credential-access-controls-session-unlock-and-key-rotation-workflow.md)

Problem:
Trading commands require API credentials, but storing secrets in plaintext config or shell history is unsafe.

Outcome:
`wbcli auth set/list/remove/test` work end-to-end with profile isolation and secure-by-default behavior on `os-keychain`, with no secret leakage in files, logs, or command output.

Scope:
- implement `auth set/list/remove/test` against `os-keychain` backend
- store only non-secret profile metadata in local config
- enforce safe input and redaction rules in command handlers and tests

Out Of Scope:
- encrypted-file fallback backend implementation (handled by PROJ-2026-014)
- session unlock TTL, key rotation, and revoke workflows (handled by PROJ-2026-015)
- general WhiteBIT trading client beyond what `auth test` minimally requires (handled by PROJ-2026-002)

Dependencies:
- PROJ-2026-006 for profile metadata persistence model
- PROJ-2026-002 for authenticated connectivity verification behavior in `auth test`
- PROJ-2026-014 for explicit fallback backend policy (must not be silently auto-enabled here)

Acceptance Criteria:
- [ ] `auth set` input mode is secure by default:
  - [ ] `--api-key` is required.
  - [ ] API secret input is hidden prompt by default (non-echo).
  - [ ] optional non-interactive input path exists (`--api-secret-stdin`) for automation.
  - [ ] plaintext `--api-secret` flag is not used.
- [ ] `auth set` validation and storage behavior:
  - [ ] invalid/empty profile fails with clear error.
  - [ ] empty API key/secret fails with clear error.
  - [ ] credentials are written to `os-keychain` only.
  - [ ] if keychain is unavailable, command fails closed with actionable message (no silent insecure fallback).
- [ ] `auth list` behavior:
  - [ ] lists configured profiles and non-secret metadata only.
  - [ ] never prints API secret, payload, signature, or full API key.
- [ ] `auth remove` behavior:
  - [ ] removes credential record for a profile.
  - [ ] operation is idempotent (missing profile does not leak internals and is handled cleanly).
- [ ] `auth test` behavior:
  - [ ] reads credentials from secure store and performs authenticated connectivity check.
  - [ ] error mapping distinguishes auth, transport, and unknown failures.
  - [ ] output and logs never expose `X-TXC-PAYLOAD`, `X-TXC-SIGNATURE`, API secret, or full API key.
- [ ] Persistence boundaries:
  - [ ] no secret material is written to repo-tracked files or plain profile config.
  - [ ] profile config stores metadata only (profile name, timestamps, backend marker).
- [ ] Unit tests and command tests include success and negative paths for each command part.

Test Matrix:
- [ ] `auth set`: interactive secret input success, stdin secret input success, missing key, empty secret, keychain unavailable, permission denied.
- [ ] `auth list`: returns metadata-only rows, redaction assertions, empty state.
- [ ] `auth remove`: existing profile removal, missing profile idempotency, permission denied.
- [ ] `auth test`: valid credentials, auth failure, network timeout, backend unavailable, redaction in error output.
- [ ] cross-cutting: verify no sensitive value appears in stdout/stderr/log buffers under verbose and non-verbose modes.

Risks:
- Secret backend availability differs by OS and CI environment.
- Improper logging can still leak headers or partial secrets.

Rollout Plan:
1. Define credential domain types and secret store interface (no command logic yet).
2. Implement `os-keychain` adapter with deterministic error mapping.
3. Wire profile metadata integration from PROJ-2026-006.
4. Implement `auth set` flags/input path only (hidden prompt + stdin option) with tests.
5. Implement `auth set` validation + storage write path with tests.
6. Implement `auth list` metadata-only read path with redaction tests.
7. Implement `auth remove` delete path and idempotency behavior with tests.
8. Implement `auth test` minimal authenticated connectivity path with redaction and failure-classification tests.
9. Add command-level docs/help text updates for secure usage.
10. Run full verification and capture test evidence in ticket status note.

Verification Evidence (Required In Review):
- `go test ./...`
- `go build .`
- command-level tests for `auth set/list/remove/test` with redaction assertions
- explicit proof that no secret values are persisted in config files

Rollback Plan:
1. Disable write operations and keep read-only mode if backend instability appears.
2. Revert to previous credential path only if a secure alternative is confirmed.

Status Notes:
- 2026-02-26: Created in Ready.
- 2026-02-26: Refined with explicit secure I/O and testability criteria.
- 2026-02-26: Linked to encryption/access hardening tickets.
- 2026-02-26: Expanded into atomic command-part acceptance criteria, explicit scope boundaries, dependency mapping, and detailed rollout/test matrix.

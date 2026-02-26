# PROJ-2026-001: Implement secure API key storage adapter

ID: PROJ-2026-001
Title: Implement secure API key storage adapter
Priority: P1
Status: Ready
Owner: nocle
Due Date: 2026-03-02
Created: 2026-02-26
Updated: 2026-02-26
Links: [CLI Design](../../docs/cli-design.md), [PROJ-2026-002](./PROJ-2026-002-implement-whitebit-signed-http-client.md), [PROJ-2026-006](./PROJ-2026-006-add-profile-config-store-for-non-secret-cli-metadata.md), [PROJ-2026-014](./PROJ-2026-014-define-credential-encryption-policy-and-encrypted-file-fallback-backend.md), [PROJ-2026-015](../backlog/PROJ-2026-015-implement-credential-access-controls-session-unlock-and-key-rotation-workflow.md)

Problem:
Trading commands require API credentials, but storing secrets in plaintext config or shell history is unsafe.

Outcome:
`wbcli auth login/use/profiles list/logout/current/test` work end-to-end with profile isolation and secure-by-default behavior on `os-keychain`, with no secret leakage in files, logs, or command output.

Scope:
- implement `auth login/use/profiles list/logout/current/test` against `os-keychain` backend
- store only non-secret profile metadata in local config
- enforce safe input and redaction rules in command handlers and tests

Command Model Examples (Approved):
- `wbcli auth profiles list`
- `wbcli auth use <profile>`
- `wbcli auth login --profile <profile>`
- `wbcli auth logout --profile <profile>`
- `wbcli auth current`

Auth Login Architecture (Hexagonal, Mandatory):
- Inbound adapter (`cmd/auth_login.go`):
  - parses flags and input mode (`--api-key`, `--profile`, `--api-secret-stdin`)
  - reads secret from hidden prompt or stdin (never from plaintext flag)
  - performs syntax-level validation only (required flags, non-empty values)
  - calls app use-case port and renders safe output
- Application auth service (`internal/app/services/auth/login.go`):
  - defines `AuthLoginService` with `Execute(ctx, req)` boundary
  - orchestrates domain validation + outbound ports
  - maps lower-level adapter errors into stable app errors
  - returns transport-agnostic result DTO
  - must be implementable without WhiteBIT client readiness
- Domain (`internal/domain/auth`):
  - credential value object and invariants (non-empty profile/key/secret)
  - profile naming policy and normalization rules
  - no CLI, keychain, filesystem, or network dependencies
- Outbound ports (`internal/app/ports/auth.go`):
  - `CredentialStore` port: `Save(ctx, profile, credential)`
  - `ProfileStore` port: `UpsertProfile(ctx, meta)` + `SetActiveProfile(ctx, profile)` (if login updates active profile)
  - `Clock` port for deterministic timestamps in metadata
  - `AuthProbe` port for `auth test` only (implemented later, after WhiteBIT client is ready)
- Outbound adapters:
  - `internal/adapters/secretstore`: os-keychain implementation of `CredentialStore`
  - profile metadata adapter from PROJ-2026-006 implementing `ProfileStore`
- Composition root (`main.go` + `cmd`):
  - instantiate adapters, wire use-case, inject into `auth login` command
  - no business logic in wiring layer

Auth Login Request/Response Contract:
- Request:
  - `profile` (default `default`)
  - `api_key` (required)
  - `api_secret` (required, ephemeral in memory)
  - `input_mode` (`prompt` | `stdin`) for auditing/tests
- Response (safe):
  - `profile`
  - `backend` (for example `os-keychain`)
  - `saved_at`
  - optional `warnings[]`
- Never include secret, payload, signature, or full API key in response.

Auth Login Error Contract:
- `ERR_PROFILE_INVALID`
- `ERR_API_KEY_REQUIRED`
- `ERR_API_SECRET_REQUIRED`
- `ERR_SECRETSTORE_UNAVAILABLE`
- `ERR_SECRETSTORE_PERMISSION_DENIED`
- `ERR_PROFILESTORE_WRITE_FAILED`
- `ERR_INTERNAL`
- CLI output must map these to actionable user messages while preserving redaction.

Auth Login Security Contract:
- secret input uses non-echo prompt by default
- `--api-secret-stdin` is allowed for automation only
- no plaintext `--api-secret` flag
- secret is held only for request lifetime and cleared where practical after use
- logs/output redact sensitive fields in normal and verbose modes
- fail closed when `os-keychain` is unavailable for this ticket scope

Atomic Commit Slices For `auth login` (Required):
1. add command skeleton + flags (`login` only, no storage call)
2. add secret input reader (prompt + stdin) + tests
3. add request validation + domain value object tests
4. add app use-case and port interfaces + unit tests with fakes
5. add os-keychain adapter implementation + adapter tests
6. wire command to use-case + safe success/error rendering tests
7. metadata update integration (`ProfileStore`) + tests
8. docs/help update for secure usage

Out Of Scope:
- encrypted-file fallback backend implementation (handled by PROJ-2026-014)
- session unlock TTL, key rotation, and revoke workflows (handled by PROJ-2026-015)
- general WhiteBIT trading client beyond what `auth test` minimally requires (handled by PROJ-2026-002)

Dependencies:
- PROJ-2026-006 for profile metadata persistence model
- PROJ-2026-002 for `auth test` authenticated connectivity verification behavior (final phase only)
- PROJ-2026-014 for explicit fallback backend policy (must not be silently auto-enabled here)

Acceptance Criteria:
- [ ] `auth login` input mode is secure by default:
  - [ ] `--api-key` is required.
  - [ ] API secret input is hidden prompt by default (non-echo).
  - [ ] optional non-interactive input path exists (`--api-secret-stdin`) for automation.
  - [ ] plaintext `--api-secret` flag is not used.
  - [ ] plaintext secret flags are prohibited for all auth commands (including legacy command paths).
  - [ ] prompt mode is used only when input is interactive (TTY); non-interactive mode must fail with a clear hint to use `--api-secret-stdin`.
  - [ ] `--api-secret-stdin` reads secret only from stdin, does not echo it, and fails with clear error on empty input.
  - [ ] prompt mode and stdin mode are mutually exclusive and produce clear validation errors if misused.
- [ ] `auth login` validation and storage behavior:
  - [ ] invalid/empty profile fails with clear error.
  - [ ] empty API key/secret fails with clear error.
  - [ ] credentials are written to `os-keychain` only.
  - [ ] if keychain is unavailable, command fails closed with actionable message (no silent insecure fallback).
- [ ] Platform support requirements:
  - [ ] must work on macOS.
  - [ ] must work on Linux.
  - [ ] Windows support is optional for this ticket (best-effort only).
- [ ] `auth use` behavior:
  - [ ] selects active profile from existing profile set.
  - [ ] fails clearly if profile is missing or has no stored credentials.
  - [ ] updates only non-secret active-profile metadata in config.
- [ ] `auth profiles list` behavior:
  - [ ] lists configured profiles and non-secret metadata only.
  - [ ] never prints API secret, payload, signature, or full API key.
- [ ] `auth logout` behavior:
  - [ ] removes credential record for a profile.
  - [ ] operation is idempotent (missing profile does not leak internals and is handled cleanly).
- [ ] `auth current` behavior:
  - [ ] prints only current active profile and safe metadata.
  - [ ] never prints secret material.
- [ ] `auth test` behavior:
  - [ ] reads credentials from secure store and performs authenticated connectivity check.
  - [ ] error mapping distinguishes auth, transport, and unknown failures.
  - [ ] output and logs never expose `X-TXC-PAYLOAD`, `X-TXC-SIGNATURE`, API secret, or full API key.
  - [ ] implemented last, after PROJ-2026-002 is ready.
- [ ] Persistence boundaries:
  - [ ] no secret material is written to repo-tracked files or plain profile config.
  - [ ] profile config stores metadata only (profile name, timestamps, backend marker).
- [ ] Unit tests and command tests include success and negative paths for each command part.
- [ ] Documentation requirements for implementation:
  - [ ] README must explain auth secret input modes in simple, easy-to-understand language.
  - [ ] README must include examples for local interactive use (prompt) and automation/CI use (`--api-secret-stdin`).
  - [ ] README must explicitly warn not to pass secrets via command arguments.
  - [ ] Cobra command help must include practical `Example` blocks for `auth login/use/profiles list/logout/current/test`.

Test Matrix:
- [ ] `auth login`: interactive secret input success, stdin secret input success, missing key, empty secret, keychain unavailable, permission denied.
- [ ] platform compatibility:
  - [ ] run `auth login/use/profiles list/logout/current` verification on macOS.
  - [ ] run `auth login/use/profiles list/logout/current` verification on Linux.
  - [ ] Windows verification is optional and recorded only if performed.
- [ ] `auth use`: existing profile selection, missing profile failure, active-profile metadata update.
- [ ] `auth profiles list`: returns metadata-only rows, redaction assertions, empty state.
- [ ] `auth logout`: existing profile removal, missing profile idempotency, permission denied.
- [ ] `auth current`: correct active profile output, empty state behavior, redaction assertions.
- [ ] `auth test`: valid credentials, auth failure, network timeout, backend unavailable, redaction in error output.
- [ ] cross-cutting: verify no sensitive value appears in stdout/stderr/log buffers under verbose and non-verbose modes.

Risks:
- Secret backend availability differs by OS and CI environment.
- Improper logging can still leak headers or partial secrets.

Rollout Plan:
1. Define credential domain types and auth service interfaces (`CredentialStore`, `ProfileStore`, `Clock`) with no WhiteBIT dependency.
2. Implement `os-keychain` adapter with deterministic error mapping.
3. Wire profile metadata integration from PROJ-2026-006.
4. Implement `auth login` flags/input path only (hidden prompt + stdin option) with tests.
5. Implement `auth login` validation + storage write path via `AuthLoginService` with tests.
6. Implement `auth profiles list` metadata-only read path with redaction tests.
7. Implement `auth use` active-profile selection path with metadata-only update tests.
8. Implement `auth logout` delete path and idempotency behavior with tests.
9. Implement `auth current` metadata read path with safe output tests.
10. Add command-level docs/help text updates for secure usage.
10.1. Update README with clear simple-language explanation of secret input behavior (prompt vs stdin), with safe examples and warning against command-argument secrets.
10.2. Add/verify `cobra.Command.Example` text for all auth commands with safe usage examples.
11. Run verification for `login/use/profiles list/logout/current` and capture evidence.
12. After PROJ-2026-002 is ready, implement `auth test` command using `AuthProbe` port, then add redaction/failure-classification tests as the final step.
13. Run full verification including `auth test` and capture final evidence.

Verification Evidence (Required In Review):
- `go test ./...`
- `go build .`
- command-level tests for `auth login/use/profiles list/logout/current/test` with redaction assertions
- explicit proof that no secret values are persisted in config files
- README excerpt/evidence showing simple-language auth input guidance and safe usage examples
- CLI help evidence showing `Example` blocks for auth commands
- platform evidence:
  - macOS pass evidence is required
  - Linux pass evidence is required
  - Windows evidence is optional

Rollback Plan:
1. Disable write operations and keep read-only mode if backend instability appears.
2. Revert to previous credential path only if a secure alternative is confirmed.

Status Notes:
- 2026-02-26: Created in Ready.
- 2026-02-26: Refined with explicit secure I/O and testability criteria.
- 2026-02-26: Linked to encryption/access hardening tickets.
- 2026-02-26: Expanded into atomic command-part acceptance criteria, explicit scope boundaries, dependency mapping, and detailed rollout/test matrix.
- 2026-02-26: Updated command model examples (`auth login/use/profiles list/logout/current`) and reordered rollout to start from `auth login`.
- 2026-02-26: Added mandatory hexagonal architecture contract for `auth login` (layers, ports, DTOs, error codes, security contract, and atomic commit slices).
- 2026-02-26: Clarified `auth login` as independent app service and moved `auth test` to final gated phase after WhiteBIT client readiness.
- 2026-02-26: Set platform acceptance rule: macOS and Linux required, Windows optional.
- 2026-02-26: Added requirement to include practical usage examples directly in Cobra command help (`Example` field).

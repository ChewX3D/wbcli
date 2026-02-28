# PROJ-2026-001: Implement secure API key storage adapter

ID: PROJ-2026-001
Title: Implement secure API key storage adapter
Priority: P1
Status: Ready
Owner: nocle
Due Date: 2026-03-02
Created: 2026-02-26
Updated: 2026-02-26
Links: [CLI Design](../../docs/cli-design.md), [PROJ-2026-002](./PROJ-2026-002-implement-whitebit-signed-http-client.md), [PROJ-2026-006](./PROJ-2026-006-add-profile-config-store-for-non-secret-cli-metadata.md), [PROJ-2026-014](./PROJ-2026-014-define-credential-encryption-policy-and-encrypted-file-fallback-backend.md), [PROJ-2026-015](../backlog/PROJ-2026-015-implement-credential-access-controls-session-unlock-and-key-rotation-workflow.md), [PROJ-2026-016](./PROJ-2026-016-remove-auth-profiles-and-switch-to-single-session.md)

Problem:
Trading commands require API credentials, but storing secrets in plaintext config or shell history is unsafe.

Outcome:
`wbcli auth login/logout/status/test` work end-to-end in single-session mode with secure-by-default behavior on `os-keychain`, with no secret leakage in files, logs, or command output.

Scope:
- implement `auth login/logout/status/test` against `os-keychain` backend
- store only non-secret session metadata in local config
- enforce safe input and redaction rules in command handlers and tests

Command Model Examples (Approved):
- `wbcli auth login`
- `wbcli auth logout`
- `wbcli auth status`
- `wbcli auth test`
- `printf '%s\n%s\n' "$WBCLI_API_KEY" "$WBCLI_API_SECRET" | wbcli auth login`
- legacy `wbcli auth set` must be removed (no compatibility shim required for this project)

Auth Login Architecture (Hexagonal, Mandatory):
- Inbound adapter (`cmd/auth_login.go`):
  - parses flags and input contract (`--profile`, optional overwrite controls)
  - reads credentials from stdin only (`api_key` first line, `api_secret` second line)
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
  - `profile` (required, explicit; no implicit default profile)
  - `api_key` (required, from stdin first line)
  - `api_secret` (required, from stdin second line; ephemeral in memory)
  - `input_mode` fixed to `stdin` for auditing/tests
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
- credential input for `auth login` is stdin-only (local and CI)
- no `--api-key` flag
- no plaintext `--api-secret` flag
- prompt-based secret input is not supported
- secret is held only for request lifetime and cleared where practical after use
- logs/output redact sensitive fields in normal and verbose modes
- fail closed when `os-keychain` is unavailable for this ticket scope

Atomic Commit Slices For `auth login` (Required):
1. add command skeleton + flags (`login` only, no storage call)
2. add stdin credential reader (`api_key` + `api_secret`) + tests
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

Scope Update (2026-02-28):
- profile-oriented command model and checklist items in this ticket are superseded by PROJ-2026-016 (single-session auth model)
- completion of this ticket must be evaluated together with PROJ-2026-016 acceptance criteria and remaining security/test coverage work

Reconciled Checklist (Canonical, 2026-02-28):
- [x] auth command model is single-session (`login/logout/status/test`) with no profile command/flag dependencies.
- [x] `auth login` uses stdin-only two-line credential input (`api_key`, `api_secret`) and rejects invalid payloads.
- [x] `auth login` overwrites existing session credentials by default (no `--force` flow).
- [x] credentials are written to `os-keychain` only; unavailable and permission-denied paths return actionable errors.
- [x] config remains metadata-only at `~/.wbcli/config.yaml` and uses `0600` permissions on macOS/Linux.
- [x] legacy `auth set/use/list/current` are removed from command tree/help output.
- [ ] `auth test` remains deferred until PROJ-2026-002 is implemented.
- [ ] required macOS/Linux manual security evidence is still pending capture in ticket review evidence.

Legacy Checklist Note:
- acceptance criteria and test matrix blocks below were authored for profile-based auth
- they are retained as historical context only and are superseded by the canonical reconciled checklist above plus PROJ-2026-016

Acceptance Criteria:
- [x] `auth login` input mode is secure by default:
  - [x] `--profile` is required and must be provided explicitly.
  - [x] implicit default profile selection is prohibited.
  - [x] API key and API secret are accepted only through stdin payload.
  - [x] no `--api-key` flag is supported.
  - [x] plaintext `--api-secret` flag is not used.
  - [x] prompt-based secret input is not supported.
  - [x] plaintext secret flags are prohibited for all auth commands (including legacy command paths).
  - [x] `auth login` fails with clear error when stdin payload is missing or empty.
  - [x] legacy `auth set` command is removed from CLI command tree and help output.
  - [x] stdin credential parsing contract is explicit and deterministic:
    - [x] read stdin input once with bounded maximum size.
    - [x] accept exactly two non-empty logical lines: first `api_key`, second `api_secret`.
    - [x] trim exactly one trailing line ending (`\\n` or `\\r\\n`) for common shell piping compatibility.
    - [x] reject empty effective key or secret value after parsing.
    - [x] reject missing second line and reject extra lines.
    - [x] parsing failures must not echo or log raw stdin content.
- [ ] `auth login` validation and storage behavior:
  - [x] invalid/empty profile fails with clear error.
  - [x] profile format is validated in `auth login` using charset `[a-zA-Z0-9._-]` and length `1..64`.
  - [x] profile values with leading/trailing spaces are rejected.
  - [x] no profile normalization is allowed; invalid profile input fails hard with explicit error.
  - [ ] profile format validation scope for this ticket is `auth login` only (profile creation path).
  - [x] empty API key/secret fails with clear error.
  - [x] if credentials already exist for the profile, overwrite is blocked unless explicitly confirmed.
  - [x] non-interactive overwrite requires explicit `--force`; otherwise fail with actionable error.
  - [ ] interactive overwrite requires explicit confirmation prompt.
  - [x] overwrite flow must never print old/new secret values.
  - [x] credentials are written to `os-keychain` only.
  - [x] if keychain is unavailable, command fails closed with actionable message (no silent insecure fallback).
  - [x] fallback to insecure storage (plaintext file/env/arg) is strictly prohibited.
  - [x] any non-keychain fallback must be explicitly configured and is out of scope for this ticket.
- [ ] Platform support requirements:
  - [ ] must work on macOS.
  - [ ] must work on Linux.
  - [ ] Windows support is optional for this ticket (best-effort only).
  - [ ] security verification must run on both macOS and Linux for auth commands (`login/use/list/logout/current`; `test` when implemented).
  - [ ] platform checks must include keychain unavailable/permission-denied scenarios per OS.
  - [ ] platform checks must include redaction and config-boundary assertions per OS.
- [x] `auth use` behavior:
  - [x] selects active profile from existing profile set.
  - [x] fails clearly if profile is missing or has no stored credentials.
  - [x] updates only non-secret active-profile metadata in config.
- [x] `auth list` behavior:
  - [x] lists configured profiles and non-secret metadata only.
  - [x] never prints API secret, payload, signature, or full API key.
- [x] `auth logout` behavior:
  - [x] removes credential record for a profile.
  - [x] operation is idempotent (missing profile does not leak internals and is handled cleanly).
- [x] `auth current` behavior:
  - [x] prints only current active profile and safe metadata.
  - [x] never prints secret material.
- [ ] `auth test` behavior:
  - [ ] reads credentials from secure store and performs authenticated connectivity check.
  - [ ] error mapping distinguishes auth, transport, and unknown failures.
  - [ ] output and logs never expose `X-TXC-PAYLOAD`, `X-TXC-SIGNATURE`, API secret, or full API key.
  - [ ] implemented last, after PROJ-2026-002 is ready.
- [ ] Redaction policy requirements:
  - [ ] API secret, payload, signature, and full API key are never printed to stdout/stderr/logs.
  - [ ] redaction policy is enforced in both normal and verbose modes.
  - [ ] all auth commands use shared redaction helpers (no per-command ad-hoc masking logic).
  - [ ] failures and debug output preserve diagnostic value without leaking sensitive material.
- [ ] Secret memory-lifetime requirements:
  - [x] secret is held in memory only for minimal required path (read -> validate -> store/sign -> clear).
  - [x] use `[]byte` for sensitive secret handling where practical.
  - [x] clear sensitive buffers after use (best-effort wipe) and avoid unnecessary copies.
  - [ ] avoid propagating secret values into long-lived structs, error objects, or log contexts.
- [x] Persistence boundaries:
  - [x] no secret material is written to repo-tracked files or plain profile config.
  - [x] profile config stores metadata only (profile name, timestamps, backend marker, active profile).
  - [x] on macOS/Linux, config path is `~/.wbcli/config.yaml`.
  - [x] config file at `~/.wbcli/config.yaml` is created/updated with owner-only permissions (`0600`).
- [ ] Unit tests and command tests include success and negative paths for each command part.
- [ ] Documentation requirements for implementation:
  - [x] README must explain auth secret input modes in simple, easy-to-understand language.
  - [x] README must include stdin-only login examples for local shell usage and automation/CI usage.
  - [x] README must explicitly warn not to pass credentials via command arguments.
  - [x] Cobra command help must include practical `Example` blocks for `auth login/use/list/logout/current/test`.
  - [ ] README must include operational key-hygiene guidance:
    - [ ] one API key per profile/environment.
    - [ ] least-privilege key scopes for WhiteBIT permissions.
    - [ ] exchange-side IP allowlist recommendation where available.
    - [ ] rotation cadence and emergency revoke flow (`auth logout` + exchange-side revoke checklist).
  - [ ] Cobra auth command help must include concise security-hygiene notes for safe operations.

Test Matrix:
- [ ] `auth login`: stdin credential payload success, missing profile, invalid profile format (spaces/special chars/unicode/too long), missing key, empty secret, keychain unavailable, permission denied.
- [ ] stdin parsing behavior:
  - [x] valid two-line stdin payload (`api_key` + `api_secret`) succeeds.
  - [x] single trailing newline is handled per contract.
  - [x] empty stdin fails with clear error.
  - [x] missing second line fails with clear error.
  - [x] extra lines after second credential line fail with clear error.
  - [x] oversized stdin input fails with clear error.
  - [ ] stdin parse errors do not leak raw input to stdout/stderr/logs.
- [x] command migration behavior:
  - [x] `wbcli auth set ...` is unavailable after migration and returns unknown-command error.
  - [x] `wbcli auth login ...` and `wbcli auth use ...` are available and shown in help output.
- [ ] overwrite control behavior:
  - [ ] new profile create path succeeds without overwrite confirmation.
  - [ ] existing profile update without confirmation/`--force` is rejected.
  - [ ] interactive confirmed overwrite succeeds.
  - [ ] non-interactive overwrite with `--force` succeeds.
  - [ ] overwrite path does not leak old/new secret values.
- [ ] fail-closed storage behavior:
  - [ ] when keychain is unavailable, command returns actionable error and exits non-zero.
  - [ ] verify there is no fallback write to plaintext config, env-based cache, or other implicit storage.
- [ ] config boundary behavior:
  - [x] verify `~/.wbcli/config.yaml` contains metadata only and no secret values.
  - [x] verify `~/.wbcli/config.yaml` permissions are `0600` on macOS/Linux.
- [ ] platform compatibility:
  - [ ] run `auth login/use/list/logout/current` verification on macOS.
  - [ ] run `auth login/use/list/logout/current` verification on Linux.
  - [ ] Windows verification is optional and recorded only if performed.
  - [ ] include keychain unavailable and permission-denied cases for both macOS and Linux.
  - [ ] include redaction no-leak assertions on both macOS and Linux.
  - [ ] include `~/.wbcli/config.yaml` metadata-only and `0600` assertions on both macOS and Linux.
- [ ] `auth use`: existing profile selection, missing profile failure, active-profile metadata update.
- [ ] `auth list`: returns metadata-only rows, redaction assertions, empty state.
- [ ] `auth logout`: existing profile removal, missing profile idempotency, permission denied.
- [ ] `auth current`: correct active profile output, empty state behavior, redaction assertions.
- [ ] `auth test`: valid credentials, auth failure, network timeout, backend unavailable, redaction in error output.
- [ ] cross-cutting: verify no sensitive value appears in stdout/stderr/log buffers under verbose and non-verbose modes.
- [ ] redaction contract:
  - [ ] assert shared redaction helper is used by all auth command output paths.
  - [ ] assert sensitive values are absent from command output and logs in success and error scenarios.
- [ ] secret memory contract:
  - [x] verify secret buffers are cleared after use on primary paths.
  - [ ] verify errors and debug paths do not retain or expose secret values.

Risks:
- Secret backend availability differs by OS and CI environment.
- Improper logging can still leak headers or partial secrets.

Rollout Plan:
1. Define credential domain types and auth service interfaces (`CredentialStore`, `ProfileStore`, `Clock`) with no WhiteBIT dependency.
2. Implement `os-keychain` adapter with deterministic error mapping.
2.1. Add fail-closed guardrails to prohibit implicit fallback when keychain is unavailable.
3. Wire profile metadata integration from PROJ-2026-006.
3.1. Enforce metadata config path (`~/.wbcli/config.yaml`) and owner-only file mode (`0600`) on macOS/Linux.
4. Implement `auth login` stdin-only credential input path (`api_key` + `api_secret`) with tests.
4.1. Remove legacy `auth set` command wiring and ensure `auth login` + `auth use` are the active entrypoints.
4.2. Implement strict stdin credential parser contract (bounded read, two-line contract, newline handling, no-leak errors).
5. Implement `auth login` validation + storage write path via `AuthLoginService` with tests.
5.1. Add secret buffer lifecycle handling (minimal lifetime + best-effort wipe + no unnecessary copies).
5.2. Add credential overwrite guardrails (interactive confirm and non-interactive `--force` behavior) with leak-safe output handling.
6. Implement `auth list` metadata-only read path with redaction tests.
7. Implement `auth use` active-profile selection path with metadata-only update tests.
8. Implement `auth logout` delete path and idempotency behavior with tests.
9. Implement `auth current` metadata read path with safe output tests.
9.1. Implement shared auth redaction helper and wire all auth command outputs/logging through it.
10. Add command-level docs/help text updates for secure usage.
10.1. Update README with clear simple-language explanation of stdin-only credential input behavior, with safe examples and warning against command-argument credentials.
10.2. Add/verify `cobra.Command.Example` text for all auth commands with safe usage examples.
10.3. Add operational key-hygiene section to README and concise security-hygiene notes in Cobra auth help text.
11. Run verification for `login/use/list/logout/current` and capture evidence.
11.1. Align auth security verification with existing CI test/build pipelines and extend CI where needed (do not create a separate disconnected verification flow).
11.2. Ensure CI artifacts/summaries capture platform-specific security evidence for macOS and Linux.
12. After PROJ-2026-002 is ready, implement `auth test` command using `AuthProbe` port, then add redaction/failure-classification tests as the final step.
13. Run full verification including `auth test` and capture final evidence.
14. Add integration tests for auth flows using mock secret-store/keychain adapters (do not use real OS keychain in tests/CI).

Verification Evidence (Required In Review):
- `go test ./...`
- `go build .`
- command-level tests for `auth login/use/list/logout/current/test` with redaction assertions
- explicit proof that no secret values are persisted in config files
- explicit proof that `~/.wbcli/config.yaml` contains metadata only and no secrets
- explicit proof that `~/.wbcli/config.yaml` uses `0600` permissions on macOS/Linux
- explicit proof that keychain-unavailable path fails closed with no implicit fallback write
- explicit proof that redaction contract is enforced in normal and verbose modes
- explicit proof that secret memory-lifetime controls are implemented (minimal lifetime, cleared buffers, no secret in errors/logs)
- explicit proof that stdin parser contract is enforced (two-line key/secret contract, newline handling, empty/missing/extra-line/oversize rejection, no-leak errors)
- README excerpt/evidence showing simple-language auth input guidance and safe usage examples
- CLI help evidence showing `Example` blocks for auth commands
- README excerpt/evidence showing operational key-hygiene guidance (per-profile keys, least privilege, allowlist, rotation/revoke)
- CLI help evidence showing concise auth security-hygiene notes
- platform evidence:
  - macOS pass evidence is required
  - Linux pass evidence is required
  - Windows evidence is optional
- CI alignment evidence:
  - show how auth security checks integrate with existing CI test/build workflows
  - include artifact/log links for macOS and Linux security verification results

Rollback Plan:
1. Disable write operations and keep read-only mode if backend instability appears.
2. Revert to previous credential path only if a secure alternative is confirmed.

Status Notes:
- 2026-02-26: Created in Ready.
- 2026-02-26: Refined with explicit secure I/O and testability criteria.
- 2026-02-26: Linked to encryption/access hardening tickets.
- 2026-02-26: Expanded into atomic command-part acceptance criteria, explicit scope boundaries, dependency mapping, and detailed rollout/test matrix.
- 2026-02-26: Updated command model examples (`auth login/use/list/logout/current`) and reordered rollout to start from `auth login`.
- 2026-02-26: Added mandatory hexagonal architecture contract for `auth login` (layers, ports, DTOs, error codes, security contract, and atomic commit slices).
- 2026-02-26: Clarified `auth login` as independent app service and moved `auth test` to final gated phase after WhiteBIT client readiness.
- 2026-02-26: Set platform acceptance rule: macOS and Linux required, Windows optional.
- 2026-02-26: Added explicit redaction security point (shared helper, verbose-safe behavior, and leak-prevention evidence).
- 2026-02-26: Added requirement to include practical usage examples directly in Cobra command help (`Example` field).
- 2026-02-26: Added explicit-profile rule (no implicit default profile; user must set profile explicitly).
- 2026-02-26: Added security boundary for metadata config path `~/.wbcli/config.yaml` with required `0600` permissions and no secret persistence.
- 2026-02-26: Added strict profile validation rule for `auth login` only (no normalization; hard errors for invalid values).
- 2026-02-26: Added secret memory-lifetime security point (minimal in-memory lifetime, best-effort buffer wipe, and no secret propagation in errors/logs).
- 2026-02-26: Added cross-platform security verification requirement (macOS/Linux) and mandated alignment with existing CI test/build workflows.
- 2026-02-26: Added overwrite security control for `auth login` (explicit confirmation/`--force`, no silent overwrite, no secret leakage on update path).
- 2026-02-26: Added operational key-hygiene documentation requirement (per-profile keys, least privilege, allowlist, rotation/revoke) for README and Cobra help.
- 2026-02-26: Confirmed migration strategy for this project: remove legacy `auth set` and use `auth login` + `auth use` only.
- 2026-02-26: Changed `auth login` to stdin-only credential input (removed `--api-key`, removed prompt input path, and replaced `--api-secret-stdin` contract with a strict two-line stdin contract).
- 2026-02-28: Scope updated by PROJ-2026-016 to remove profiles and move auth to single-session mode; profile-based checklist items are superseded.

Final Note (Mandatory):
- add integration tests with mock secret-store/keychain adapters; do not run integration tests against real OS keychains.

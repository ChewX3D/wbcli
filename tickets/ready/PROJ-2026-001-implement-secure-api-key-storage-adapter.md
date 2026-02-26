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
  - `profile` (required, explicit; no implicit default profile)
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
  - [ ] `--profile` is required and must be provided explicitly.
  - [ ] implicit default profile selection is prohibited.
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
  - [ ] profile format is validated in `auth login` using charset `[a-zA-Z0-9._-]` and length `1..64`.
  - [ ] profile values with leading/trailing spaces are rejected.
  - [ ] no profile normalization is allowed; invalid profile input fails hard with explicit error.
  - [ ] profile format validation scope for this ticket is `auth login` only (profile creation path).
  - [ ] empty API key/secret fails with clear error.
  - [ ] if credentials already exist for the profile, overwrite is blocked unless explicitly confirmed.
  - [ ] non-interactive overwrite requires explicit `--force`; otherwise fail with actionable error.
  - [ ] interactive overwrite requires explicit confirmation prompt.
  - [ ] overwrite flow must never print old/new secret values.
  - [ ] credentials are written to `os-keychain` only.
  - [ ] if keychain is unavailable, command fails closed with actionable message (no silent insecure fallback).
  - [ ] fallback to insecure storage (plaintext file/env/arg) is strictly prohibited.
  - [ ] any non-keychain fallback must be explicitly configured and is out of scope for this ticket.
- [ ] Platform support requirements:
  - [ ] must work on macOS.
  - [ ] must work on Linux.
  - [ ] Windows support is optional for this ticket (best-effort only).
  - [ ] security verification must run on both macOS and Linux for auth commands (`login/use/profiles list/logout/current`; `test` when implemented).
  - [ ] platform checks must include keychain unavailable/permission-denied scenarios per OS.
  - [ ] platform checks must include redaction and config-boundary assertions per OS.
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
- [ ] Redaction policy requirements:
  - [ ] API secret, payload, signature, and full API key are never printed to stdout/stderr/logs.
  - [ ] redaction policy is enforced in both normal and verbose modes.
  - [ ] all auth commands use shared redaction helpers (no per-command ad-hoc masking logic).
  - [ ] failures and debug output preserve diagnostic value without leaking sensitive material.
- [ ] Secret memory-lifetime requirements:
  - [ ] secret is held in memory only for minimal required path (read -> validate -> store/sign -> clear).
  - [ ] use `[]byte` for sensitive secret handling where practical.
  - [ ] clear sensitive buffers after use (best-effort wipe) and avoid unnecessary copies.
  - [ ] avoid propagating secret values into long-lived structs, error objects, or log contexts.
- [ ] Persistence boundaries:
  - [ ] no secret material is written to repo-tracked files or plain profile config.
  - [ ] profile config stores metadata only (profile name, timestamps, backend marker, active profile).
  - [ ] on macOS/Linux, config path is `~/.wbcli/config.yaml`.
  - [ ] config file at `~/.wbcli/config.yaml` is created/updated with owner-only permissions (`0600`).
- [ ] Unit tests and command tests include success and negative paths for each command part.
- [ ] Documentation requirements for implementation:
  - [ ] README must explain auth secret input modes in simple, easy-to-understand language.
  - [ ] README must include examples for local interactive use (prompt) and automation/CI use (`--api-secret-stdin`).
  - [ ] README must explicitly warn not to pass secrets via command arguments.
  - [ ] Cobra command help must include practical `Example` blocks for `auth login/use/profiles list/logout/current/test`.

Test Matrix:
- [ ] `auth login`: interactive secret input success, stdin secret input success, missing profile, invalid profile format (spaces/special chars/unicode/too long), missing key, empty secret, keychain unavailable, permission denied.
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
  - [ ] verify `~/.wbcli/config.yaml` contains metadata only and no secret values.
  - [ ] verify `~/.wbcli/config.yaml` permissions are `0600` on macOS/Linux.
- [ ] platform compatibility:
  - [ ] run `auth login/use/profiles list/logout/current` verification on macOS.
  - [ ] run `auth login/use/profiles list/logout/current` verification on Linux.
  - [ ] Windows verification is optional and recorded only if performed.
  - [ ] include keychain unavailable and permission-denied cases for both macOS and Linux.
  - [ ] include redaction no-leak assertions on both macOS and Linux.
  - [ ] include `~/.wbcli/config.yaml` metadata-only and `0600` assertions on both macOS and Linux.
- [ ] `auth use`: existing profile selection, missing profile failure, active-profile metadata update.
- [ ] `auth profiles list`: returns metadata-only rows, redaction assertions, empty state.
- [ ] `auth logout`: existing profile removal, missing profile idempotency, permission denied.
- [ ] `auth current`: correct active profile output, empty state behavior, redaction assertions.
- [ ] `auth test`: valid credentials, auth failure, network timeout, backend unavailable, redaction in error output.
- [ ] cross-cutting: verify no sensitive value appears in stdout/stderr/log buffers under verbose and non-verbose modes.
- [ ] redaction contract:
  - [ ] assert shared redaction helper is used by all auth command output paths.
  - [ ] assert sensitive values are absent from command output and logs in success and error scenarios.
- [ ] secret memory contract:
  - [ ] verify secret buffers are cleared after use on primary paths.
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
4. Implement `auth login` flags/input path only (hidden prompt + stdin option) with tests.
5. Implement `auth login` validation + storage write path via `AuthLoginService` with tests.
5.1. Add secret buffer lifecycle handling (minimal lifetime + best-effort wipe + no unnecessary copies).
5.2. Add credential overwrite guardrails (interactive confirm and non-interactive `--force` behavior) with leak-safe output handling.
6. Implement `auth profiles list` metadata-only read path with redaction tests.
7. Implement `auth use` active-profile selection path with metadata-only update tests.
8. Implement `auth logout` delete path and idempotency behavior with tests.
9. Implement `auth current` metadata read path with safe output tests.
9.1. Implement shared auth redaction helper and wire all auth command outputs/logging through it.
10. Add command-level docs/help text updates for secure usage.
10.1. Update README with clear simple-language explanation of secret input behavior (prompt vs stdin), with safe examples and warning against command-argument secrets.
10.2. Add/verify `cobra.Command.Example` text for all auth commands with safe usage examples.
11. Run verification for `login/use/profiles list/logout/current` and capture evidence.
11.1. Align auth security verification with existing CI test/build pipelines and extend CI where needed (do not create a separate disconnected verification flow).
11.2. Ensure CI artifacts/summaries capture platform-specific security evidence for macOS and Linux.
12. After PROJ-2026-002 is ready, implement `auth test` command using `AuthProbe` port, then add redaction/failure-classification tests as the final step.
13. Run full verification including `auth test` and capture final evidence.

Verification Evidence (Required In Review):
- `go test ./...`
- `go build .`
- command-level tests for `auth login/use/profiles list/logout/current/test` with redaction assertions
- explicit proof that no secret values are persisted in config files
- explicit proof that `~/.wbcli/config.yaml` contains metadata only and no secrets
- explicit proof that `~/.wbcli/config.yaml` uses `0600` permissions on macOS/Linux
- explicit proof that keychain-unavailable path fails closed with no implicit fallback write
- explicit proof that redaction contract is enforced in normal and verbose modes
- explicit proof that secret memory-lifetime controls are implemented (minimal lifetime, cleared buffers, no secret in errors/logs)
- README excerpt/evidence showing simple-language auth input guidance and safe usage examples
- CLI help evidence showing `Example` blocks for auth commands
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
- 2026-02-26: Updated command model examples (`auth login/use/profiles list/logout/current`) and reordered rollout to start from `auth login`.
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

# PROJ-2026-016: Remove auth profiles and switch to single-session auth

ID: PROJ-2026-016
Title: Remove auth profiles and switch to single-session auth
Priority: P1
Status: Done
Owner: nocle
Due Date: 2026-03-03
Created: 2026-02-28
Updated: 2026-02-28
Links: [PROJ-2026-001](./PROJ-2026-001-implement-secure-api-key-storage-adapter.md), [CLI Design](../../docs/cli-design.md)

Problem:
Current auth flow is profile-based (`login/use/list/current/logout`) and requires profile selection/validation. Product direction is a simplified single-user CLI where user state is only `logged in` or `logged out`.

Outcome:
`wbcli` auth works in single-session mode with no profile concept in commands, services, ports, adapters, config schema, tests, and docs.

Relation To PROJ-2026-001:
- This ticket is a scope override for the profile-oriented parts of `PROJ-2026-001`.
- `PROJ-2026-001` must still be finished, but its remaining acceptance criteria/checklist must be updated to match this single-session model.
- Any `PROJ-2026-001` requirements that depend on profiles (`auth use/list/current`, `--profile`, active profile metadata) are replaced by this ticket's requirements.

In Scope:
- remove profile checks and profile validation from auth flows
- remove profile-based auth commands and flags
- refactor auth storage interfaces to single credential slot
- refactor config metadata from profile map to single auth-session record
- update tests, docs, and ticket checklists to single-session behavior

Out Of Scope:
- OAuth implementation
- WhiteBIT partner/OAuth integration
- backward compatibility for removed profile commands

Target Command Model:
- `wbcli auth login`
- `wbcli auth logout`
- `wbcli auth status`
- `wbcli auth test`

Command Changes:
- remove: `wbcli auth use`
- remove: `wbcli auth list`
- remove: `wbcli auth current`
- remove `--profile` from auth commands
- keep stdin-only contract for `auth login` (`api_key` first line, `api_secret` second line)

Acceptance Criteria:
- [x] CLI surface is single-session:
  - [x] `auth use/list/current` are removed from command tree and help output.
  - [x] `auth status` exists and reports logged-in/logged-out state with safe metadata only.
  - [x] no auth command accepts `--profile`.
- [x] Auth domain/service logic has no profile dependency:
  - [x] profile name validation is removed from auth flows.
  - [x] no active-profile selection logic remains.
  - [x] auth state is represented as single session only.
- [x] Ports and adapters are single-session:
  - [x] `CredentialStore` no longer takes profile parameters.
  - [x] profile metadata store is replaced by single-session metadata store.
  - [x] keychain adapter uses one fixed credential record key for `wbcli`.
- [x] Persistence behavior:
  - [x] config contains single auth-session metadata, no profile map.
  - [x] config remains metadata-only and never stores secrets.
  - [x] config path remains `~/.wbcli/config.yaml` with `0600` on macOS/Linux.
- [x] Command behavior:
  - [x] `auth login` writes credentials to secure store and sets logged-in session metadata.
  - [x] `auth login` overwrites existing session credential by default.
  - [x] `auth logout` clears credential and session metadata; idempotent when already logged out.
  - [x] `auth test` keeps current scope and dependency on PROJ-2026-002.
- [x] Test and docs coverage:
  - [x] tests cover single-session command success and failure paths.
  - [x] docs/examples are updated to no-profile commands.
  - [x] `PROJ-2026-001` checklist is reconciled to this new scope.

Test Matrix:
- [x] `auth login` valid stdin payload creates session.
- [x] `auth login` invalid stdin payload fails with clear error.
- [x] `auth login` overwrites existing session when already logged in.
- [x] `auth logout` succeeds when logged in.
- [x] `auth logout` succeeds (idempotent) when logged out.
- [x] `auth status` reports logged-out state.
- [x] `auth status` reports logged-in safe metadata.
- [x] keychain unavailable/permission-denied paths return actionable errors.
- [x] config assertions confirm metadata-only and `0600` permissions.
- [x] removed commands (`use/list/current`) return unknown command errors.

Risks:
- profile-based user scripts break immediately due to command removal
- refactor can leave stale profile references in docs/tests if not fully audited
- keychain record key migration may require re-login after deployment

Rollout Plan:
1. Remove profile flags and profile commands from `cmd/auth*`.
2. Introduce `auth status` command and output contract.
3. Refactor app ports/services from profile-based to single-session interfaces.
4. Refactor keychain and config adapters to single-session storage model.
5. Update tests for command/service/adapter behavior in single-session mode.
6. Update README and `docs/cli-design.md` to single-session commands.
7. Reconcile and update `PROJ-2026-001` checklist items affected by this scope change.

Rollback Plan:
1. Revert commits from this ticket.
2. Restore previous profile-based auth command set and interfaces.
3. Require user to re-login if keychain/config format changed during rollout.

Status Notes:
- 2026-02-28: Created in Ready.
- 2026-02-28: Scoped as mandatory simplification to single-session auth (no profiles, no backward compatibility).
- 2026-02-28: Marked as scope override for profile-related parts of PROJ-2026-001.
- 2026-02-28: Implemented single-session auth architecture and command surface (`login/logout/status/test`), with profile flags and profile commands removed.
- 2026-02-28: Removed `--force` from `auth login`; login now overwrites existing session by default.
- 2026-02-28: Added command-level tests for `auth logout` (logged-in + idempotent), `auth status` logged-in output, and actionable unavailable/permission-denied auth errors; reconciled PROJ-2026-001 checklist to single-session scope.
- 2026-02-28: Closed: implementation complete and manually confirmed login/logout working on macOS; Linux manual verification tracked separately in PROJ-2026-017.

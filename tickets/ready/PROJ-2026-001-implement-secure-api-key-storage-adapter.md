# PROJ-2026-001: Implement secure API key storage adapter

ID: PROJ-2026-001
Title: Implement secure API key storage adapter
Priority: P1
Status: Ready
Owner: chewbaccalol
Due Date: 2026-03-02
Created: 2026-02-26
Updated: 2026-02-26
Links: [CLI Design](../../docs/cli-design.md), [PROJ-2026-014](./PROJ-2026-014-define-credential-encryption-policy-and-encrypted-file-fallback-backend.md), [PROJ-2026-015](../backlog/PROJ-2026-015-implement-credential-access-controls-session-unlock-and-key-rotation-workflow.md)

Problem:
Trading commands require API credentials, but storing secrets in plaintext config or shell history is unsafe.

Outcome:
`whitbit keys` commands manage credentials via OS secret storage with profile isolation and safe error handling, aligned with encryption and access hardening tickets.

Acceptance Criteria:
- [ ] `whitbit keys set/list/remove/test` commands are implemented.
- [ ] API key/secret are written to secure storage, never persisted in repo-tracked files.
- [ ] Storage backend preference is enforced (`os-keychain` first, fallback backend from PROJ-2026-014 only when needed).
- [ ] Command input supports non-echo secret entry to avoid shell history leakage.
- [ ] `keys test` performs authenticated connectivity check without printing sensitive headers.
- [ ] Unit tests cover success, missing profile, missing secret backend, and permission denied cases.

Risks:
- Secret backend availability differs by OS and CI environment.
- Improper logging can still leak headers or partial secrets.

Rollout Plan:
1. Implement secret store interface and platform adapter.
2. Integrate with profile metadata store and backend selection policy.
3. Add command handlers and tests.

Rollback Plan:
1. Disable write operations and keep read-only mode if backend instability appears.
2. Revert to previous credential path only if a secure alternative is confirmed.

Status Notes:
- 2026-02-26: Created in Ready.
- 2026-02-26: Refined with explicit secure I/O and testability criteria.
- 2026-02-26: Linked to encryption/access hardening tickets.

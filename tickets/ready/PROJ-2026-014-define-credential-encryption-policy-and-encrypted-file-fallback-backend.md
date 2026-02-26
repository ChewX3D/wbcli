# PROJ-2026-014: Define credential encryption policy and encrypted-file fallback backend

ID: PROJ-2026-014
Title: Define credential encryption policy and encrypted-file fallback backend
Priority: P1
Status: Ready
Owner: chewbaccalol
Due Date: 2026-03-01
Created: 2026-02-26
Updated: 2026-02-26
Links: [CLI Design](../../docs/cli-design.md), [WhiteBIT Integration](../../docs/whitebit-integration.md), [WhiteBIT HTTP Auth](https://docs.whitebit.com/private/http-auth)

Problem:
The current plan says to use secure storage, but it does not define concrete encryption rules for fallback environments where OS secret store is unavailable.

Outcome:
Credential storage follows a strict policy: OS keychain first, encrypted-file fallback second, with auditable cryptographic parameters and migration rules.

Acceptance Criteria:
- [ ] Security policy document defines supported backends and selection order: `os-keychain` default, `encrypted-file` fallback only.
- [ ] Encrypted-file backend uses AES-256-GCM with per-record random nonce and authenticated metadata.
- [ ] Encryption key is derived from passphrase using Argon2id with configurable memory/time/parallelism and random salt.
- [ ] Credential file permissions are locked to owner-only (`0600`) and format versioning supports future migration.
- [ ] Migration command exists to move credentials from encrypted-file to OS keychain when keychain becomes available.
- [ ] Tests validate decrypt failure on tampering, wrong passphrase, and stale schema version.

Risks:
- Weak KDF settings can make offline brute force practical if file is leaked.
- Custom crypto implementation bugs can silently corrupt credentials.

Rollout Plan:
1. Finalize security policy and cryptographic parameter defaults.
2. Implement encrypted-file backend behind explicit configuration.
3. Add migration path and negative-path tests.

Rollback Plan:
1. Disable encrypted-file backend and require OS keychain only.
2. Keep migration utility for existing encrypted-file users.

Status Notes:
- 2026-02-26: Created in Ready.
- 2026-02-26: Added to formalize encryption at rest and fallback behavior.

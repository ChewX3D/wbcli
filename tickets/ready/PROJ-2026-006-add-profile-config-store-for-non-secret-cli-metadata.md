# PROJ-2026-006: Add profile config store for non-secret CLI metadata

ID: PROJ-2026-006
Title: Add profile config store for non-secret CLI metadata
Priority: P1
Status: Ready
Owner: chewbaccalol
Due Date: 2026-03-01
Created: 2026-02-26
Updated: 2026-02-26
Links: [CLI Design](../../docs/cli-design.md)

Problem:
We need profile management for account/environment targeting, but secrets must stay outside git-tracked config.

Outcome:
CLI persists only non-secret profile metadata locally (default profile, timestamps, labels), while secret values remain in secret storage.

Acceptance Criteria:
- [ ] Config file is stored in OS-appropriate app config path and ignored by git.
- [ ] Profile metadata schema is versioned for future migrations.
- [ ] Commands support creating/listing/selecting/removing profiles without storing API secrets.
- [ ] Corrupt or missing config file paths return actionable errors.
- [ ] Unit tests cover load/save/migrate/error behavior.

Risks:
- Cross-platform path handling can break profile discovery.
- Schema changes may introduce migration bugs if versioning is skipped.

Rollout Plan:
1. Define metadata schema and storage path resolver.
2. Implement repository interface with file-backed adapter.
3. Wire profile commands and add tests.

Rollback Plan:
1. Disable profile persistence and fall back to explicit `--profile` only mode.
2. Revert schema migration changes and restore last stable format.

Status Notes:
- 2026-02-26: Created in Ready.
- 2026-02-26: Sequenced after PROJ-2026-005 as part of Phase 0 minimal setup.

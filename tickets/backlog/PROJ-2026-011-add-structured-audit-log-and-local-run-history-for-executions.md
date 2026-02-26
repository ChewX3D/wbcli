# PROJ-2026-011: Add structured audit log and local run history for executions

ID: PROJ-2026-011
Title: Add structured audit log and local run history for executions
Priority: P2
Status: Backlog
Owner: chewbaccalol
Due Date: 2026-03-12
Created: 2026-02-26
Updated: 2026-02-26
Links: [CLI Design](../../docs/cli-design.md)

Problem:
Operational support needs traceability of trading actions; without run history, debugging failed or partial submissions is difficult.

Outcome:
Each command execution records sanitized audit events and a retrievable local history for support and incident analysis.

Acceptance Criteria:
- [ ] Every live order command writes a structured audit record with `request_id`, profile, mode, and outcome.
- [ ] Sensitive fields (API keys, signatures, secrets) are never persisted.
- [ ] `wbcli history list` supports filtering by date/profile/mode.
- [ ] Retention setting exists for pruning old records.
- [ ] Tests verify redaction and history query behavior.

Risks:
- Improper redaction can leak sensitive values.
- History growth may impact local performance if retention is unmanaged.

Rollout Plan:
1. Define audit event schema and local storage backend.
2. Add write hooks in order execution paths.
3. Implement history query command and retention settings.

Rollback Plan:
1. Disable persistence and keep in-memory logging only.
2. Preserve command output while removing history command until fixed.

Status Notes:
- 2026-02-26: Created in Backlog.
- 2026-02-26: Added to support incident analysis and reproducibility in full-scale phase.

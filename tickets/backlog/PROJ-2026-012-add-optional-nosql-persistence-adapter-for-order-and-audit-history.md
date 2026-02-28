# PROJ-2026-012: Add optional NoSQL persistence adapter for order and audit history

ID: PROJ-2026-012
Title: Add optional NoSQL persistence adapter for order and audit history
Priority: P3
Status: Backlog
Owner: chewbaccalol
Due Date: 2026-03-16
Created: 2026-02-26
Updated: 2026-02-28
Links: [CLI Design](../../docs/cli-design.md)

Problem:
If team workflows require shared history across machines/users, local-only storage is insufficient.

Outcome:
A feature-gated NoSQL adapter can persist order/audit records centrally while keeping local mode as default. DynamoDB is a supported option for this ticket.

Acceptance Criteria:
- [ ] Storage interface supports local backend and NoSQL backend without changing command behavior.
- [ ] DynamoDB adapter is implemented as the first NoSQL option.
- [ ] NoSQL backend is activated only when explicit config enables it.
- [ ] Table/index design supports query patterns for `request_id`, profile, and timestamp.
- [ ] Write paths are idempotent by execution key to avoid duplicates.
- [ ] Failure to connect to NoSQL backend degrades safely with clear error and no partial secret leakage.

Risks:
- Added infrastructure complexity may not be justified for solo/local use.
- Network/database outage can block history writes unless fallback policy is clear.

Rollout Plan:
1. Design storage abstraction and local-first default.
2. Implement DynamoDB adapter behind feature flag.
3. Keep adapter contracts NoSQL-oriented so other NoSQL backends can be added later.
4. Add integration tests for the enabled NoSQL backend in optional CI job.

Rollback Plan:
1. Disable NoSQL backend flag and continue with local storage backend.
2. Keep schema migration scripts reversible for rollback.

Status Notes:
- 2026-02-26: Created in Backlog.
- 2026-02-26: Classified as optional infrastructure step for multi-user scaling.
- 2026-02-28: Reframed from single-database scope to NoSQL adapter scope; DynamoDB marked as primary supported option.

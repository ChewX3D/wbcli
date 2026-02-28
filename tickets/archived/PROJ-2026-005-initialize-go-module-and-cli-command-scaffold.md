# PROJ-2026-005: Initialize Go module and CLI command scaffold

ID: PROJ-2026-005
Title: Initialize Go module and CLI command scaffold
Priority: P0
Status: Archived
Owner: chewbaccalol
Due Date: 2026-02-28
Created: 2026-02-26
Updated: 2026-02-28
Links: [CLI Design](../../docs/cli-design.md)

Problem:
The repository has planning docs and ticket automation, but no Go codebase to implement the CLI.

Outcome:
A buildable Go CLI skeleton exists with clear package boundaries (`cmd`, `internal/domain`, `internal/app`, `internal/adapters`) so feature tickets can be implemented incrementally.

Acceptance Criteria:
- [ ] `go.mod` and `go.sum` are committed with module path and pinned dependencies.
- [ ] Root package builds and supports `wbcli --help` with `keys` and `order` command groups stubbed.
- [ ] Project layout includes baseline packages matching architecture in `README.md`.
- [ ] A `make build` or equivalent documented command builds the CLI locally.
- [ ] Initial README usage section is updated for local build/run.

Risks:
- CLI framework choice could create later migration cost if not evaluated early.
- Early package boundaries can become rigid if adapters are over-coupled to CLI flags.

Rollout Plan:
1. Create module and baseline directories.
2. Add root command and subcommand stubs.
3. Verify build path and update README.

Rollback Plan:
1. Revert scaffold commit.
2. Keep ticket open with documented blocker and revised structure proposal.

Status Notes:
- 2026-02-26: Created in Ready.
- 2026-02-26: Marked as Phase 0 foundation for all implementation tickets.
- 2026-02-28: Canceled by requester.
- 2026-02-28: Canceled by requester (confirmed).

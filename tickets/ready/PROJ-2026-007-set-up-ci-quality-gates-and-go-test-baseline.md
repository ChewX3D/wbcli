# PROJ-2026-007: Set up CI quality gates and Go test baseline

ID: PROJ-2026-007
Title: Set up CI quality gates and Go test baseline
Priority: P1
Status: Ready
Owner: chewbaccalol
Due Date: 2026-03-02
Created: 2026-02-26
Updated: 2026-02-26
Links: [AGENTS](../../AGENTS.md), [Ticket Workflow](../README.md)

Problem:
Without automated quality gates, implementation speed will increase defect risk and break Definition of Done.

Outcome:
Protected-branch quality checks are in place for format, lint/vet, tests, and secret/dependency scanning.

Acceptance Criteria:
- [ ] CI workflow runs `gofmt` check, `go vet`, and `go test ./...` on PRs.
- [ ] Secret scan and dependency vulnerability scan run in CI.
- [ ] CI status is referenced in `README.md` with contributor instructions.
- [ ] Failing checks block merge for protected branches.
- [ ] One intentionally failing sample check has been validated and removed after verification.

Risks:
- Slow CI may reduce iteration speed and encourage bypass behavior.
- Overly strict checks early can stall bootstrap if not tuned.

Rollout Plan:
1. Add CI workflow and baseline scripts.
2. Validate pass/fail behavior with sample PR branch.
3. Document local pre-check commands.

Rollback Plan:
1. Temporarily disable non-critical scan steps if they are unstable.
2. Keep core build/test gates enabled while remediating scanners.

Status Notes:
- 2026-02-26: Created in Ready.
- 2026-02-26: Added as Phase 0 quality gate before full MVP feature velocity.

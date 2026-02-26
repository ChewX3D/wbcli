# PROJ-2026-008: Implement range live submission via collateral bulk order endpoint

ID: PROJ-2026-008
Title: Implement range live submission via collateral bulk order endpoint
Priority: P1
Status: Backlog
Owner: chewbaccalol
Due Date: 2026-03-09
Created: 2026-02-26
Updated: 2026-02-26
Links: [CLI Design](../../docs/cli-design.md), [WhiteBIT Integration](../../docs/whitebit-integration.md)

Problem:
Dry-run planning is not enough for production trading workflows; range plans must be submitted safely with controlled execution.

Outcome:
CLI submits validated range plans to the bulk endpoint with confirmation, chunking, and partial-failure reporting.

Acceptance Criteria:
- [ ] Planner output maps correctly to `/api/v4/collateral-bulk-limit-order` payload.
- [ ] Live submission requires explicit `--confirm` or interactive confirmation.
- [ ] Large batches are chunked with deterministic ordering and request correlation IDs.
- [ ] Partial failures are reported per order with retry guidance.
- [ ] Integration tests with mock server cover chunk success, chunk failure, and mixed outcomes.

Risks:
- Partial successes can leave position exposure inconsistent with expected plan.
- Endpoint limits may change and break chunk sizing assumptions.

Rollout Plan:
1. Implement submission orchestrator and chunk policy.
2. Add confirmation and safety prompts.
3. Add error reconciliation output and tests.

Rollback Plan:
1. Revert command to dry-run-only mode.
2. Disable automatic chunk retries while keeping reporting.

Status Notes:
- 2026-02-26: Created in Backlog.
- 2026-02-26: Planned as first full-scale execution ticket after dry-run planner stabilizes.

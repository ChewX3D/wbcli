# PROJ-2026-010: Add cancel and reconcile commands for batch orders

ID: PROJ-2026-010
Title: Add cancel and reconcile commands for batch orders
Priority: P2
Status: Backlog
Owner: chewbaccalol
Due Date: 2026-03-11
Created: 2026-02-26
Updated: 2026-02-26
Links: [CLI Design](../../docs/cli-design.md)

Problem:
After submitting many orders, operators need fast cancellation and reconciliation flows to manage risk and confirm execution state.

Outcome:
CLI supports batch-safe cancellation and status reconciliation based on client-order-id prefixes and execution IDs.

Acceptance Criteria:
- [ ] `whitbit order cancel` supports cancel by explicit order id and by client-order-id prefix.
- [ ] `whitbit order reconcile` retrieves status for planned/submitted orders and prints mismatch summary.
- [ ] Commands include `--dry-run` preview for cancellation targets.
- [ ] Idempotent behavior is documented for already-canceled/already-filled orders.
- [ ] Integration tests cover cancel success, partial cancel, and not-found paths.

Risks:
- Cancellation race conditions may produce stale reconciliation snapshots.
- Bulk cancel targeting errors can cancel unintended orders.

Rollout Plan:
1. Add cancellation/reconciliation use-cases in app layer.
2. Implement adapter calls and safe target filtering.
3. Add summary renderer and test coverage.

Rollback Plan:
1. Restrict cancellation to explicit IDs only.
2. Disable prefix cancellation until filtering is verified.

Status Notes:
- 2026-02-26: Created in Backlog.
- 2026-02-26: Sequenced after bulk range submission to close lifecycle loop.

# PROJ-2026-009: Add geometric and capped-geometric amount modes with exposure cap checks

ID: PROJ-2026-009
Title: Add geometric and capped-geometric amount modes with exposure cap checks
Priority: P2
Status: Backlog
Owner: chewbaccalol
Due Date: 2026-03-10
Created: 2026-02-26
Updated: 2026-02-26
Links: [CLI Design](../../docs/cli-design.md)

Problem:
Advanced scaling strategies are needed for full implementation, but unbounded geometric growth can exceed risk tolerance quickly.

Outcome:
Range planner supports `geometric` and `capped-geometric` modes with explicit risk caps and clear validation feedback.

Acceptance Criteria:
- [ ] `--amount-mode geometric --ratio` is implemented with deterministic multiplier math.
- [ ] `--amount-mode capped-geometric --ratio --max-multiplier` is implemented and enforced.
- [ ] Exposure guardrails enforce max total amount and max notional per plan.
- [ ] Dry-run output shows multiplier and cumulative exposure per step.
- [ ] Unit tests validate progression formulas and cap enforcement.

Risks:
- Misinterpreted multipliers can create unintended large orders.
- Exposure calculation errors can bypass safety limits.

Rollout Plan:
1. Extend domain amount-mode model and validators.
2. Update range planner output contract.
3. Add formula/exposure regression tests.

Rollback Plan:
1. Hide new modes behind feature flag.
2. Default unsupported modes back to validation error.

Status Notes:
- 2026-02-26: Created in Backlog.
- 2026-02-26: Marked as post-MVP enhancement with strict risk controls.

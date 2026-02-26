# PROJ-2026-004: Add order range planner with constant and arithmetic amount modes

ID: PROJ-2026-004
Title: Add order range planner with constant and arithmetic amount modes
Priority: P1
Status: Backlog
Owner: chewbaccalol
Due Date: 2026-03-06
Created: 2026-02-26
Updated: 2026-02-26
Links: [CLI Design](../../docs/cli-design.md)

Problem:
Users need laddered range plans, but generating orders manually is slow and error-prone.

Outcome:
`whitbit order range` can generate deterministic constant/arithmetic plans in local dry-run mode with exposure summary.

Acceptance Criteria:
- [ ] Range generation supports `start-price`, `end-price`, `step`, `amount-mode constant|arithmetic`.
- [ ] Arithmetic mode supports `start-multiplier` and `step-multiplier`.
- [ ] Planner validates max order count and rejects invalid price-step combinations.
- [ ] `--dry-run` prints per-order preview plus aggregate totals without API submission.
- [ ] Unit tests cover ascending/descending ranges and edge-case validation failures.

Risks:
- Floating-point precision drift can create invalid prices/amounts.
- Large ranges could consume excessive memory if plan generation is not bounded.

Rollout Plan:
1. Implement range planner in domain layer.
2. Add command wiring and dry-run renderer.
3. Add unit tests for all formulas and bounds checks.

Rollback Plan:
1. Temporarily support constant mode only if arithmetic behavior is unstable.
2. Revert planner formula changes and retain validation safeguards.

Status Notes:
- 2026-02-26: Created in Backlog.
- 2026-02-26: Tagged as Phase 1.5 after single-order MVP.

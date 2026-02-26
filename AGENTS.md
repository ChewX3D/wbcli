# AGENTS.md

## Purpose And Principles

This file defines the operating system for project execution across any software repo.

Principles:

- ship in small, reviewable increments
- keep work traceable from idea to release
- enforce quality gates before merge
- keep docs and code synchronized
- optimize for sustainability, not heroics

## Roles

Use these roles even on a one-person team (one person may hold multiple roles):

- `Requester`: defines outcome and business context
- `Owner`: accountable for ticket delivery
- `Reviewer`: validates correctness, risk, and maintainability
- `Release Steward`: handles release notes, tagging, and rollback readiness

## Ticket System

Every non-trivial change needs a ticket.

Required fields:

- `ID`: `PROJ-YYYY-NNN` (example: `CORE-2026-014`)
- `Title`: action-oriented, specific
- `Priority`: `P0`, `P1`, `P2`, `P3`
- `Status`: `Backlog`, `Ready`, `In Progress`, `Review`, `Blocked`, `Done`, `Archived`
- `Owner`: single accountable person
- `Due Date`: ISO date (`YYYY-MM-DD`) or `None`
- `Acceptance Criteria`: measurable checklist
- `Links`: design doc, PR, incident, related tickets

Rules:

- no ticket ID, no merge (except emergency hotfix with retroactive ticket)
- each PR references exactly one primary ticket
- status changes must include a short note

## Backlog And Todo Workflow

### Triage

- triage new items at least 2 times per week
- reject or archive vague items without clear value
- convert ideas into tickets only when problem + outcome are explicit

### Planning

- maintain a `Ready` queue with top priorities for next 1-2 weeks
- assign owner before moving to `In Progress`
- split tickets estimated over 2 days of focused work

### WIP Limits

- max `In Progress` per owner: `2`
- no new start while blocked ticket lacks unblock action

### Archive Rules

- move done tickets to `Archived` after 30 days
- close stale backlog items after 60 days without activity

## Definition Of Ready (DoR)

A ticket is `Ready` only if:

- problem statement is clear
- acceptance criteria are testable
- dependencies are identified
- rollout and rollback notes exist (for risky changes)
- owner and priority are set

## Definition Of Done (DoD)

A ticket is `Done` only if:

- code merged to main branch
- required tests pass
- docs/changelog updated
- monitoring/alerts considered (if production-impacting)
- acceptance criteria verified

## Branching And Commit Conventions

Branch naming:

- `feature/<ticket-id>-<slug>`
- `fix/<ticket-id>-<slug>`
- `chore/<ticket-id>-<slug>`
- `hotfix/<ticket-id>-<slug>`

Commit format:

- `<type>(<scope>): <short summary>`
- types: `feat`, `fix`, `chore`, `docs`, `test`, `refactor`, `perf`, `build`, `ci`
- include ticket ID in commit body or footer (`Refs: PROJ-2026-014`)

Rules:

- prefer small atomic commits
- no direct commits to protected main branch

## PR Checklist And Review Standards

PR must include:

- linked ticket
- what changed and why
- risk level (`low`, `medium`, `high`)
- test evidence
- rollback plan (if needed)

Reviewer checks:

- correctness and edge cases
- security/privacy impact
- migration compatibility
- test quality (not only coverage)
- docs consistency

Merge policy:

- at least 1 reviewer approval
- all required CI checks green
- unresolved comments prohibited

## Testing Expectations By Change Type

- `docs-only`: spelling/link checks
- `ui/ux`: component tests + visual/manual sanity
- `api/logic`: unit + integration tests for success/failure paths
- `data/schema`: migration test + backward compatibility check
- `infra/ci`: dry-run in non-prod path where possible
- `hotfix`: minimal reproducer test added within 24 hours post-fix

## Documentation Update Policy

When behavior changes, update in same PR:

- user-facing usage docs
- architecture/design notes when boundaries change
- runbook/ops docs when operational steps change

If no doc impact, PR must explicitly state: `Docs impact: none`.

## Release Cadence And Changelog

Cadence:

- default: weekly or bi-weekly release train
- hotfix releases as needed

Changelog process:

- maintain `Unreleased` section continuously
- group entries by `Added`, `Changed`, `Fixed`, `Removed`
- tag release with version and date
- link tickets and PRs in release notes

## Incident And Bug Handling

Severity:

- `SEV1`: critical outage/security
- `SEV2`: major degradation
- `SEV3`: minor issue/workaround exists

Flow:

1. open incident ticket immediately
2. stabilize service (mitigate first, optimize later)
3. capture timeline and decisions
4. ship fix with linked ticket/PR
5. complete postmortem for SEV1/SEV2 within 5 business days

## Metrics And KPIs

Delivery metrics:

- lead time (ticket start to production)
- cycle time by status
- throughput (done tickets per week)
- WIP aging

Quality metrics:

- change failure rate
- escaped defects by severity
- mean time to recovery (MTTR)
- flaky test rate
- documentation freshness (tickets merged with doc updates where required)

## Automation Expectations

CI minimum gates:

- build/lint/test required for protected branches
- secret scanning and dependency vulnerability scan
- PR template enforcement

Workflow automation:

- auto-label by path/type where possible
- stale ticket reminder at 7 days inactivity
- auto-close stale tickets at 30 days after reminder unless exempt
- scheduled reminder for weekly and monthly rituals

## Weekly And Monthly Rituals

Weekly (30-45 min):

- review in-progress items and blockers
- check aging tickets and rebalance priorities
- review escaped bugs and flaky tests
- confirm release readiness

Monthly (60 min):

- review KPI trends
- prune backlog and archive stale work
- audit test gaps and incident follow-ups
- update process rules causing repeated friction

## Templates

### Ticket Template

```md
ID: PROJ-YYYY-NNN
Title:
Priority: P0|P1|P2|P3
Status: Backlog|Ready|In Progress|Review|Blocked|Done|Archived
Owner:
Due Date: YYYY-MM-DD | None
Links: [Design Doc] [PR] [Related Tickets]

Problem:

Outcome:

Acceptance Criteria:
- [ ]
- [ ]

Risks:

Rollout Plan:

Rollback Plan:
```

### Todo Item Template

```md
- [ ] <short action> | Owner: <name> | Ticket: <ID> | Due: <YYYY-MM-DD>
```

### PR Template

```md
## Summary

## Ticket
- Primary: <ID>

## Changes
- 

## Risk
- Level: low|medium|high
- Key risks:

## Testing
- [ ] Unit
- [ ] Integration
- [ ] Manual
- Evidence:

## Docs
- [ ] Updated
- [ ] Not needed (reason): 

## Rollback

## Reviewer Checklist
- [ ] Acceptance criteria met
- [ ] Edge cases covered
- [ ] Security/privacy considered
- [ ] Observability impact considered
```

### Weekly Review Template

```md
Week Of:

Completed:
- 

In Progress:
- 

Blocked:
- 

KPIs Snapshot:
- Lead time:
- Throughput:
- Escaped defects:
- Flaky tests:

Next Week Focus:
- 
```

### Postmortem Template

```md
Incident ID:
Severity: SEV1|SEV2|SEV3
Date:
Owners:

Summary:

Customer Impact:

Timeline:
- HH:MM -

Root Cause:

Contributing Factors:

What Went Well:

What Failed:

Action Items:
- [ ] <action> | Owner: <name> | Due: <date> | Ticket: <ID>
```

## Quick Start In 30 Minutes

1. Create ticket board columns: `Backlog`, `Ready`, `In Progress`, `Review`, `Blocked`, `Done`, `Archived`.
2. Add ticket/PR/postmortem templates from this file to your platform.
3. Define branch protection: PR required, CI required, no direct main pushes.
4. Enable CI checks: lint, test, build, secret scan.
5. Start triage: label current work with IDs, priorities, owners, due dates.
6. Set WIP rule (`max 2`) and stale reminder automation.
7. Schedule weekly (30-45 min) and monthly (60 min) operating reviews.

## Common Failure Modes And Prevention Rules

- Work starts without clear acceptance criteria:
  - prevention: enforce DoR before `In Progress`
- Too many parallel tasks, little completion:
  - prevention: strict WIP limit and weekly aging review
- PRs merge with hidden risk:
  - prevention: mandatory risk section + rollback plan
- Docs drift from behavior:
  - prevention: docs update required in DoD
- Repeated incidents with no learning:
  - prevention: required postmortem actions tracked as tickets
- Release pain accumulates:
  - prevention: fixed cadence + changelog discipline

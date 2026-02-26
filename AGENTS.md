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

Repository implementation:

- ticket files live under `tickets/<status>/`
- generated board is `tickets/board.md`
- create ticket via `./scripts/tickets/new.sh`
- move ticket via `./scripts/tickets/move.sh`
- rebuild board via `./scripts/tickets/board.sh`

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

Team model:

- team size: 2 developers
- trunk-based development: both developers commit directly to `main`
- keep commits small and focused to reduce merge friction

Commit format:

- `<type>(<scope>): <short summary>`
- types: `feat`, `fix`, `chore`, `docs`, `test`, `refactor`, `perf`, `build`, `ci`
- include ticket ID in commit body or footer (`Refs: PROJ-2026-014`)

Rules:

- prefer small atomic commits
- before starting work, run `git pull --rebase origin main`
- after local commit and before push, run `git pull --rebase origin main`
- if pull introduces conflicts, resolve them locally, re-run required checks, then push
- push only when local `main` is up to date with `origin/main`

## Change Checklist And Review Standards

For trunk-based direct commits to `main`, each change set (commit or small commit series) must include:

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

- prefer 1 peer review before push when practical; if not practical, request review immediately after push
- all required CI checks must pass on `main`
- if CI fails after push, prioritize immediate fix or revert

## Testing Expectations By Change Type

- `docs-only`: spelling/link checks
- `ui/ux`: component tests + visual/manual sanity
- `api/logic`: unit + integration tests for success/failure paths
- `data/schema`: migration test + backward compatibility check
- `infra/ci`: dry-run in non-prod path where possible
- `hotfix`: minimal reproducer test added within 24 hours post-fix

## Go Code Style (Google)

VERY IMPORTANT (MANDATORY):

- all Go code must follow Google Go style guidance
- style violations are blocking issues and must be fixed before merge

For all Go code in this repository, follow the Google Go style guides in full:

- Go Style Guide: https://google.github.io/styleguide/go/
- Style Decisions: https://google.github.io/styleguide/go/decisions
- Best Practices: https://google.github.io/styleguide/go/best-practices

These references are the source of truth. If this repo guidance and Google guidance differ, Google Go guidance wins unless a repo-specific exception is explicitly documented.

Required implementation rules for day-to-day work:

- Formatting:
  - run `gofmt` on changed Go files before commit
  - keep imports gofmt-compatible and clean (no unused imports)
  - keep package layout simple and cohesive; one package responsibility per directory
- Naming and declarations:
  - use idiomatic Go names (`MixedCaps`, no underscores in identifiers)
  - use clear package names; avoid stuttered APIs (`foo.FooType`)
  - keep exported identifiers documented with comments that start with the identifier name
  - use consistent initialisms (`ID`, `HTTP`, `URL`, `JSON`)
- Error handling:
  - check and handle every returned error; do not ignore errors silently
  - wrap errors with context using `%w` when propagating (`fmt.Errorf("...: %w", err)`)
  - keep error strings lowercase and without trailing punctuation
  - reserve `panic` for truly unrecoverable programmer/runtime faults, not normal control flow
- Context and API shape:
  - pass `context.Context` as the first argument (`ctx context.Context`) when cancellation/deadlines are relevant
  - do not store `Context` in structs
  - accept interfaces where useful, but return concrete types from constructors/functions
  - design zero-value-safe types where practical
- Concurrency:
  - make goroutine ownership/lifecycle explicit and prevent leaks
  - close channels only from the sender side; document channel ownership
  - protect shared mutable state with synchronization and avoid data races
- Tests:
  - prefer table-driven tests for multi-case logic
  - make failures actionable (`got` vs `want` in error messages)
  - keep tests deterministic; avoid time-based flakiness
  - include coverage for success paths, error paths, and edge cases
- Documentation:
  - add package comments for non-trivial packages
  - document exported behavior, invariants, side effects, and concurrency expectations
  - keep examples/docs synchronized with current CLI/API behavior

## Architecture Standard (Hexagonal / Ports And Adapters)

VERY IMPORTANT (MANDATORY):

- all new features and refactors must follow hexagonal architecture (ports and adapters)
- architecture violations are blocking issues and must be fixed before merge

Authoritative resources (required reading for contributors):

- Alistair Cockburn, original Hexagonal Architecture article:
  - https://alistair.cockburn.us/hexagonal-architecture
- AWS Prescriptive Guidance, Hexagonal architecture pattern:
  - https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/hexagonal-architecture.html
- Robert C. Martin, Clean Architecture and dependency rule context:
  - https://blog.cleancoder.com/uncle-bob/2011/11/22/Clean-Architecture.html

Concept model:

- center (inside): business core
  - domain model + business rules (`internal/domain`)
  - application/use-case orchestration (`internal/app`)
- boundary: ports
  - inbound ports: use-case interfaces exposed by the core
  - outbound ports: interfaces the core needs for side effects (exchange API, secret store, DB, clock, etc.)
- edge (outside): adapters
  - primary/inbound adapters drive use cases (CLI, HTTP handlers, jobs)
  - secondary/outbound adapters implement outbound ports (WhiteBIT client, keychain, persistence, messaging)

Dependency rule (non-negotiable):

- dependencies point inward only
- `internal/domain` must not import adapter/infrastructure packages
- `internal/app` may depend on `domain` and port interfaces, never concrete infrastructure
- adapters depend on core ports and translate external concerns to/from core models
- `cmd/wbcli` is the composition root that wires concrete adapters into use cases

Implementation contract by layer:

- Domain (`internal/domain`)
  - pure business logic, entities/value objects, invariants, validation rules
  - no transport, IO, DB, CLI flag parsing, or framework code
- Application (`internal/app`)
  - use-cases, orchestration, transaction boundaries, policy sequencing
  - depends on domain + abstract ports only
  - defines request/response DTOs suitable for adapters
- Ports
  - declare explicit interfaces at the core boundary
  - keep interfaces small and behavior-focused (avoid god interfaces)
  - include context and error contracts where relevant
- Adapters (`internal/adapters/*`)
  - map external formats/protocols to core DTOs and back
  - isolate exchange-specific/auth-specific/storage-specific details
  - do not leak transport models into domain
- Composition root (`cmd/wbcli`)
  - instantiate adapters
  - inject them into application services
  - select runtime configuration/profile/environment

Flow for a command/request:

1. inbound adapter parses input and performs syntax-level validation
2. adapter calls an application use-case through an inbound port
3. use-case applies domain rules and invokes outbound ports
4. outbound adapters execute side effects (API/storage/etc.)
5. use-case returns stable result model
6. inbound adapter renders output (`table`/`json`) without business logic leakage

Testing strategy tied to architecture:

- domain tests: pure unit tests, no mocks for infrastructure
- application tests: use fakes/stubs for outbound ports, validate orchestration and error mapping
- adapter tests: integration-focused contract tests against real protocol boundaries where feasible
- end-to-end smoke tests: CLI path through composition root

Anti-patterns (prohibited):

- business rules inside CLI handlers or transport structs
- domain package importing infrastructure libraries
- direct SDK/DB calls from use-case layer without outbound port
- shared mutable globals crossing layers
- adapter-specific error/data types leaking into domain APIs

PR architecture checklist (mandatory):

- does this change keep dependencies pointing inward?
- are side effects behind outbound ports?
- is business logic placed in `domain`/`app` instead of adapters?
- can the use-case run with test doubles for infrastructure?
- are adapters thin translators rather than decision engines?

## Documentation Update Policy

When behavior changes, update docs in the same change set:

- user-facing usage docs
- architecture/design notes when boundaries change
- runbook/ops docs when operational steps change

If no doc impact, change notes must explicitly state: `Docs impact: none`.

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
3. Define branch protection for trunk mode: direct pushes to `main` allowed for the 2-developer team, CI required on `main`.
4. Enable CI checks: lint, test, build, secret scan.
5. Start triage: label current work with IDs, priorities, owners, due dates.
6. Set WIP rule (`max 2`) and stale reminder automation.
7. Schedule weekly (30-45 min) and monthly (60 min) operating reviews.

## Common Failure Modes And Prevention Rules

- Work starts without clear acceptance criteria:
  - prevention: enforce DoR before `In Progress`
- Too many parallel tasks, little completion:
  - prevention: strict WIP limit and weekly aging review
- Changes land with hidden risk:
  - prevention: mandatory risk section + rollback plan
- Docs drift from behavior:
  - prevention: docs update required in DoD
- Repeated incidents with no learning:
  - prevention: required postmortem actions tracked as tickets
- Release pain accumulates:
  - prevention: fixed cadence + changelog discipline

# PROJ-2026-013: Add release packaging and install documentation for CLI

ID: PROJ-2026-013
Title: Add release packaging and install documentation for CLI
Priority: P2
Status: Backlog
Owner: chewbaccalol
Due Date: 2026-03-17
Created: 2026-02-26
Updated: 2026-02-26
Links: [README](../../README.md), [AGENTS](../../AGENTS.md)

Problem:
Without release automation and install docs, onboarding and repeatable delivery remain manual and error-prone.

Outcome:
CLI can be packaged and installed consistently, with release notes tied to tickets and rollback instructions documented.

Acceptance Criteria:
- [ ] Build artifacts are produced for target platforms (Linux/macOS at minimum).
- [ ] Install paths are documented (`go install` and binary download flow).
- [ ] Changelog `Unreleased` process is defined and linked to ticket IDs.
- [ ] Release checklist includes verification and rollback steps.
- [ ] A dry-run release pipeline is executed successfully.

Risks:
- Packaging misconfiguration can ship broken binaries.
- Missing checksum/signature verification reduces distribution trust.

Rollout Plan:
1. Add release workflow/scripts.
2. Add install + upgrade documentation.
3. Validate dry-run release and update changelog process.

Rollback Plan:
1. Stop automated publish and keep manual artifact verification only.
2. Revert release workflow changes while preserving changelog updates.

Status Notes:
- 2026-02-26: Created in Backlog.
- 2026-02-26: Scheduled for end of full-scale phase after command set stabilizes.

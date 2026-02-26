# PROJ-2026-015: Implement credential access controls session unlock and key rotation workflow

ID: PROJ-2026-015
Title: Implement credential access controls session unlock and key rotation workflow
Priority: P1
Status: Backlog
Owner: chewbaccalol
Due Date: 2026-03-07
Created: 2026-02-26
Updated: 2026-02-26
Links: [CLI Design](../../docs/cli-design.md), [WhiteBIT Integration](../../docs/whitebit-integration.md), [WhiteBIT HTTP Auth](https://docs.whitebit.com/private/http-auth)

Problem:
Even with encrypted storage, runtime secret access can leak through logs, long-lived sessions, or poor rotation hygiene.

Outcome:
Secret access is time-scoped and auditable, and API key rotation/revocation workflows are first-class CLI operations.

Acceptance Criteria:
- [ ] CLI supports session unlock TTL for credential reads (for example 15 minutes) with explicit lock command.
- [ ] Secrets are held in memory only as long as needed for signing and are zeroed where practical.
- [ ] Logs and error outputs redact API key, payload, and signature values.
- [ ] `whitbit keys rotate` workflow documents and validates safe cutover from old key to new key.
- [ ] `whitbit keys revoke` flow removes local credential and prints exchange-side revoke checklist.
- [ ] Tests cover TTL expiry, redaction guarantees, and rotation success/failure paths.

Risks:
- Overly strict unlock policies can hurt automation workflows.
- Rotation flows can cause downtime if cutover sequencing is wrong.

Rollout Plan:
1. Add runtime credential session manager and TTL enforcement.
2. Add redaction middleware for command output and logs.
3. Implement rotation/revoke commands with dry-run support.

Rollback Plan:
1. Disable session cache and force per-command credential prompt.
2. Keep existing key commands while disabling rotation automation until stabilized.

Status Notes:
- 2026-02-26: Created in Backlog.
- 2026-02-26: Added as follow-up hardening for runtime access and lifecycle security.

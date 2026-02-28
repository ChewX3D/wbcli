# Auth Subcommands Specification

Status: current snapshot  
Last updated: 2026-02-28

## Scope

`wbcli auth` is single-session only: user is either logged in or logged out.  
There is no profile model in current command behavior.

Removed commands:

- `wbcli auth set`
- `wbcli auth use`
- `wbcli auth list`
- `wbcli auth current`

## Commands

### `wbcli auth login`

Purpose: store WhiteBIT credentials in secure backend and mark session as logged in.

Input contract:

- reads credentials from stdin only
- accepts exactly two non-empty logical lines:
  - line 1: `api_key`
  - line 2: `api_secret`
- maximum payload size: `16 KiB`
- no `--api-key`, no `--api-secret`, no `--profile`

Behavior:

- overwrites existing stored credential/session by default
- writes secrets to secure backend only
- updates metadata session state

Success output format:

```text
logged_in=true backend=<backend> api_key=<masked_hint> saved_at=<RFC3339 timestamp>
```

### `wbcli auth logout`

Purpose: clear stored credential and auth session metadata.

Behavior:

- idempotent; succeeds even when user is already logged out

Success output format:

```text
logged_out=true
```

### `wbcli auth status`

Purpose: show current auth state using safe metadata only.

Output formats:

```text
logged_in=false
```

or

```text
logged_in=true backend=<backend> api_key=<masked_hint> updated_at=<RFC3339 timestamp>
```

### `wbcli auth test`

Purpose: validate stored credentials against WhiteBIT via `AuthProbe`.

Current behavior:

- returns not-implemented error while WhiteBIT probe is not wired
- once probe is wired, success output is:

```text
auth test passed
```

## Storage And Security Boundaries

- secret credentials are stored in `os-keychain` backend only
- non-secret metadata is stored in `~/.wbcli/config.yaml`
- config file permission target on macOS/Linux: `0600`
- config must not contain API secret or raw credential material
- command output and error mapping must not expose secret values, payloads, or signatures

## Error Contract (User-Facing)

Actionable auth errors include:

- stdin payload missing/invalid/too large
- not logged in
- keychain unavailable
- keychain permission denied
- auth test not implemented yet (until WhiteBIT client/probe is ready)

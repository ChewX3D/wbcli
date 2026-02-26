# Tickets System

This repository uses markdown tickets with lightweight automation scripts.

## Structure

- `tickets/backlog`: queued work not yet ready
- `tickets/ready`: refined work ready to start
- `tickets/in-progress`: currently being implemented
- `tickets/review`: awaiting review/merge
- `tickets/blocked`: paused with explicit blocker
- `tickets/done`: completed recently
- `tickets/archived`: old completed work
- `tickets/board.md`: generated status board
- `tickets/templates`: reusable templates

## ID Convention

- Format: `PROJ-YYYY-NNN` (example: `PROJ-2026-001`)
- Default prefix is `PROJ` and can be overridden with `--prefix` or `TICKET_PREFIX`.

## Commands

Create ticket:

```bash
./scripts/tickets/new.sh --title "Implement secure key storage" --priority P1 --owner chewbaccalol --status Ready
```

Move ticket:

```bash
./scripts/tickets/move.sh PROJ-2026-001 "In Progress" "Implementation started."
```

Rebuild board:

```bash
./scripts/tickets/board.sh
```

## Rules

- Every non-trivial change starts with a ticket.
- Ticket must include acceptance criteria before moving to `Ready`.
- Max 2 tickets in `In Progress` per owner.
- `Done` tickets older than 30 days move to `Archived`.

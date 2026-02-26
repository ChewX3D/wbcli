#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TICKETS_DIR="$ROOT_DIR/tickets"
BOARD_SCRIPT="$ROOT_DIR/scripts/tickets/board.sh"

usage() {
  cat <<'EOF'
Usage:
  ./scripts/tickets/new.sh --title "Short title" [options]

Options:
  --priority P0|P1|P2|P3     Default: P2
  --owner NAME               Default: unassigned
  --due YYYY-MM-DD|None      Default: None
  --status STATUS            Default: Backlog
  --prefix PREFIX            Default: PROJ (or $TICKET_PREFIX)
  -h, --help

Statuses:
  Backlog, Ready, In Progress, Review, Blocked, Done, Archived
EOF
}

normalize_status() {
  local raw="${1:-}"
  local norm
  norm="$(echo "$raw" | tr '[:upper:]' '[:lower:]' | tr -d '[:space:]' | tr '_' '-')"
  case "$norm" in
    backlog) echo "backlog|Backlog" ;;
    ready) echo "ready|Ready" ;;
    in-progress|inprogress) echo "in-progress|In Progress" ;;
    review) echo "review|Review" ;;
    blocked) echo "blocked|Blocked" ;;
    done) echo "done|Done" ;;
    archived) echo "archived|Archived" ;;
    *)
      echo "Invalid status: $raw" >&2
      exit 1
      ;;
  esac
}

TITLE=""
PRIORITY="P2"
OWNER="unassigned"
DUE_DATE="None"
STATUS="Backlog"
PREFIX="${TICKET_PREFIX:-PROJ}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --title) TITLE="${2:-}"; shift 2 ;;
    --priority) PRIORITY="${2:-}"; shift 2 ;;
    --owner) OWNER="${2:-}"; shift 2 ;;
    --due) DUE_DATE="${2:-}"; shift 2 ;;
    --status) STATUS="${2:-}"; shift 2 ;;
    --prefix) PREFIX="${2:-}"; shift 2 ;;
    -h|--help) usage; exit 0 ;;
    *)
      echo "Unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$TITLE" ]]; then
  echo "--title is required." >&2
  usage
  exit 1
fi

case "$PRIORITY" in
  P0|P1|P2|P3) ;;
  *)
    echo "Invalid priority: $PRIORITY" >&2
    exit 1
    ;;
esac

if [[ "$DUE_DATE" != "None" ]] && ! [[ "$DUE_DATE" =~ ^[0-9]{4}-[0-9]{2}-[0-9]{2}$ ]]; then
  echo "Due date must be YYYY-MM-DD or None." >&2
  exit 1
fi

if ! [[ "$PREFIX" =~ ^[A-Z][A-Z0-9]+$ ]]; then
  echo "Prefix must match ^[A-Z][A-Z0-9]+$ (example: CORE)." >&2
  exit 1
fi

status_info="$(normalize_status "$STATUS")"
status_dir="${status_info%%|*}"
status_label="${status_info##*|}"

year="$(date +%Y)"
max_seq="$(
  find "$TICKETS_DIR" -type f -name "${PREFIX}-${year}-*.md" \
    | sed -E "s#.*/${PREFIX}-${year}-([0-9]{3}).*#\1#" \
    | sort -n \
    | tail -1
)"
max_seq="${max_seq:-000}"
next_seq="$(printf "%03d" "$((10#$max_seq + 1))")"
ticket_id="${PREFIX}-${year}-${next_seq}"

slug="$(
  echo "$TITLE" \
    | tr '[:upper:]' '[:lower:]' \
    | sed -E 's/[^a-z0-9]+/-/g; s/^-+//; s/-+$//; s/-+/-/g'
)"
slug="${slug:-ticket}"

file_path="$TICKETS_DIR/$status_dir/${ticket_id}-${slug}.md"
today="$(date +%F)"

cat > "$file_path" <<EOF
# ${ticket_id}: ${TITLE}

ID: ${ticket_id}
Title: ${TITLE}
Priority: ${PRIORITY}
Status: ${status_label}
Owner: ${OWNER}
Due Date: ${DUE_DATE}
Created: ${today}
Updated: ${today}
Links: []

Problem:

Outcome:

Acceptance Criteria:
- [ ]

Risks:

Rollout Plan:

Rollback Plan:

Status Notes:
- ${today}: Created in ${status_label}.
EOF

"$BOARD_SCRIPT" >/dev/null
echo "Created $ticket_id"
echo "File: $file_path"

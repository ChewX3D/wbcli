#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TICKETS_DIR="$ROOT_DIR/tickets"
BOARD_SCRIPT="$ROOT_DIR/scripts/tickets/board.sh"

usage() {
  cat <<'EOF'
Usage:
  ./scripts/tickets/move.sh TICKET_ID TARGET_STATUS [NOTE]

Example:
  ./scripts/tickets/move.sh PROJ-2026-001 "In Progress" "Started implementation."
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

if [[ $# -lt 2 ]]; then
  usage
  exit 1
fi

ticket_id="$1"
target_status="$2"
note="${3:-Moved to ${target_status}.}"

status_info="$(normalize_status "$target_status")"
target_dir="${status_info%%|*}"
target_label="${status_info##*|}"

source_file="$(
  find "$TICKETS_DIR" -mindepth 2 -maxdepth 2 -type f -name "${ticket_id}-*.md" \
    | head -n1
)"

if [[ -z "$source_file" ]]; then
  echo "Ticket not found: $ticket_id" >&2
  exit 1
fi

target_file="$TICKETS_DIR/$target_dir/$(basename "$source_file")"
if [[ "$source_file" != "$target_file" ]]; then
  mv "$source_file" "$target_file"
fi

today="$(date +%F)"

if grep -q '^Status:' "$target_file"; then
  sed -i -E "s/^Status:.*/Status: ${target_label}/" "$target_file"
else
  printf '\nStatus: %s\n' "$target_label" >> "$target_file"
fi

if grep -q '^Updated:' "$target_file"; then
  sed -i -E "s/^Updated:.*/Updated: ${today}/" "$target_file"
else
  printf 'Updated: %s\n' "$today" >> "$target_file"
fi

if ! grep -q '^Status Notes:' "$target_file"; then
  printf '\nStatus Notes:\n' >> "$target_file"
fi
printf -- '- %s: %s\n' "$today" "$note" >> "$target_file"

"$BOARD_SCRIPT" >/dev/null
echo "Moved $ticket_id to $target_label"
echo "File: $target_file"

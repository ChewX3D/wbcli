#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TICKETS_DIR="$ROOT_DIR/tickets"
BOARD_FILE="$TICKETS_DIR/board.md"

get_field() {
  local key="$1"
  local file="$2"
  local value
  value="$(sed -n "s/^${key}:[[:space:]]*//p" "$file" | head -n1)"
  if [[ -n "$value" ]]; then
    echo "$value"
  else
    echo "n/a"
  fi
}

print_status_section() {
  local dir="$1"
  local label="$2"
  local files=()
  local file

  echo "## ${label}"
  mapfile -t files < <(find "$TICKETS_DIR/$dir" -maxdepth 1 -type f -name '*.md' | sort)

  if [[ "${#files[@]}" -eq 0 ]]; then
    echo "- _(empty)_"
    echo
    return
  fi

  for file in "${files[@]}"; do
    local id title priority owner due rel
    id="$(get_field "ID" "$file")"
    title="$(get_field "Title" "$file")"
    priority="$(get_field "Priority" "$file")"
    owner="$(get_field "Owner" "$file")"
    due="$(get_field "Due Date" "$file")"
    rel="./$dir/$(basename "$file")"
    echo "- [${id}](${rel}) | ${priority} | Owner: ${owner} | Due: ${due} | ${title}"
  done
  echo
}

{
  echo "# Ticket Board"
  echo
  echo "_Last updated: $(date -Iseconds)_"
  echo
  print_status_section "backlog" "Backlog"
  print_status_section "ready" "Ready"
  print_status_section "in-progress" "In Progress"
  print_status_section "review" "Review"
  print_status_section "blocked" "Blocked"
  print_status_section "done" "Done"
  print_status_section "archived" "Archived"
} > "$BOARD_FILE"

echo "Updated $BOARD_FILE"

#!/usr/bin/env bash
# Copies the project state to the clipboard (macOS).
# Usage: scripts/context.sh [extra files...]
set -u
cd "$(git rev-parse --show-toplevel)" || exit 1

section() { printf '\n===== %s =====\n' "$1"; }
show_file() { section "$1"; cat "$1" 2>/dev/null || echo "(missing)"; }

{
  show_file docs/PROGRESS.md
  show_file docs/DECISIONS.md
  section "git log (last 15)"; git log --oneline -15 2>&1
  section "git status";        git status --short 2>&1
  section "file list";         git ls-files --cached --others --exclude-standard
  show_file go.mod
  for f in $(find internal -path '*/port/*.go' 2>/dev/null | sort); do show_file "$f"; done
  for f in "$@"; do show_file "$f"; done
  section "go test"; go test -count=1 ./... 2>&1
} | pbcopy

echo "Copied to clipboard."

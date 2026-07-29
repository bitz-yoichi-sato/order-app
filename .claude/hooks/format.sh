#!/usr/bin/env bash
set -euo pipefail

input="$(cat)"

if command -v jq >/dev/null 2>&1; then
  file_path="$(printf '%s' "$input" | jq -r '.tool_input.file_path // empty')"
else
  file_path="$(printf '%s' "$input" | sed -n 's/.*"file_path"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1)"
fi

if [ -z "${file_path:-}" ]; then
  exit 0
fi

case "$file_path" in
  *.go)
    gofmt -w "$file_path"
    if command -v golangci-lint >/dev/null 2>&1; then
      if ! golangci-lint run "$file_path"; then
        exit 2
      fi
    fi
    ;;
  *.ts | *.tsx)
    if command -v prettier >/dev/null 2>&1; then
      npx prettier --write "$file_path"
    fi
    ;;
  *)
    exit 0
    ;;
esac

exit 0

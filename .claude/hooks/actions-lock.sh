#!/usr/bin/env bash
# PostToolUse hook: after Claude edits a workflow, run `gh actions-lock` so that
# .github/workflows/actions.lock follows the `uses:` of the workflows (the CI verifies the lockfile).
set -euo pipefail

file=$(jq -r '.tool_input.file_path // empty')
case "${file}" in
  */.github/workflows/*.yml | */.github/workflows/*.yaml) ;;
  *) exit 0 ;;
esac

cd "${CLAUDE_PROJECT_DIR:-.}"
# Exit status 2 shows the output to Claude, so that it can fix the workflow.
gh actions-lock >&2 || exit 2

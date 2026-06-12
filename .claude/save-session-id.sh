#!/usr/bin/env bash
set -euo pipefail
session_id=$(jq -r '.session_id // empty')
[[ -n "$session_id" ]] && echo "$session_id" > "$(dirname "$0")/last_session_id"

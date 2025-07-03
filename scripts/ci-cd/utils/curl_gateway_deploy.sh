#!/bin/bash
set -euo pipefail

# --- Input
GATEWAY_URL="${GATEWAY_URL:-}"
GATEWAY_USER="${GATEWAY_USER:-}"
GATEWAY_PASSWORD="${GATEWAY_PASSWORD:-}"
CI_PROJECT_ID="${CI_PROJECT_ID:-}"
CI_PROJECT_PATH="${CI_PROJECT_PATH:-}"

# --- Validation
if [[ -z "$GATEWAY_URL" || -z "$GATEWAY_USER" || -z "$GATEWAY_PASSWORD" || -z "$CI_PROJECT_ID" || -z "$CI_PROJECT_PATH" ]]; then
  echo "[ERROR] Missing required environment variables."
  exit 1
fi

# --- jq check
if ! command -v jq >/dev/null; then
  echo "[INFO] Installing jq…"
  apt-get update -qq && apt-get install -yqq jq
fi

# --- Get Token
echo "[INFO] Requesting gateway token…"
TOKEN=$(curl --silent --fail -X POST "$GATEWAY_URL/api/v1/auth/token" \
  -H "Content-Type: application/json" \
  --data "{\"username\":\"$GATEWAY_USER\",\"password\":\"$GATEWAY_PASSWORD\"}" | jq -r '.token')

if [[ -z "$TOKEN" || "$TOKEN" == "null" ]]; then
  echo "[ERROR] Unable to retrieve token."
  exit 1
fi

# --- Read version.json
echo "[INFO] Reading version.json…"
ENVS=$(jq -r '.envs' version.json | tr -d '"')
REGIONS=$(jq -r '.regions // empty' version.json | tr -d '"')

# --- Compose payload
BODY=$(jq -n \
  --arg envs "$ENVS" \
  --arg regions "$REGIONS" \
  --arg git_id "$CI_PROJECT_ID" \
  --arg gitlab_path "$CI_PROJECT_PATH" \
  '{envs:$envs, git_id:$git_id, gitlab_path:$gitlab_path, regions:$regions}')

echo "[INFO] Sending sync request to $GATEWAY_URL/api/v1/argocd/refresh-sync…"

RESPONSE=$(curl --silent -X POST "$GATEWAY_URL/api/v1/argocd/refresh-sync" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  --data "$BODY" \
  --max-time 360)

echo "[DEBUG] Response:"
echo "$RESPONSE" | jq .

STATUS=$(echo "$RESPONSE" | jq -r '.status // "UNKNOWN"')

if [[ "$STATUS" != "SUCCESS" ]]; then
  echo "[ERROR] Sync failed: $STATUS"
  exit 1
fi

echo "[SUCCESS] Sync completed."

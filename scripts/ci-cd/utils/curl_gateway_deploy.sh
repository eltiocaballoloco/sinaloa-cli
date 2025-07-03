#!/bin/bash
set -euo pipefail

############################
# 1) Input argument (env)  #
############################
ENV="$1"                       # es. "prod" oppure "perf"

if [[ -z "$ENV" ]]; then
  echo "[ERROR] No environment passed."
  exit 1
fi

############################
# 2) Mandatory variables   #
############################
: "${GATEWAY_URL:?Missing GATEWAY_URL}"
: "${GATEWAY_USER:?Missing GATEWAY_USER}"
: "${GATEWAY_PASSWORD:?Missing GATEWAY_PASSWORD}"
: "${CI_PROJECT_ID:?Missing CI_PROJECT_ID}"
: "${CI_PROJECT_PATH:?Missing CI_PROJECT_PATH}"

############################
# 3) Ensure jq is present  #
############################
if ! command -v jq >/dev/null; then
  echo "[INFO] Installing jq…"
  apt-get update -qq && apt-get install -yqq jq
fi

############################
# 4) Get Gateway token     #
############################
echo "[INFO] Requesting gateway token…"
TOKEN=$(curl --silent --fail \
  -X POST "$GATEWAY_URL/api/v1/auth/token" \
  -H "Content-Type: application/json" \
  --data "{\"username\":\"$GATEWAY_USER\",\"password\":\"$GATEWAY_PASSWORD\"}" |
  jq -r '.token')

if [[ -z "$TOKEN" || "$TOKEN" == "null" ]]; then
  echo "[ERROR] Unable to retrieve token."
  exit 1
fi

#############################################
# 5) Read region for the specific ENV       #
#############################################
REGION=$(jq -r --arg env "$ENV" '.regions[$env] // ""' version.json | xargs)

if [[ -z "$REGION" ]]; then
  echo "[ERROR] Region not defined for env '$ENV' in version.json"
  exit 1
fi

echo "[INFO] Env: $ENV  ->  Region string: '$REGION'"

############################
# 6) Compose JSON payload  #
############################
BODY=$(jq -n \
  --arg envs "$ENV" \
  --arg regions "$REGION" \
  --arg git_id "$CI_PROJECT_ID" \
  --arg gitlab_path "$CI_PROJECT_PATH" \
  '{envs:$envs, git_id:$git_id, gitlab_path:$gitlab_path, regions:$regions}')

############################
# 7) Call refresh‑sync     #
############################
echo "[INFO] Sending sync request…"
RESPONSE=$(curl --silent --fail \
  -X POST "$GATEWAY_URL/api/v1/argocd/refresh-sync" \
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

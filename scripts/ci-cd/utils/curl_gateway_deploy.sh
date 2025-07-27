#!/bin/bash
set -euo pipefail


############################
# 1) Input argument (env)  #
############################
ENV="$1"  # eg. "prod" or "perf"

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
echo "[INFO] Requesting gateway token..."
TOKEN=$(curl --silent --show-error --fail \
  -X POST "$GATEWAY_URL/api/v1/auth/token" \
  -H "Content-Type: application/json" \
  --data "{\"username\":\"$GATEWAY_USER\",\"password\":\"$GATEWAY_PASSWORD\"}" \
  || { echo "[ERROR] Failed to get token. Check GATEWAY_URL or credentials."; exit 1; })

# Extract the token value from the json
TOKEN=$(echo "$TOKEN" | jq -r '.token')

# Check token
if [[ -z "$TOKEN" || "$TOKEN" == "null" ]]; then
  echo "[ERROR] Unable to retrieve token."
  exit 1
fi


#############################################
# 5) Read region for the specific ENV       #
#############################################
REGION=$(jq -r --arg env "$ENV" '.regions[$env] // ""' version.json | xargs)

# If we find in region "backup,main",
# set the correct region to sync
if [[ "$REGION" == *"backup,main"* ]]; then
  if [[ "$ENV" == "prod" ]]; then
    REGION="main"
  elif [[ "$ENV" == "prod-backup" ]]; then
    REGION="backup"
  fi
fi

# Set the env
if [[ "$ENV" == "prod" || "$ENV" == "prod-backup" ]]; then
  ENV_ARG="prod"
else
  ENV_ARG="$ENV"
fi

echo "[INFO] Env: '$ENV'  ->  Regions: '$REGION'"


############################
# 6) Compose JSON payload  #
############################
BODY=$(jq -n \
  --arg envs "$ENV_ARG" \
  --arg regions "$REGION" \
  --arg git_id "$CI_PROJECT_ID" \
  --arg gitlab_path "$CI_PROJECT_PATH" \
  '{envs:$envs, git_id:$git_id, gitlab_path:$gitlab_path, regions:$regions}')


############################
# 7) Call refresh‑sync     #
############################
echo "[INFO] Sending the sync request to deploy the apps, please wait..."

RESPONSE=$(curl --silent --show-error --location \
  --request POST "$GATEWAY_URL/api/v1/argocd/refresh-sync" \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data "$BODY" \
  --connect-timeout 30 \
  --max-time 360
) || {
  CODE=$?
  echo "[ERROR] Curl request failed (exit code $CODE)"
  echo "[ERROR] Payload was:"
  echo "$BODY" | jq .
  exit 1
}

# Verify json response
if ! echo "$RESPONSE" | jq . >/dev/null 2>&1; then
  echo "[ERROR] Invalid JSON received from API."
  echo "$RESPONSE"
  exit 1
fi

# Parse status
STATUS=$(echo "$RESPONSE" | jq -r '.status // "UNKNOWN"')
MESSAGE=$(echo "$RESPONSE" | jq -r '.message // "No message provided."')

if [[ "$STATUS" != "SUCCESS" ]]; then
  echo "[ERROR] Sync failed: $STATUS"
  echo "[ERROR] Message: $MESSAGE"
  exit 1
fi

echo "[INFO] Sync completed."
echo "[INFO] $MESSAGE"

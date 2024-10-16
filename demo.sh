#!/usr/bin/env bash

# session manager demo script

# Requires: bash curl jq

set -e
set -o pipefail

test -n "$AB_CLIENT_ID" || (echo "AB_CLIENT_ID is not set"; exit 1)
test -n "$AB_CLIENT_SECRET" || (echo "AB_CLIENT_SECRET is not set"; exit 1)
test -n "$AB_NAMESPACE" || (echo "AB_NAMESPACE is not set"; exit 1)

if [ -z "$GRPC_SERVER_URL" ] && [ -z "$EXTEND_APP_NAME" ]; then
  echo "GRPC_SERVER_URL or EXTEND_APP_NAME is not set"
  exit 1
fi

DEMO_PREFIX='sm_grpc_demo'

api_curl()
{
  curl -s -D api_curl_http_header.out -o api_curl_http_response.out -w '%{http_code}' "$@" > api_curl_http_code.out
  echo >> api_curl_http_response.out
  cat api_curl_http_response.out
}

clean_up()
{
  echo Deleting player ...

  api_curl -X DELETE "${AB_BASE_URL}/iam/v3/admin/namespaces/$AB_NAMESPACE/users/$USER_ID/information" -H "Authorization: Bearer $ACCESS_TOKEN"

  echo Deleting session configuration ...

  api_curl -X DELETE -s "${AB_BASE_URL}/session/v1/admin/namespaces/$AB_NAMESPACE/configurations/$DEMO_PREFIX" -H "Authorization: Bearer $ACCESS_TOKEN" -H 'Content-Type: application/json' >/dev/null
}

trap clean_up EXIT

echo Logging in client ...

ACCESS_TOKEN="$(api_curl -s ${AB_BASE_URL}/iam/v3/oauth/token -H 'Content-Type: application/x-www-form-urlencoded' -u "$AB_CLIENT_ID:$AB_CLIENT_SECRET" -d "grant_type=client_credentials" | jq --raw-output .access_token)"

if [ "$(cat api_curl_http_code.out)" -ge "400" ]; then
  cat api_curl_http_response.out
  exit 1
fi

if [ -n "$GRPC_SERVER_URL" ]; then
  echo Registering session manager $GRPC_SERVER_URL in the session configuration ...

  api_curl -X POST -s "${AB_BASE_URL}/session/v1/admin/namespaces/$AB_NAMESPACE/configuration/" -H "Authorization: Bearer $ACCESS_TOKEN" -H 'Content-Type: application/json' -d "{\"name\":\"${DEMO_PREFIX}\",\"joinability\":\"OPEN\",\"maxPlayers\":2,\"minPlayers\":1,\"autoJoin\":true,\"type\":\"NONE\",\"dsSource\":\"custom\",\"dsManualSetReady\":false,\"requestedRegions\":[\"us-west-2\"],\"grpcSessionConfig\":{\"customURL\":\"${GRPC_SERVER_URL}\",\"functionFlag\":7}}" >/dev/null

  if [ "$(cat api_curl_http_code.out)" -ge "400" ]; then
    exit 1
  fi
elif [ -n "$EXTEND_APP_NAME" ]; then
  echo Registering session manager $EXTEND_APP_NAME in the session configuration ...

  api_curl -X POST -s "${AB_BASE_URL}/session/v1/admin/namespaces/$AB_NAMESPACE/configuration" -H "Authorization: Bearer $ACCESS_TOKEN" -H 'Content-Type: application/json' -d "{\"name\":\"${DEMO_PREFIX}\",\"joinability\":\"OPEN\",\"maxPlayers\":2,\"minPlayers\":1,\"autoJoin\":true,\"type\":\"NONE\",\"dsSource\":\"custom\",\"dsManualSetReady\":false,\"requestedRegions\":[\"us-west-2\"],\"grpcSessionConfig\":{\"appName\":\"${EXTEND_APP_NAME}\",\"functionFlag\":7}}" >/dev/null

  if [ "$(cat api_curl_http_code.out)" -ge "400" ]; then
    exit 1
  fi
else
  echo "GRPC_SERVER_URL or EXTEND_APP_NAME is not set"
  exit 1
fi

echo Creating PLAYER ...

USER_ID="$(api_curl -s "${AB_BASE_URL}/iam/v4/public/namespaces/$AB_NAMESPACE/users" -H "Authorization: Bearer $ACCESS_TOKEN" -H 'Content-Type: application/json' -d "{\"authType\":\"EMAILPASSWD\",\"country\":\"ID\",\"dateOfBirth\":\"1995-01-10\",\"displayName\":\"Session Manager gRPC Player\",\"uniqueDisplayName\":\"Session Manager gRPC Player\",\"emailAddress\":\"${DEMO_PREFIX}_player@test.com\",\"password\":\"GFPPlmdb2-\",\"username\":\"${DEMO_PREFIX}_player\"}" | jq --raw-output .userId)"

if [ "$(cat api_curl_http_code.out)" -ge "400" ]; then
  cat api_curl_http_response.out
  exit 1
fi

echo Test OnSessionCreated ...

SESSION_ID="$(api_curl -s "${AB_BASE_URL}/session/v1/public/namespaces/$AB_NAMESPACE/gamesession" -H "Authorization: Bearer $ACCESS_TOKEN" -H 'Content-Type: application/json' -d "{\"configurationName\":\"${DEMO_PREFIX}\",\"teams\":[{\"parties\":[{\"partyID\":\"\",\"userIDs\":[\"${USER_ID}\"]}],\"userIDs\":[\"${USER_ID}\"]}]}" | jq --raw-output .id)"

echo

if [ "$(cat api_curl_http_code.out)" -ge "400" ]; then
  exit 1
fi

echo Test OnSessionUpdated ...

api_curl -X PATCH -s "${AB_BASE_URL}/session/v1/public/namespaces/$AB_NAMESPACE/gamesessions/$SESSION_ID" -H "Authorization: Bearer $ACCESS_TOKEN" -H 'Content-Type: application/json' -d "{\"attributes\":{\"SAMPLE\":\"value from GRPC server updated\"},\"version\":1}"

echo

if [ "$(cat api_curl_http_code.out)" -ge "400" ]; then
  exit 1
fi

echo Test OnSessionDeleted ...

api_curl -X DELETE -s "${AB_BASE_URL}/session/v1/public/namespaces/$AB_NAMESPACE/gamesessions/$SESSION_ID" -H "Authorization: Bearer $ACCESS_TOKEN" -H 'Content-Type: application/json'
echo

if [ "$(cat api_curl_http_code.out)" -ge "400" ]; then
  exit 1
fi

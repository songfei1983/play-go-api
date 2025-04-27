#!/bin/zsh

# Exit on error
set -e

# Generate random username
generate_username() {
    echo "user_$(date +%s)_${RANDOM}"
}

# Base URL for the API
BASE_URL="http://localhost:8080"

# Wait for service to be ready
echo "Waiting for service to be ready..."
max_attempts=30
attempt=0

# curl flags explanation:
# -s, --silent: Don't show progress meter or error messages
# -f, --fail: Return error on HTTP errors (non 2xx responses)
# -i: Include HTTP headers in the output
while true; do
    HEALTH_CHECK=$(curl -s -f -i "${BASE_URL}/health" 2>&1)
    STATUS=$?

    if [ $STATUS -eq 0 ]; then
        if echo "$HEALTH_CHECK" | grep -q "200 OK"; then
            echo "Service is healthy!"
            break
        fi
    fi

    attempt=$((attempt+1))
    if [ $attempt -eq $max_attempts ]; then
        echo "Service failed to become ready after $max_attempts attempts"
        echo "Last health check response:"
        echo "$HEALTH_CHECK"
        exit 1
    fi
    printf '.'
    sleep 1
done
echo "Service is ready!"

echo "Starting API tests based on OpenAPI specification..."

# Test health endpoint
echo "\n[Test] GET /health - Health check"
echo "Executing: http --print="b" GET ${BASE_URL}/health"
HEALTH_RESPONSE=$(http --print="b" GET "${BASE_URL}/health")
echo "Response body:"
echo "$HEALTH_RESPONSE" | jq '.'
if ! echo "$HEALTH_RESPONSE" | grep -q '"status": "ok"'; then
    echo "❌ Health check failed"
    exit 1
fi
echo "✅ Health check passed"

# Test user registration (POST /api/v1/register)
echo "\n[Test] POST /api/v1/register - User registration"
TEST_USERNAME=$(generate_username)
echo "Executing: http --print="b" POST ${BASE_URL}/api/v1/register username=${TEST_USERNAME} password=securepass123 email=${TEST_USERNAME}@example.com"
echo "Request body:"
jq -n \
  --arg username "$TEST_USERNAME" \
  --arg email "$TEST_USERNAME@example.com" \
  '{username: $username, password: "securepass123", email: $email}'

REGISTER_RESPONSE=$(http --print="b" POST "${BASE_URL}/api/v1/register" \
    username="${TEST_USERNAME}" \
    password="securepass123" \
    email="${TEST_USERNAME}@example.com")
echo "Response body:"
echo "$REGISTER_RESPONSE" | jq '.'

USER_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.id')
if [ -z "$USER_ID" ] || [ "$USER_ID" = "null" ]; then
    echo "❌ User registration failed"
    echo "Error response:"
    echo "$REGISTER_RESPONSE" | jq '.'
    exit 1
fi
echo "✅ User registration passed (ID: ${USER_ID})"

# Test user login (POST /api/v1/login)
echo "\n[Test] POST /api/v1/login - User authentication"
echo "Executing: http --print="b" POST ${BASE_URL}/api/v1/login username=${TEST_USERNAME} password=securepass123"
echo "Request body:"
jq -n \
  --arg username "$TEST_USERNAME" \
  '{username: $username, password: "securepass123"}'

LOGIN_RESPONSE=$(http --print="b" POST "${BASE_URL}/api/v1/login" \
    username="${TEST_USERNAME}" \
    password="securepass123")
echo "Response body:"
echo "$LOGIN_RESPONSE" | jq '.'

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')
if [ -z "$TOKEN" ]; then
    echo "❌ Login failed"
    echo "Error response:"
    echo "$LOGIN_RESPONSE" | jq '.'
    exit 1
fi
echo "✅ Login passed (Token received)"

# Set auth header for subsequent requests
AUTH_HEADER="Authorization: Bearer ${TOKEN}"

# Test get user by ID (GET /api/v1/users/{id})
echo "[Test] GET /api/v1/users/${USER_ID} - Get user by ID"
echo "Executing: http --print="b" GET ${BASE_URL}/api/v1/users/${USER_ID} ${AUTH_HEADER}"
USER_RESPONSE=$(http --print="b" GET "${BASE_URL}/api/v1/users/${USER_ID}" "$AUTH_HEADER")
if ! echo "$USER_RESPONSE" | jq -e ".id==${USER_ID}" > /dev/null; then
    echo "❌ Get user by ID failed"
    exit 1
fi
echo "✅ Get user by ID passed"

# Test update user (PUT /api/v1/users/{id})
echo "[Test] PUT /api/v1/users/${USER_ID} - Update user"
UPDATE_RESPONSE=$(http --print="b" PUT "${BASE_URL}/api/v1/users/${USER_ID}" \
    "$AUTH_HEADER" \
    username="${TEST_USERNAME}" \
    password="newpass123" \
    email="updated_${TEST_USERNAME}@example.com")
echo "Response body:"
echo "$UPDATE_RESPONSE" | jq '.'
if ! echo "$UPDATE_RESPONSE" | jq -e ".email | contains(\"updated_\")" > /dev/null; then
    echo "❌ Update user failed"
    exit 1
fi
echo "✅ Update user passed"

# Test delete user (DELETE /api/v1/users/{id})
echo "[Test] DELETE /api/v1/users/${USER_ID} - Delete user"
if ! http --print="b" DELETE "${BASE_URL}/api/v1/users/${USER_ID}" "$AUTH_HEADER"; then
    echo "❌ Delete user failed"
    exit 1
fi
echo "✅ Delete user passed"

# Verify deletion (GET /api/v1/users/{id} should return 404)
echo "[Test] GET /api/v1/users/${USER_ID} - Verify deletion"
VERIFY_RESPONSE=$(http GET "${BASE_URL}/api/v1/users/${USER_ID}" "$AUTH_HEADER" || true)
if echo "$VERIFY_RESPONSE" | grep -q '"error": "User not found"'; then
    echo "✅ User deletion verified (404 Not Found)"
else
    echo "❌ User still exists or unexpected response after deletion"
    echo "Response body:"
    echo "$VERIFY_RESPONSE" | jq '.'
    exit 1
fi
echo "✅ User deletion verified"

echo "All API tests completed successfully! ✨"

#!/bin/bash

# 性能测试脚本，自动登录获取JWT Token，并测试主要用户API
KONG_HOST="http://localhost:8000"
LOGIN_PATH="/api/v1/login"
REGISTER_PATH="/api/v1/register"
GETUSER_PATH="/api/v1/users/1"
GETUSERS_PATH="/api/v1/users"
UPDATEUSER_PATH="/api/v1/users/1"
CONCURRENCY=100
REQUESTS=10000

# 1. 生成随机用户名和邮箱
RAND_SUFFIX=$(cat /proc/sys/kernel/random/uuid 2>/dev/null || uuidgen)
RAND_SUFFIX=${RAND_SUFFIX:-$RANDOM}
USERNAME="benchuser_${RAND_SUFFIX}"
EMAIL="${USERNAME}@example.com"

# 2. 注册一个用户（只做一次，避免重复注册失败）
echo "Registering user..."
curl -s -X POST "$KONG_HOST$REGISTER_PATH" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"benchpass\",\"email\":\"$EMAIL\"}"

# 3. 登录获取token
echo "Logging in..."
LOGIN_RESP=$(curl -s -X POST "$KONG_HOST$LOGIN_PATH" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"benchpass\"}")
TOKEN=$(echo $LOGIN_RESP | jq -r '.token')
if [ -z "$TOKEN" ]; then
  echo "登录失败，无法获取token，退出测试。"
  echo "登录接口返回：$LOGIN_RESP"
  exit 1
fi

# 安装jq如果不存在
if ! command -v jq &> /dev/null; then
  echo "安装jq..."
  brew install jq
fi
AUTH_HEADER="Authorization: Bearer $TOKEN"
echo "Token: $TOKEN"

# 3. Benchmark Register（可选，通常注册只测一次）
# echo "Benchmarking Register ..."
# ab -c $CONCURRENCY -n $REQUESTS -p <(echo '{"username":"benchuser2","password":"benchpass","email":"benchuser2@example.com"}') -T 'application/json' "$KONG_HOST$REGISTER_PATH"
# echo "Benchmark for Register completed!"

# 4. Benchmark GetUser
echo "Benchmarking GetUser ..."
ab -c $CONCURRENCY -n $REQUESTS -H "$AUTH_HEADER" "$KONG_HOST$GETUSER_PATH"
echo "Benchmark for GetUser completed!"

# 5. Benchmark GetUsers
echo "Benchmarking GetUsers ..."
ab -c $CONCURRENCY -n $REQUESTS -H "$AUTH_HEADER" "$KONG_HOST$GETUSERS_PATH"
echo "Benchmark for GetUsers completed!"

# 6. Benchmark UpdateUser（用PATCH方式，ab不支持PATCH，使用curl循环模拟）
echo "Benchmarking UpdateUser ..."
for i in $(seq 1 $REQUESTS); do
  curl -s -X PATCH "$KONG_HOST$UPDATEUSER_PATH" \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d '{"first_name":"Bench","last_name":"User"}' > /dev/null &
  if (( $i % $CONCURRENCY == 0 )); then
    wait
  fi
done
wait
echo "Benchmark for UpdateUser completed!"

echo "All selected user API benchmarks completed via Kong gateway!"

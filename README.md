# Rate-Limiting API Gateway

A **lightweight API Gateway** built in **Golang** with **distributed token bucket rate limiting**.  
Supports **per-user and per-API rate limits**, **dynamic configuration updates**, and **atomic token consumption using Redis**.

---

## Features

- **Per-user and per-API rate limiting** using the **token bucket algorithm**
- **Dynamic configuration** stored in Redis (update limits without restarting the gateway)
- **Distributed state management** with Redis to handle multiple instances
- **Atomic token operations** using Lua scripts for consistency under high concurrency
- Lightweight **reverse proxy** forwarding requests to backend services
- Gin-based **middleware integration** for API request handling

---

## Tech Stack

- **Golang** – Backend language
- **Gin** – HTTP router and middleware
- **Redis** – In-memory database for token buckets and configuration
- **Lua** – Atomic token bucket operations
- **Docker** – Optional for running Redis

---

## Setup

### 1. Clone the repository

```bash
git clone https://github.com/ashifR17/Rate-Limiting-API-Gateway.git
cd Rate-Limiting-API-Gateway
```
### 2. Start Redis (Docker recommended)
```bash
docker run -d -p 6379:6379 redis
```

### 3. Set initial rate limit configuration
```bash
redis-cli SET rate_limit:config '{
  "global_capacity": 10,
  "global_rate": 1,
  "user_capacity": 5,
  "user_rate": 1,
  "user_api_capacity": 3,
  "user_api_rate": 1
}'
```
### 4. Start backend server
```bash
go run ./cmd/backend
```
### 5. Run the API Gateway
```bash
PORT=8081 go run main.go
The gateway will run on http://localhost:8081.
```

##Usage
### Test rate limiting
```bash
for i in {1..10}; do
  curl -s -o /dev/null -w "%{http_code}\n" -H "X-User-Id: user1" http://localhost:8081/api/ping
done
```
Returns 200 for allowed requests and 429 when limits are exceeded.
Tokens refill over time according to configuration.

### Inspect Redis Buckets
```bash
docker exec -it redis-gateway redis-cli keys "*"
docker exec -it redis-gateway redis-cli HGETALL rate:user:user1:api:/api/*proxyPath
```

## Future Enhancements
1. Add metrics dashboard for monitoring token usage

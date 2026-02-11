-- KEYS:
-- KEYS[1] = global bucket
-- KEYS[2] = user bucket
-- KEYS[3] = user_api bucket

-- ARGV:
-- 1 = global capacity
-- 2 = global rate (tokens/sec)
-- 3 = user capacity
-- 4 = user rate (tokens/sec)
-- 5 = user_api capacity
-- 6 = user_api rate (tokens/sec)
-- 7 = current timestamp (seconds)
-- 8 = TTL (seconds)

local now = tonumber(ARGV[7])
local ttl = tonumber(ARGV[8])

-- Helper function: refill and consume a bucket
local function refill_and_consume(key, capacity, rate, now)
    local bucket = redis.call("HMGET", key, "tokens", "last")
    local tokens = tonumber(bucket[1])
    local last = tonumber(bucket[2])

    if tokens == nil then
        tokens = capacity
        last = now
    else
        local delta = math.max(0, now - last)
        tokens = math.min(capacity, tokens + (rate * delta))
    end

    if tokens < 1 then
        -- Persist the updated bucket even if no tokens left
        redis.call("HSET", key, "tokens", tokens, "last", now)
        return false
    end

    -- Consume 1 token
    tokens = tokens - 1
    redis.call("HSET", key, "tokens", tokens, "last", now)
    return true
end

-- Refill & consume each bucket
local g_ok = refill_and_consume(KEYS[1], tonumber(ARGV[1]), tonumber(ARGV[2]), now)
local u_ok = refill_and_consume(KEYS[2], tonumber(ARGV[3]), tonumber(ARGV[4]), now)
local ua_ok = refill_and_consume(KEYS[3], tonumber(ARGV[5]), tonumber(ARGV[6]), now)

-- Set TTL for all buckets
redis.call("EXPIRE", KEYS[1], ttl)
redis.call("EXPIRE", KEYS[2], ttl)
redis.call("EXPIRE", KEYS[3], ttl)

-- Return 1 = allowed, 0 = rate-limited
if g_ok and u_ok and ua_ok then
    return 1
else
    return 0
end

package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiterConfig struct {
	MaxRequests int
	Window      time.Duration
}

type RedisRateLimiter struct {
	redis  *redis.Client
	config RateLimiterConfig
}

func NewRedisRateLimiter(redisClient *redis.Client, cfg RateLimiterConfig) *RedisRateLimiter {
	return &RedisRateLimiter{
		redis:  redisClient,
		config: cfg,
	}
}

// Token bucket implemented via a Lua script for atomicity.
// Lua scripts in Redis run atomically — no race conditions between
// read-modify-write steps across concurrent requests.
var tokenBucketScript = redis.NewScript(`
	local key        = KEYS[1]
	local capacity   = tonumber(ARGV[1])
	local refill_rate = tonumber(ARGV[2]) 
	local now        = tonumber(ARGV[3])   
	
	local data = redis.call("HMGET", key, "tokens", "last_refill")
	
	local tokens     = tonumber(data[1])
	local last_refill = tonumber(data[2])
	
	if tokens == nil then
		-- first request: initialise bucket with full capacity
		tokens      = capacity
		last_refill = now
	end
	
	-- calculate how many tokens to add based on elapsed time
	local elapsed = (now - last_refill) / 1000.0          -- convert ms to seconds
	local refill  = math.floor(elapsed * refill_rate)
	
	tokens = math.min(capacity, tokens + refill)
	
	if refill > 0 then
		last_refill = now
	end
	
	if tokens >= 1 then
		tokens = tokens - 1
		redis.call("HMSET", key, "tokens", tokens, "last_refill", last_refill)
		redis.call("PEXPIRE", key, math.ceil((capacity / refill_rate) * 1000))
		return 1   -- allowed
	else
		redis.call("HMSET", key, "tokens", tokens, "last_refill", last_refill)
		redis.call("PEXPIRE", key, math.ceil((capacity / refill_rate) * 1000))
		return 0   -- denied
	end`,
)

func (rl *RedisRateLimiter) allow(ctx context.Context, key string) (bool, error) {
	capacity := int64(rl.config.MaxRequests)
	refillRate := float64(rl.config.MaxRequests) / rl.config.Window.Seconds()
	nowMs := time.Now().UnixMilli()

	result, err := tokenBucketScript.Run(ctx, rl.redis, []string{key}, capacity, refillRate, nowMs).Int()
	if err != nil {
		return true, err // fail open — don't block traffic if Redis is down
	}

	return result == 1, nil
}

func (rl *RedisRateLimiter) LimitByIP() func(http.Handler) http.Handler {
	return rl.LimitByKey(func(r *http.Request) string {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		return fmt.Sprintf("rl:ip:%s", ip)
	})
}

func (rl *RedisRateLimiter) LimitByKey(keyFn func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFn(r)
			allowed, err := rl.allow(r.Context(), key)
			if err != nil {
				fmt.Printf("rate limiter error: %v\n", err)
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.config.MaxRequests))
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "rate limit exceeded",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

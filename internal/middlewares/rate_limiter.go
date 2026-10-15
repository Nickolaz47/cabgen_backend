package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type ipRateLimiter struct {
	mu      sync.Mutex
	entries map[string]*ipLimiterEntry
	rate    rate.Limit
	burst   int
	lastGC  time.Time
}

func (l *ipRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	if now.Sub(l.lastGC) > ipLimiterCleanupInterval {
		for k, e := range l.entries {
			if now.Sub(e.lastSeen) > ipLimiterIdleTimeout {
				delete(l.entries, k)
			}
		}
		l.lastGC = now
	}

	e, ok := l.entries[ip]
	if !ok {
		e = &ipLimiterEntry{limiter: rate.NewLimiter(l.rate, l.burst)}
		l.entries[ip] = e
	}
	e.lastSeen = now

	return e.limiter.Allow()
}

const (
	ipLimiterCleanupInterval = 5 * time.Minute
	ipLimiterIdleTimeout     = 10 * time.Minute
)

func RateLimitPerMinute(max int) gin.HandlerFunc {
	limiter := &ipRateLimiter{
		entries: make(map[string]*ipLimiterEntry),
		rate:    rate.Every(time.Minute / time.Duration(max)),
		burst:   max,
	}

	return func(c *gin.Context) {
		if !limiter.allow(c.ClientIP()) {
			localizer := translation.GetLocalizerFromContext(c)
			c.JSON(http.StatusTooManyRequests, responses.APIResponse{
				Error: responses.GetResponse(localizer,
					responses.TooManyRequestsError),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

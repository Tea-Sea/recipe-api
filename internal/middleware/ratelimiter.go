package middleware

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	Get    *rate.Limiter
	Post   *rate.Limiter
	Put    *rate.Limiter
	Delete *rate.Limiter
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		Get:    rate.NewLimiter(rate.Every(2*time.Second), 5),
		Post:   rate.NewLimiter(rate.Every(10*time.Second), 2),
		Put:    rate.NewLimiter(rate.Every(10*time.Second), 2),
		Delete: rate.NewLimiter(rate.Every(20*time.Second), 1),
	}
}

func (rates *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var limiter *rate.Limiter

		switch r.Method {
		case http.MethodGet:
			limiter = rates.Get
		case http.MethodPost:
			limiter = rates.Post
		case http.MethodPut:
			limiter = rates.Put
		case http.MethodDelete:
			limiter = rates.Delete
		default:
			next.ServeHTTP(w, r)
			return
		}

		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

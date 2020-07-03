package ratelimiter

import (
	"context"
	"net/http"

	"golang.org/x/time/rate"
)

var rateLimiterKey = struct{}{}

// Transport implements the http.RoundTripper interface.
//
// Once this Transport is registerred in your HTTP Client, The RoundTrip method
// gets executed before every HTTP Call.
// With this method we can automate the rate limiting process under the hood
// which reduces boilplate code in the HTTP Client Libraries.
//
// you can find more about the RoundTrippper interface here:
// https://golang.org/pkg/net/http/#RoundTripper
type Transport struct {
	limiters *Limiters
}

// NewTransport returns a new Transport instance with a limiters store.
func NewTransport(limiters *Limiters) *Transport {
	return &Transport{
		limiters: limiters,
	}
}

// RoundTrip wraps http.DefaultTransport.RoundTrip
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	limiter, ok := t.fromContext(req.Context())
	if ok {
		limiter.Wait(req.Context())
	}
	return http.DefaultTransport.RoundTrip(req)
}

// ToContext adds a rate limiter key to the context so that the Transport can find it.
func ToContext(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, rateLimiterKey, key)
}

func (t *Transport) fromContext(ctx context.Context) (*rate.Limiter, bool) {
	key, ok := ctx.Value(rateLimiterKey).(string)
	if !ok {
		return nil, false
	}
	limiter, err := t.limiters.Get(key)
	if err != nil {
		return nil, false
	}
	return limiter, true
}

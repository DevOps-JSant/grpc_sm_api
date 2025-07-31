package interceptors

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type rateLimiter struct {
	mu        sync.Mutex
	visitors  map[string]int
	limit     int
	resetTime time.Duration
}

func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors:  make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}
	go rl.resetVisitorCount()
	return rl
}

func (rl *rateLimiter) resetVisitorCount() {
	for {
		time.Sleep(rl.resetTime)
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) RateLimiterInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	log.Println("Received request in RateLimiter interceptor")

	rl.mu.Lock()
	defer rl.mu.Unlock()

	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to get client IP")
	}

	// TODO: Get ip from metadata incoming request instead of peer.Addr
	host, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse client IP")
	}

	visitorIP := host // simple visitor ip extraction
	rl.visitors[visitorIP]++
	log.Printf("Visitor count from %v is %v\n", visitorIP, rl.visitors[visitorIP])

	if rl.visitors[visitorIP] > rl.limit {
		// http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return nil, status.Error(codes.ResourceExhausted, "too many requests")
	}

	// Call the handler to proceed with the client request
	log.Println("Sending response from RateLimiter interceptor")
	return handler(ctx, req)

}

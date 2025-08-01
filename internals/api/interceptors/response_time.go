package interceptors

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ResponseTimeInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	log.Println("Received request in ResponseTime interceptor")
	// Record the start time
	start := time.Now()

	// Call the handler to proceed with the client request
	resp, err := handler(ctx, req)

	// Calculate time duration
	duration := time.Since(start)

	// Log the request details with duration
	st, _ := status.FromError(err)
	log.Printf("Method: %s, Status: %d Duration: %v\n", info.FullMethod, st.Code(), duration)

	md := metadata.Pairs("X-Response-Time", duration.String())
	grpc.SetHeader(ctx, md)

	log.Println("Sending response from ResponseTime interceptor")
	return resp, err
}

package grpcfuncs_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	pb "gophkeeper/proto"
)

func TestAuth(t *testing.T) {
	// Start the gRPC server in a separate goroutine
	go func() {
		// Start your gRPC server implementation
	}()

	// Connect the client to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	// Create a new gRPC client
	client := pb.NewGophkeeperClient(conn)

	// Perform the gRPC request
	request := &pb.AuthLoginRequest{
		// Set the request fields accordingly
	}
	response, err := client.Auth(context.Background(), request)

	// Assert the expected response
	assert.NoError(t, err)
	assert.NotNil(t, response)
	// Add more assertions as needed
}

// Add more test functions for other gRPC endpoints

package grpcfuncs_test

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"gophkeeper/internal/grpcfuncs"
	pb "gophkeeper/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestAuth(t *testing.T) {
	// Start the gRPC server in a separate goroutine
	g := grpcfuncs.NewGophKeeperServer()
	go func() {

		listen, err := net.Listen("tcp", ":3200")
		if err != nil {
			log.Fatal(err)
		}

		s := grpc.NewServer()
		pb.RegisterGophkeeperServer(s, &g)

		if err := s.Serve(listen); err != nil {
			log.Fatal(err)
		}
		// Start your gRPC server implementation
	}()

	// Connect the client to the server
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer conn.Close()

	// Create a new gRPC client
	client := pb.NewGophkeeperClient(conn)

	// Perform the gRPC request
	request := &pb.AuthLoginRequest{
		Login:    "test",
		Password: "password",
	}
	var header metadata.MD
	response, err := client.Auth(context.Background(), request, grpc.Header(&header))

	// Assert the expected response
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, header)
	// Add more assertions as needed
}

// Add more test functions for other gRPC endpoints

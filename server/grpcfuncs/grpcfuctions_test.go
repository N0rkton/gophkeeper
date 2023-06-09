package grpcfuncs_test

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"gophkeeper/server/grpcfuncs"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	pb "gophkeeper/proto"
)

func TestAuth(t *testing.T) {
	// Start the gRPC server in a separate goroutine

	go func() {
		grpcfuncs.Init()
		listen, err := net.Listen("tcp", ":3200")
		if err != nil {
			log.Fatal(err)
		}

		s := grpc.NewServer()
		pb.RegisterGophkeeperServer(s, &grpcfuncs.GophKeeperServer{})

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

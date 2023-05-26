package main

import (
	"google.golang.org/grpc"
	pb "gophkeeper/proto"
	"log"
	"net"
)

func main() {

	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))

	pb.RegisterGophkeeperServer(s, &GophKeeperServer{})

	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

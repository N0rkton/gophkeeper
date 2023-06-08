package main

import (
	"google.golang.org/grpc"
	pb "gophkeeper/proto"
	"gophkeeper/server/grpcfuncs"
	"log"
	"net"
)

func main() {
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
}

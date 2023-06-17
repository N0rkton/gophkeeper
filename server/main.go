package main

import (
	"gophkeeper/internal/grpcfuncs"
	pb "gophkeeper/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	gophKeeper := grpcfuncs.NewGophKeeperServer()
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterGophkeeperServer(s, &gophKeeper)

	if err = s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

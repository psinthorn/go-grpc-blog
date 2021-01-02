package main

import (
	"context"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/grpc"
)

type authServer struct {
}

func (authServer) Login(_ context.Context, in *proto.LoginRequest) (*proto.AuthResponse, error) {
	return &proto.AuthResponse{}, nil
}

func main() {
	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, authServer{})
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}

}

package main

import (
	"context"
	"log"
	"net"
	"time"

	authpb "github.com/psinthorn/go-grpc-blog/proto"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang/protobuf/protoc-gen-go/grpc"
	"github.com/psinthorn/go-grpc-blog/global"
	"go.mongodb.org/mongo-driver/bson"
)

type authServer struct{}

func (*authServer) Login(_ context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	login, password := req.GetLogin(), req.GetPassword()
	ctx, cancle := global.NewDBContext(5 * time.Second)
	defer cancle()

	var user global.User
	user := global.DB.Collection("user").FindOne(ctx, []bson.M{"$or": bson.M["username":login], bson.M["email":login]}).Decode(&user)
	if user == NilUser {
		return &authpb.AuthResponse{}, error.New("Wrong login credentials provided")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return &authpb.AuthResponse{}, error.New("Wrong login credentials provided")
	}

	return &authpb.AuthResponse{token: user.GetToken()}, nil

}

func main() {

	// สร้าง gRPC Server ขึ้นมาใหม่
	server := grpc.NewServer()

	// Register Auth Server
	authpb.RegisterAuthServiceServer(server, &authServer{})

	// สร้างค่า ตัวแปรสำหรับ Server
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("Error: ", err.Error())
	}

	// สั่งให้ server เริ่มต้นทำงาน
	server.Serve(listener)

}

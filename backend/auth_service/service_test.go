package main

import (
	"context"
	"testing"

	authpb "github.com/psinthorn/go-grpc-blog/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/psinthorn/go-grpc-blog/global"
	"golang.org/x/crypto/bcrypt"
)

func Test_authServer_Login(t *testing.T) {
	global.ConnectToTestDB()
	pw, _ := bcrypt.GenerateFromPassword([]byte("example"), bcrypt.DefaultCost)
	global.DB.Collection("user").InsertOne(context.Background(), global.User{ID: primitive.NewObjectID(), Email: "test@gmail.com", UserName: "test", Password: string(pw)})

	server := authServer{}
	_, err := server.Login(context.Background(), &authpb.LoginRequest{Login: "test@gmail.com", Passoword: "example"})
	if err != nil {
		t.Error("1: Error was return: ", err.Error())
	}
	// t.Error(res)
	_, err = server.Login(context.Background(), &authpb.LoginRequest{Login: "somethings", Passoword: "somethings"})
	if err == nil {
		t.Error("2: Error was return: ", err.Error())
	}
	// t.Error(res)

	_, err = server.Login(context.Background(), &authpb.LoginRequest{Login: "test", Passoword: "example"})
	if err != nil {
		t.Error("3: Error was return: ", err.Error())
	}
}

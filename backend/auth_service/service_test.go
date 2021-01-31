package main

import (
	"context"
	"testing"

	"github.com/psinthorn/go-grpc-blog/global"
	authpb "github.com/psinthorn/go-grpc-blog/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Test_authServer_Login(t *testing.T) {
	global.ConnectToTestDB()
	pw, _ := bcrypt.GenerateFromPassword([]byte("example"), bcrypt.DefaultCost)
	global.DB.Collection("user").InsertOne(context.Background(), global.User{ID: primitive.NewObjectID(), Email: "test@gmail.com", UserName: "test", Password: string(pw)})

	server := authServer{}

	// Test login with correct user and passwor should not retuen any error
	_, err := server.Login(context.Background(), &authpb.LoginRequest{Login: "test@gmail.com", Password: "example"})
	if err != nil {
		t.Error("1: Error was return: ", err.Error())
	}

	// Test login with wrong username and password
	_, err = server.Login(context.Background(), &authpb.LoginRequest{Login: "somethings", Password: "somethings"})
	if err == nil {
		t.Error("2: Error was return: ", err.Error())
	}

	// Test login with wrong user name only
	_, err = server.Login(context.Background(), &authpb.LoginRequest{Login: "test", Password: "example"})
	if err != nil {
		t.Error("3: Error was return: ", err.Error())
	}
}

func Test_authServer_Signup(t *testing.T) {
	server := authServer{}
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{UserName: "UniqueUser", Email: "uniqueemail@test.com"})

	// Test both username and email is available (not unique)
	_, err := server.Signup(context.Background(), &authpb.SignupRequest{UserName: "notUniqueUser", Email: "notuniqueemail@test.com", Password: "thisIsPassword"})
	if err != nil {
		t.Error("1: Shold not get any errors")
	}

	// Test UserName is Unique
	_, err = server.Signup(context.Background(), &authpb.SignupRequest{UserName: "UniqueUser", Email: "Notuniqueemail@test.com", Password: "thisIsPassword"})
	if err.Error() != "Username is exist" {
		t.Error("3: Username should be exist")
	}

	// Test Email is Unique
	_, err = server.Signup(context.Background(), &authpb.SignupRequest{UserName: "NotUniqueUser", Email: "uniqueemail@test.com", Password: "thisIsPassword"})
	if err.Error() != "Email is exist" {
		t.Error("4: Email addressShold should be exist")
	}

	// Test password is less than 8 charecters
	_, err = server.Signup(context.Background(), &authpb.SignupRequest{UserName: "notUniqueUser", Email: "notuniqueemail@test.com", Password: "Passwor"})
	if err.Error() != "Password must be at lease 8 charecters" {
		t.Error("5: Should get error Password must be at lease 8 charecters")
	}

}

func Test_authServer_UniqueUserNameValidate(t *testing.T) {
	server := authServer{}
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{UserName: "test"})
	res, err := server.UniqueUserNameValidate(context.Background(), &authpb.UniqueUserNameValidateRequest{UserName: "tester"})
	if err != nil {
		t.Error("1: Error was return: ", err.Error())
	}
	if res.GetIsUnique() {
		t.Error("1: Testing with Wrong result")
	}

	res, err = server.UniqueUserNameValidate(context.Background(), &authpb.UniqueUserNameValidateRequest{UserName: "test"})
	if err != nil {
		t.Error("2: Error was return: ", err.Error())
	}
	if !res.GetIsUnique() {
		t.Error("2: Testing with Wrong result")
	}

}

func Test_authServer_UniqueEmailValidate(t *testing.T) {
	server := authServer{}
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{Email: "test@test.com"})
	res, err := server.UniqueEmailValidate(context.Background(), &authpb.UniqueEmailValidateRequest{Email: "tester@test.com"})
	if err != nil {
		t.Error("1: Error was return: ", err.Error())
	}
	if res.GetIsUnique() {
		t.Error("1: Testing with Wrong result")
	}

	res, err = server.UniqueEmailValidate(context.Background(), &authpb.UniqueEmailValidateRequest{Email: "test@test.com"})
	if err != nil {
		t.Error("2: Error was return: ", err.Error())
	}
	if !res.GetIsUnique() {
		t.Error("2: Testing with Wrong result")
	}

}

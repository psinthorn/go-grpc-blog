package main

import (
	"context"
	"errors"
	"log"
	"net"
	"regexp"
	"time"

	authpb "github.com/psinthorn/go-grpc-blog/proto"
	"google.golang.org/grpc"

	"github.com/psinthorn/go-grpc-blog/global"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type authServer struct{}

func (s *authServer) Signup(_ context.Context, req *authpb.SignupRequest) (*authpb.AuthResponse, error) {
	// สรุปขั้นตอนการ Signup
	// 1. รับค่าจาก input
	// 2. ตรวจสอบการค่าที่ได้รับว่าตรงกับข้อกำหนดหรือไม่ (validation) หากไม่ตรงกับข้อกำหนดให้แจ้ง (returns error) user ให้ทราบ
	// 2.1 ตรวจสอบว่ามี user หรือ email อยู่ในระบบปัจุบันแล้วหรือไม่หากมีแล้วให้ส่งข้อมูลหรือแจ้งให้ (returns &authpb.AuthResponse or  error) user ได้ทราบ (*สร้าง Function ในการตรวจ  current user และ current email )
	// 2.2 หากไม่พบข้อผิดพลาดให้ดำเนินการขั้นต่อไป
	// 3. การเชื่อมต่ฐานข้อมูลและนำข้อมูลเก็บลงฐานข้อมูล
	// 4. ทำการ Login อัตโนมัติให้

	// เริ่มปฎิบัติตามลำดับขั้นตอนด้านบน
	// 1. รับค่าจาก input
	userName, email, password := req.GetUserName(), req.GetEmail(), req.GetPassword()

	// 2. ตรวจสอบการค่าที่ได้รับว่าตรงกับข้อกำหนดหรือไม่ (validation) หากไม่ตรงกับข้อกำหนดให้แจ้ง (returns error) user ให้ทราบ

	// 2.1 ตรวจสอบความถูกต้องของ email
	match, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", email)
	if !match {
		return &authpb.AuthResponse{}, errors.New("Email address is not valid")
	}
	if len(email) < 8 {
		return &authpb.AuthResponse{}, errors.New("Email address must be at lease 8 charecters")
	}

	// 2.1 ตรวจสอบความถูกต้องของ userName
	if len(userName) < 8 {
		return &authpb.AuthResponse{}, errors.New("User name must be at lease 8 charecters")
	}
	if len(password) < 8 {
		return &authpb.AuthResponse{}, errors.New("Password must be at lease 8 charecters")
	}

	// // Unique Username and Email Validation
	// _, err := s.UniqueUserNameValidate(context.Background(), &authpb.UniqueUserNameValidateRequest{UserName: userName})
	// if err != nil {
	// 	return &authpb.AuthResponse{}, err
	// }
	// _, err = s.UniqueEmailValidate(context.Background(), &authpb.UniqueEmailValidateRequest{Email: userName})
	// if err != nil {
	// 	if err != nil {
	// 		return &authpb.AuthResponse{}, err
	// 	}
	// }

	// 3. การเชื่อมต่ฐานข้อมูลและนำข้อมูลเก็บลงฐานข้อมูล
	var newUser global.User
	ctx, cancle := global.NewDBContext(10 * time.Second)
	defer cancle()

	res, err := s.UniqueUserNameValidate(context.Background(), &authpb.UniqueUserNameValidateRequest{UserName: userName})
	if err != nil {
		log.Println("Somthings went wrong")
		return &authpb.AuthResponse{}, err
	}
	if res.GetIsUnique() {
		return &authpb.AuthResponse{}, errors.New("Username is exist")
	}

	res, err = s.UniqueEmailValidate(context.Background(), &authpb.UniqueEmailValidateRequest{Email: email})
	if err != nil {
		log.Println("Somthings went wrong")
		return &authpb.AuthResponse{}, err
	}
	if res.GetIsUnique() {
		return &authpb.AuthResponse{}, errors.New("Email is exist")
	}

	pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	newUser = global.User{ID: primitive.NewObjectID(), UserName: userName, Email: email, Password: string(pw)}

	ctx, cancle = global.NewDBContext(5 * time.Second)
	defer cancle()
	_, err = global.DB.Collection("user").InsertOne(ctx, newUser)
	if err != nil {
		log.Println("Error on insert new user:", err.Error())
		return &authpb.AuthResponse{}, err
	}

	return &authpb.AuthResponse{Token: newUser.GetToken()}, nil
}

func (s *authServer) Login(_ context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {

	// สรุปขั้นตอนในการ login
	// 1. รับค่า username และ password มาจาก user input form
	// 2. ทำการเชื่อมต่อฐานข้อมูล ขึ้นอยู่กับว่าเราเลือกใช้ฐานข้ออะไร เช่น mongodb, mySql, Posgres เป็นต้น
	// 3. สร้างตัวแปร โดยมี type เป็น User Struct ที่เราได้ทำการสร้างไว้ที่ global folder/package
	// 4. ค้นหา user ที่ฐานข้อมูลและ decode ให้อยู่ในรูปแบบ user struct ว่ามีตรงตามที่ input มาหรือไม่
	// 4.1 หากไม่มีให้คืนค่าเป็น user struct เปล่า (nil) และในส่วนของ error ให้ส่งคำอธิบาย error กลับให้รับทราบ
	// 5. หากพบ user ให้ทำการตรวจสอบเปรียบเทียบ password ว่าถูกต้องหรือไม่
	// 5.1 หากทำการตรวจสอบเปรียบเทียบ password แล้วไม่ตรงให้ทำการส่ง error กลับ
	// 6. หาก password ถูกต้องให้ทำหาร generate token และส่งกลับให้ user ต่อไป

	// 1. รับค่า username และ password มาจาก user input form
	login, password := req.GetLogin(), req.GetPassword()

	// 2. ทำการเชื่อมต่อฐานข้อมูล ขึ้นอยู่กับว่าเราเลือกใช้ฐานข้ออะไร เช่น mongodb, mySql, Posgres เป็นต้น
	ctx, cancle := global.NewDBContext(5 * time.Second)
	defer cancle()

	// 3. สร้างตัวแปร โดยมี type เป็น User Struct ที่เราได้ทำการสร้างไว้ที่ global folder/package
	var user global.User

	// 4. ค้นหา user ที่ฐานข้อมูลและ decode ให้อยู่ในรูปแบบ user struct ว่ามีตรงตามที่ input มาหรือไม่
	// 4.1 หากไม่มีให้คืนค่าเป็น user struct เปล่า (nil) และในส่วนของ error ให้ส่งคำอธิบาย error กลับให้รับทราบ
	global.DB.Collection("user").FindOne(ctx, bson.M{"$or": []bson.M{bson.M{"username": login}, bson.M{"email": login}}}).Decode(&user)
	if user == global.NilUser {
		return &authpb.AuthResponse{}, errors.New("Wrong login credentials provided")
	}

	// 5. หากพบ user ให้ทำการตรวจสอบเปรียบเทียบ password ว่าถูกต้องหรือไม่
	// 5.1 หากทำการตรวจสอบเปรียบเทียบ password แล้วไม่ตรงให้ทำการส่ง error กลับ
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return &authpb.AuthResponse{}, errors.New("Wrong login credentials provided")
	}

	// 6. หาก password ถูกต้องให้ทำหาร generate token และส่งกลับให้ user ต่อไป
	return &authpb.AuthResponse{Token: user.GetToken()}, nil

}

// Validate Unique UserName
func (s *authServer) UniqueUserNameValidate(_ context.Context, req *authpb.UniqueUserNameValidateRequest) (*authpb.UniqueValidateResponse, error) {
	var result global.User
	userName := req.GetUserName()
	ctx, cancle := global.NewDBContext(5 * time.Second)
	defer cancle()
	global.DB.Collection("user").FindOne(ctx, bson.M{"username": userName}).Decode(&result)
	return &authpb.UniqueValidateResponse{IsUnique: result != global.NilUser}, nil
}

// Validate Unique Email addtess
func (s *authServer) UniqueEmailValidate(_ context.Context, req *authpb.UniqueEmailValidateRequest) (*authpb.UniqueValidateResponse, error) {
	var result global.User
	email := req.GetEmail()
	ctx, cancle := global.NewDBContext(5 * time.Second)
	defer cancle()
	global.DB.Collection("user").FindOne(ctx, bson.M{"email": email}).Decode(&result)

	return &authpb.UniqueValidateResponse{IsUnique: result != global.NilUser}, nil
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

package main

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	authpb "github.com/psinthorn/go-grpc-blog/proto"
	"google.golang.org/grpc"

	"github.com/psinthorn/go-grpc-blog/global"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type authServer struct{}

func (*authServer) Login(_ context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {

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
	login, password := req.GetLogin(), req.GetPassoword()

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

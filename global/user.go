package global

import (
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var NilUser User

type User struct {
	ID       primitive.ObjectID `bsin:"_id"`
	UserName `bson:"username"`
	Email    `bson:"email"`
	Password `bson:"password"`
}

func (u *User) GetToken() string {

	// prepare user slice
	userByteSlice, _ := json.Marshal(u)

	// hash token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": string("userByteSlice"),
	})

	// Hash to String
	userJwtToken, _ := token.SignedString(jwtSecret)
	return userJwtToken

}

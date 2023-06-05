package entity

import (
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Firstname   string             `json:"firstname" bson:"firstname"`
	Lastname    string             `json:"lastname" bson:"lastname"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	Password    string             `json:"password" bson:"password"`
	Role        string             `json:"role" bson:"role"`
	CreatedAt   int64              `json:"created_at" bson:"created_at"`
}

func (u *User) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      u.ID.Hex(),
		"firstname":    u.Firstname,
		"lastname":     u.Lastname,
		"role":         u.Role,
		"phone_number": u.PhoneNumber,
		"created_at":   u.CreatedAt,
		"exp":          time.Now().Add(time.Hour * 24 * 31).Unix(),
	})
	return token.SignedString([]byte("EiXooniesae4aegh0av1aith2oaheesh"))
}

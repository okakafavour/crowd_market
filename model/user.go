package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name             string             `bson:"name" json:"name"`
	Email            string             `bson:"email" json:"email"`
	Password         string             `bson:"password,omitempty" json:"password,omitempty"`
	Provider         string             `bson:"provider" json:"provider"`
	IsVerified       bool               `bson:"isVerified" json:"isVerified"`
	VerificationCode string             `bson:"verificationCode,omitempty" json:"verificationCode,omitempty"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}

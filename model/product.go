package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
	Name        string             `bson:"name" json:"name" binding:"required"`
	Price       float64            `bson:"price" json:"price" binding:"required"`
	Area        string             `bson:"area" json:"area" binding:"required"`
	Description string             `bson:"description" json:"description" binding:"required"`
	ImageURL    string             `bson:"image_url" json:"image_url"`
	Category    string             `bson:"category" json:"category" binding:"category"`
	CreatedAt   *time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt   *time.Time         `bson:"updated_at" json:"updated_at"`
}

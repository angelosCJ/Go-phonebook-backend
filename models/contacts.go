package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Contact struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Phone string             `json:"phone" bson:"phone"`
	Email string             `json:"email" bson:"email"`
}

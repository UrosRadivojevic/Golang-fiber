package user_request

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserRequest struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" validate:"required"`
	Password  string             `json:"password" validate:"required,gte=8"`
	Firstname string             `json:"firstname"`
}

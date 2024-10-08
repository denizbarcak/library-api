package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Name     string             `json:"name"`
    Email    string             `json:"email"`
    Password string             `json:"password"`
    Books    []primitive.ObjectID `json:"books"`
}

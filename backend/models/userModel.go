package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName          string             `bson:"fullName" json:"fullName"`
	Email             string             `bson:"email" json:"email"`
	PasswordHash      string             `bson:"passwordHash" json:"-"`
	Country           string             `bson:"country" json:"country"`
	PreferredIndustry string             `bson:"preferredIndustry" json:"preferredIndustry"`
	UserType          string             `bson:"userType" json:"userType"`
	CreatedAt         time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type UserPublic struct {
	ID                string    `json:"id"`
	FullName          string    `json:"fullName"`
	Email             string    `json:"email"`
	Country           string    `json:"country"`
	PreferredIndustry string    `json:"preferredIndustry"`
	UserType          string    `json:"userType"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func ToUserPublic(u User) UserPublic {
	return UserPublic{
		ID:                u.ID.Hex(),
		FullName:          u.FullName,
		Email:             u.Email,
		Country:           u.Country,
		PreferredIndustry: u.PreferredIndustry,
		UserType:          u.UserType,
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
	}
}

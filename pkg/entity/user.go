package entity

import (
	"gorm.io/gorm"
)

type EditUser struct {
	FirstName string `json:"firstname" bson:"firstname" binding:"required"`
	LastName  string `json:"lastname" bson:"lastname" binding:"required"`
	Email     string `json:"email" bson:"email" binding:"required"`
}

type Signup struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type User struct {
	gorm.Model `json:"-"`
	ID         int    `gorm:"primarykey" bson:"_id,omitempty" json:"-"`
	FirstName  string `json:"firstname" bson:"firstname" binding:"required"`
	LastName   string `json:"lastname" bson:"lastname" binding:"required"`
	Email      string `json:"email" bson:"email" binding:"required"`
	Phone      string `json:"phone" bson:"phone" binding:"required"`
	Password   string `json:"-" bson:"password" binding:"required"`
	Wallet     int    `json:"wallet"`
	Permission bool   `gorm:"not null;default:true" json:"-"`
}

type Address struct {
	gorm.Model `json:"-"`
	ID         int    `gorm:"primarykey" json:"id"`
	UserId     int    `json:"-"`
	House      string `json:"house"`
	Street     string `json:"street"`
	City       string `json:"city"`
	Pincode    int    `json:"pincode"`
	Type       string `json:"type"`
}

type Login struct {
	Phone    string `json:"phone" bson:"Phone" binding:"required"`
	Password string `json:"password" bson:"Password" binding:"required"`
}

type OtpKey struct {
	gorm.Model
	Key   string `json:"key"`
	Phone string `json:"phone"`
}

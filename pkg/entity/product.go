package entity

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model  `json:"-"`
	ID          int       `gorm:"primarykey"`
	Name        string    `json:"name"`
	Price       int       `json:"price"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
	ImageURL    string    `json:"imageurl"`
	Removed     bool      `json:"-" gorm:"default false"`
	Category    string    `json:"category" gorm:"default ticket"`
	SubCategory string    `json:"subcategory"`
	AdminId     int       `json:"-"`
}

type TicketDetails struct {
	gorm.Model  `json:"-"`
	TicketId    int    `json:"ticketid"`
	Description string `json:"description"`
	Venue       string `json:"venue"`
}
type TicketInput struct {
	Ticket
	TicketDetails
	Inventory
}
type Apparel struct {
	gorm.Model  `json:"-"`
	ID          int    `gorm:"primarykey" json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	ImageURL    string `json:"imageurl"`
	Removed     bool   `json:"-" gorm:"default false"`
	Category    string `json:"category" gorm:"default apparel"`
	SubCategory string `json:"subcategory"`
	AdminId     int    `json:"-"`
}

type ApparelDetails struct {
	gorm.Model  `json:"-"`
	ApparelID   uint   `json:"apparelid"`
	Description string `json:"description"`
	Size        string `json:"size"`
	Color       string `json:"color"`
}
type ApparelInput struct {
	Apparel
	ApparelDetails
	Inventory
}

type Inventory struct {
	gorm.Model      `json:"-"`
	ProductId       int
	ProductCategory string
	Quantity        int
}

type Coupon struct {
	Id         int       `json:"-" gorm:"primarykey"`
	Code       string    `json:"code"`
	Type       string    `json:"type"`
	Amount     int       `json:"amount"`
	ValidFrom  time.Time `json:"-"`
	ValidUntil time.Time `json:"valid_until"`
	UsageLimit int       `json:"usage_limit"`
	UsedCount  int       `json:"-"`
	Category   string    `json:"category"`
	AdminId    int       `json:"-"`
}

type Offer struct {
	Id         int       `json:"-" gorm:"primarykey"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Amount     int       `json:"amount"`
	MinPrice   int       `json:"minprice"`
	ValidFrom  time.Time `json:"-"`
	ValidUntil time.Time `json:"valid_until"`
	UsageLimit int       `json:"usage_limit"`
	UsedCount  int       `json:"-"`
	Category   string    `json:"category"`
	AdminId    int       `json:"-"`
}

type UsedCoupon struct {
	UserId     int    `json:"userid"`
	CouponCode string `json:"couponcode"`
}

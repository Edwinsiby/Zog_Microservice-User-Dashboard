package entity

import "gorm.io/gorm"

type Cart struct {
	gorm.Model      `json:"-"`
	UserId          int     `json:"-"`
	ApparelQuantity int     `json:"apparelquantity"`
	TicketQuantity  int     `json:"ticketquantity"`
	TotalPrice      float64 `json:"totalprice"`
	OfferPrice      int     `json:"offerprice"`
}

type CartItem struct {
	gorm.Model  `json:"-"`
	CartId      int     `json:"-"`
	Category    string  `json:"category"`
	ProductId   int     `json:"productid"`
	Quantity    int     `json:"quantity"`
	ProductName string  `json:"productname"`
	Price       float64 `json:"price"`
}

type Wishlist struct {
	gorm.Model  `json:"-"`
	UserId      int     `json:"-"`
	Category    string  `json:"category"`
	ProductId   int     `json:"productid"`
	ProductName string  `json:"productname"`
	Price       float64 `json:"price"`
}

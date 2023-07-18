package entity

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model    `json:"-"`
	ID            int     `gorm:"primarykey" json:"id"`
	UserID        int     `json:"userid"`
	AddressId     int     `json:"adressid"`
	Total         float64 `json:"total"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"paymentmethod"`
	PaymentStatus string  `json:"paymentstatus"`
	PaymentId     string  `json:"paymentid"`
}

type OrderItem struct {
	gorm.Model `json:"-"`
	OrderID    int     `json:"order_id"`
	ProductID  int     `json:"product_id"`
	Category   string  `json:"category"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

type Return struct {
	gorm.Model `json:"-"`
	OrderId    int    `json:"orderid"`
	UserId     int    `json:"userid"`
	Reason     string `json:"reason"`
	Status     string `json:"status"`
	Refund     string `json:"refund"`
	TotalPrice int    `json:"totalprice"`
}

type Invoice struct {
	gorm.Model  `json:"-"`
	OrderId     int     `json:"orderid"`
	UserId      int     `json:"userid"`
	AddressType string  `json:"addresstype"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Payment     string  `json:"payment"`
	Status      string  `json:"status"`
	PaymentId   string  `json:"paymentid"`
	Remark      string  `json:"remark" gorm:"default zog_festiv"`
}

type SalesReport struct {
	TotalSales       float64
	TotalOrders      int64
	AverageOrder     float64
	PaymentMethod    string
	PaymentMethodCnt int
}

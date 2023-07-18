package entity

import "gorm.io/gorm"

type Admin struct {
	gorm.Model `json:"-"`
	AdminName  string `json:"adminname"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	Active     bool   `gorm:"not null;default true"`
}

type AdminDashboard struct {
	TotalUsers        int    `json:"totalusers"`
	NewUsers          int    `json:"newusers"`
	TotalProducts     int    `json:"totalproducts"`
	StocklessCategory string `json:"stocklesscategory"`
	TotalOrders       int    `json:"totalorders"`
	AverageOrderValue int    `json:"averageordervalue"`
	PendingOrders     int    `json:"pendingorders"`
	ReturnOrders      int    `json:"returnorders"`
	TotalRevenue      int    `json:"totalrevenue"`
	TotalQuery        int    `json:"totalquery"`
}

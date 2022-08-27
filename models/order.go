package models

import (
	"time"
)

type Order struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	ProductId int       `json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductId"`
	UserId    int       `json:"user_id"`
	User      User      `gorm:"foreignKey:UserId"`
}

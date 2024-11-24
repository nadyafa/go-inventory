package models

import "time"

type Inventory struct {
	ProductID uint      `json:"product_id" gorm:"index;not null;foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Stock     uint      `json:"stock"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

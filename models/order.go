package models

import (
	"time"
)

type Order struct {
	OrderID   uint      `json:"order_id" gorm:"primaryKey;autoIncrement"`
	ProductID uint      `json:"product_id" gorm:"index;not null;foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Quantity  uint      `json:"quantity"`
	OrderDate time.Time `json:"order_date" gorm:"autoCreateTime"`
}

type OrderDto struct {
	OrderID   uint      `json:"order_id"`
	ProductID uint      `json:"product_id"`
	Quantity  uint      `json:"quantity" binding:"required"`
	OrderDate time.Time `json:"order_date"`
}

func (o *OrderDto) FillFromModel(model Order) {
	o.OrderID = model.OrderID
	o.ProductID = model.ProductID
	o.Quantity = model.Quantity
	o.OrderDate = model.OrderDate
}

func (o *OrderDto) ToModel() Order {
	model := Order{
		OrderID:   o.OrderID,
		ProductID: o.ProductID,
		Quantity:  o.Quantity,
	}

	return model
}

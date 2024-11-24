package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ProductID   uint            `json:"product_id" gorm:"primaryKey;autoIncrement"`
	Name        string          `json:"name"`
	Description sql.NullString  `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Category    string          `json:"category"`
	ImagePath   sql.NullString  `json:"image_path"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Inventory   Inventory       `json:"inventory" gorm:"foreignKey:ProductID"`
}

type ProductDto struct {
	ProductID   uint      `json:"product_id"`
	Name        string    `json:"name" binding:"required"`
	Description *string   `json:"description"`
	Price       int       `json:"price"`
	Category    string    `json:"category"`
	ImagePath   *string   `json:"image_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// from DB to client
func (p *ProductDto) FillFromModel(model Product) {
	p.ProductID = model.ProductID
	p.Name = model.Name
	p.Price = int(model.Price.IntPart())
	p.Category = model.Category
	if model.Description.Valid {
		p.Description = &model.Description.String
	}
	if model.ImagePath.Valid {
		p.ImagePath = &model.ImagePath.String
	}
	p.CreatedAt = model.CreatedAt
	p.UpdatedAt = model.UpdatedAt
}

// from client to DB
func (p *ProductDto) ToModel() Product {
	model := Product{
		ProductID: p.ProductID,
		Name:      p.Name,
		Price:     decimal.NewFromInt(int64(p.Price)),
		Category:  p.Category,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
	if p.Description != nil {
		model.Description.String = *p.Description
		model.Description.Valid = true
	}
	if p.ImagePath != nil {
		model.ImagePath.String = *p.ImagePath
		model.ImagePath.Valid = true
	}

	return model
}

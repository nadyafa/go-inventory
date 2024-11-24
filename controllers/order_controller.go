package controllers

import (
	"errors"
	"fmt"
	"inventory-management/models"
	"inventory-management/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateOrder(c *gin.Context) {
	var orderDto models.OrderDto

	err := c.ShouldBindJSON(&orderDto)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse(fmt.Sprintf("Failed to bind request: %s", err.Error())),
		)
		return
	}

	// Checking if the productID exists
	var product models.Product
	if err := repositories.DB.First(&product, orderDto.ProductID).Error; err != nil {
		c.JSON(
			http.StatusNotFound,
			models.NewFailedResponse("Product not found"),
		)
		return
	}

	// Checking productID existence in inventory
	var inventory models.Inventory
	if err := repositories.DB.First(&inventory, orderDto.ProductID).Error; err != nil {
		c.JSON(
			http.StatusNotFound,
			models.NewFailedResponse("Product stock unavailable"),
		)
		return
	}

	// Ensure the order quantity is less than or equal to the available stock
	if inventory.Stock < orderDto.Quantity {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse(fmt.Sprintf("Sorry only %d unit(s) available. Please change the order details", inventory.Stock)),
		)
		return
	}

	// Begin the transaction
	tx := repositories.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update inventory stock
	inventory.Stock -= orderDto.Quantity
	if err := tx.Model(&inventory).Where("product_id = ?", orderDto.ProductID).Updates(models.Inventory{
		Stock:     inventory.Stock,
		UpdatedAt: time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse(fmt.Sprintf("Failed to update inventory: %s", err.Error())),
		)
		return
	}

	// Create the order
	order := orderDto.ToModel()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse(fmt.Sprintf("Failed to create order: %s", err.Error())),
		)
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse(fmt.Sprintf("Failed to commit transaction: %s", err.Error())),
		)
		return
	}

	// Send back the successfully saved order for the response
	orderDto.FillFromModel(order)
	c.JSON(http.StatusCreated, orderDto)
}

func ReadOrderById(c *gin.Context) {
	// parse and validate id req
	id, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse(fmt.Sprintf("Invalid id: %s", err.Error())),
		)
		return
	}

	var order models.Order
	err = repositories.DB.First(&order, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(
			http.StatusNotFound,
			models.NewFailedResponse("Order history not found"),
		)
		return
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse(fmt.Sprintf("Failed to get order history: %s", err.Error())),
		)
		return
	}

	var orderDto models.OrderDto
	orderDto.FillFromModel(order)

	c.JSON(
		http.StatusOK,
		models.NewSuccessResponse("Success", orderDto),
	)
}

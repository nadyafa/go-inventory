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

func ReadInventories(c *gin.Context) {
	var inventory []models.Inventory

	// fetching data
	query := repositories.DB.Model(&models.Inventory{})
	if err := query.Find(&inventory).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse(fmt.Sprintf("Failed to fetch inventories: %s", err.Error())),
		)
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func ReadInventoryByProductID(c *gin.Context) {
	// parse and validate id req
	id, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse(fmt.Sprintf("Invalid id: %s", err.Error())),
		)
		return
	}

	var inventory models.Inventory
	err = repositories.DB.First(&inventory, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(
			http.StatusNotFound,
			models.NewFailedResponse("Product inventory not found"),
		)
		return
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse(fmt.Sprintf("Failed to get product inventory: %s", err.Error())),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		models.NewSuccessResponse("Success", inventory),
	)
}

func UpdateInventoryByProductID(c *gin.Context) {
	// parse and validate id req
	id, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse(fmt.Sprintf("Invalid id: %s", err.Error())),
		)
		return
	}

	// bind req body input to DTO
	var newInventory models.Inventory
	if err := c.ShouldBind(&newInventory); err != nil {
		// response to bad input
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse(fmt.Sprintf("Invalid request payload: %s", err.Error())),
		)
		return
	}

	// fetch if product exist
	var existingInventory models.Inventory
	err = repositories.DB.First(&existingInventory, "product_id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				models.NewFailedResponse(fmt.Sprintf("Inventory not found: %s", err.Error())),
			)
		} else {
			c.JSON(
				http.StatusInternalServerError,
				models.NewFailedResponse(fmt.Sprintf("Failed to fetch inventory: %s", err.Error())),
			)
		}
		return
	}

	err = repositories.DB.Model(&existingInventory).Where("product_id = ?", id).Updates(models.Inventory{
		Stock:     newInventory.Stock,
		Location:  newInventory.Location,
		CreatedAt: existingInventory.CreatedAt,
		UpdatedAt: time.Now(),
	}).Error
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse(fmt.Sprintf("Failed to update product: %s", err.Error())),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		models.NewSuccessResponse("Inventory successfully updated", existingInventory),
	)
}

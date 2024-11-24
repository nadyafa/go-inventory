package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"inventory-management/models"
	"inventory-management/repositories"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddProduct(c *gin.Context) {
	var productDto models.ProductDto

	err := c.ShouldBindJSON(&productDto)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse(fmt.Sprintf("Failed to bind request: %s", err.Error())),
		)
		return
	}

	// convert client input ToModel to insert into DB
	product := productDto.ToModel()

	// save product to DB
	err = repositories.DB.Create(&product).Error
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse(fmt.Sprintf("Failed to save product: %s", err.Error())),
		)
		return
	}

	// create an inventory based on productID
	inventory := models.Inventory{
		ProductID: product.ProductID,
		Stock:     0,
		Location:  "",
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}

	err = repositories.DB.Create(&inventory).Error
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse(fmt.Sprintf("Failed to create inventory: %s", err.Error())),
		)
		return
	}

	// send back successfully saved product for respond
	productDto.FillFromModel(product)
	c.JSON(http.StatusCreated, productDto)
}

func ReadProducts(c *gin.Context) {
	var products []models.ProductDto

	// filter date based on name or category
	nameFilter := c.DefaultQuery("name", "")
	categoryFilter := c.DefaultQuery("category", "")

	query := repositories.DB.Model(&models.Product{})
	if nameFilter != "" {
		query = query.Where("name LIKE ?", "%"+nameFilter+"%")
	}

	if categoryFilter != "" {
		query = query.Where("category = ?", categoryFilter)
	}

	// fetching data
	if err := query.Find(&products).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse(fmt.Sprintf("Failed to fetch products: %s", err.Error())),
		)
		return
	}

	c.JSON(http.StatusOK, products)
}

func ReadProductById(c *gin.Context) {
	// parse and validate id req
	id, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse(fmt.Sprintf("Invalid id: %s", err.Error())),
		)
		return
	}

	var product models.Product
	err = repositories.DB.First(&product, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(
			http.StatusNotFound,
			models.NewFailedResponse("Product not found"),
		)
		return
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse(fmt.Sprintf("Failed to get product: %s", err.Error())),
		)
		return
	}

	var productDto models.ProductDto
	productDto.FillFromModel(product)

	c.JSON(
		http.StatusOK,
		models.NewSuccessResponse("Success", productDto),
	)
}

func UpdateProductById(c *gin.Context) {
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
	var productDto models.ProductDto
	if err := c.ShouldBind(&productDto); err != nil {
		// response to bad input
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse(fmt.Sprintf("Invalid request payload: %s", err.Error())),
		)
		return
	}

	// fetch if product exist
	var existingProduct models.Product
	err = repositories.DB.First(&existingProduct, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				models.NewFailedResponse(fmt.Sprintf("Product not found: %s", err.Error())),
			)
		} else {
			c.JSON(
				http.StatusInternalServerError,
				models.NewFailedResponse(fmt.Sprintf("Failed to fetch product: %s", err.Error())),
			)
		}
		return
	}

	product := productDto.ToModel()
	product.ProductID = existingProduct.ProductID
	product.CreatedAt = existingProduct.CreatedAt
	product.UpdatedAt = time.Now()

	if err := repositories.DB.Save(&product).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse("Failed to update product"),
		)
		return
	}

	productDto.FillFromModel(product)
	c.JSON(
		http.StatusOK,
		models.NewSuccessResponse("Product successfully updated", productDto),
	)
}

func DeleteProductById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse(fmt.Sprintf("Invalid id: %s", err.Error())),
		)
	}

	result := repositories.DB.Delete(&models.Product{}, id)
	if result.Error != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse(fmt.Sprintf("Failed to delete product: %s", result.Error.Error())),
		)
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(
			http.StatusNotFound,
			models.NewFailedResponse("Product not found"),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		models.NewSuccessResponse(fmt.Sprintf("Product with ID %d deleted", id), nil),
	)
}

var productUploadDir = "./uploads/products"

func UploadProductImage(c *gin.Context) {
	// check if the productId valid
	productId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse(fmt.Sprintf("Invalid id: %s", err.Error())),
		)
		return
	}

	// check if the productId exist in DB
	var product models.Product
	if err := repositories.DB.First(&product, productId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				models.NewFailedResponse("Product not found"),
			)
		} else {
			c.JSON(
				http.StatusInternalServerError,
				models.NewFailedResponse(fmt.Sprintf("Database error: %s", err.Error())),
			)
		}
		return
	}

	// retrieve the uploaded file
	formFile, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse(fmt.Sprintf("Failed to retrieve uploaded file: %s", err.Error())),
		)
		return
	}
	defer formFile.Close()

	// validate file size max 1 mb = 2^20 = 1,048,576 bytes
	const maxFileSize = 1 << 20
	if fileHeader.Size > maxFileSize {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse("File size exceeds the maximum size 1 MB"),
		)
		return
	}

	// validate extension type and MIME type
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	validExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !validExtensions[ext] {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse("Only supports file types: .jpg, .jpeg, .png"),
		)
		return
	}

	buffer := make([]byte, 512)
	if _, err := formFile.Read(buffer); err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse("Failed to read file content"),
		)
		return
	}

	if mimeType := http.DetectContentType(buffer); !strings.HasPrefix(mimeType, "image/") {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse("Uploaded file is not an image"),
		)
		return
	}

	// upload to an extract dir or make new if there is no existing dr
	if err := os.MkdirAll(productUploadDir, os.ModePerm); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse(fmt.Sprintf("Failed to create directory folder: %s", err.Error())),
		)
	}

	// generate unique file name using timestamp to prevent replacing existing image uploaded
	baseName := strings.TrimSuffix(fileHeader.Filename, ext)
	timestamp := time.Now().Unix()
	savePath := filepath.Join(productUploadDir, fmt.Sprintf("%d_%s%s", timestamp, baseName, ext))
	savePath = strings.ReplaceAll(savePath, "\\", "/")

	// save image to server
	if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse("Failed to save uploaded file"),
		)
		return
	}

	// update product image_path
	product.ImagePath = sql.NullString{
		String: savePath,
		Valid:  true,
	}

	if err := repositories.DB.Save(&product).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewFailedResponse("Failed to update product image"),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		models.NewSuccessResponse("File successfully uploaded", map[string]string{
			"image_path": savePath,
		}))
}

func DownloadProductImage(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewFailedResponse("Invalid product_id"),
		)
		return
	}

	// check if product exist
	var product models.Product
	if err := repositories.DB.First(&product, productId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				models.NewFailedResponse("Product not found"),
			)
		} else {
			c.JSON(
				http.StatusInternalServerError,
				models.NewFailedResponse(fmt.Sprintf("Database error: %s", err.Error())),
			)
		}
		return
	}

	// check if the product has an image path
	if !product.ImagePath.Valid {
		c.JSON(
			http.StatusNotFound,
			models.NewFailedResponse("Product image not found"),
		)
		return
	}

	imagePath := product.ImagePath.String
	c.File(imagePath)
}

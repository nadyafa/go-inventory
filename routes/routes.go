package routes

import (
	"inventory-management/controllers"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	// products
	r.POST("/products", controllers.AddProduct)
	r.GET("/products", controllers.ReadProducts)
	r.GET("/products/:product_id", controllers.ReadProductById)
	r.PUT("/products/:product_id", controllers.UpdateProductById)
	r.DELETE("/products/:product_id", controllers.DeleteProductById)

	r.POST("/products/:product_id/productImage", controllers.UploadProductImage)
	r.GET("/products/:product_id/productImage", controllers.DownloadProductImage)

	// inventory
	r.GET("/inventories", controllers.ReadInventories)
	r.GET("/inventories/:product_id", controllers.ReadInventoryByProductID)
	r.PUT("/inventories/:product_id", controllers.UpdateInventoryByProductID)

	// order
	r.POST("/orders", controllers.CreateOrder)
	r.GET("/orders/:order_id", controllers.ReadOrderById)

	return r
}

package main

import (
	repository "inventory-management/repositories"
	"inventory-management/routes"
	"log"
)

func main() {
	repository.Init()

	r := routes.Router()

	// start route
	err := r.Run(":8080")
	if err != nil {
		log.Fatalln(err)
		return
	}
}

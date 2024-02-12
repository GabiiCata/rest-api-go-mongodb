package main

import (
	"log"
	

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "github.com/api-rest-go/docs" // Usa el nombre de tu módulo aquí


	"github.com/joho/godotenv"


	//import pkg/db folder
	"github.com/api-rest-go/pkg/db"
	"github.com/api-rest-go/internal/handlers"

)


func main() {

	// Cargar el archivo .env
	// especificar ubicacion del archivo .env en caso de que no se encuentre en la raiz del proyecto
	// godotenv.Load("/path/to/your/.env")
	if err := godotenv.Load("./internal/env/.env"); err != nil {
		log.Print("No .env file found")
	}
	_, err := db.Connect()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	router := gin.Default()

	// Swagger endpoint
    url := ginSwagger.URL("http://localhost:3001/swagger/doc.json") // The url pointing to API definition
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))


	router.GET("/person", handlers.GetPerson)
	router.POST("/person", handlers.PostPerson)
	router.GET("/person/:id", handlers.GetPersonByID)
	router.Run(":3001")
}

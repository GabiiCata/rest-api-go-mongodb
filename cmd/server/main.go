package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "github.com/api-rest-go/docs" // Usa el nombre de tu módulo aquí


	"github.com/joho/godotenv"


	//import pkg/db folder
	"github.com/api-rest-go/pkg/db"

)

type Person struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Age  int    `json:"age" bson:"age"`
}



// getPerson godoc
// @Summary Lista todas las personas
// @Description Obtiene un listado de todas las personas
// @Tags personas
// @Accept  json
// @Produce  json
// @Success 200 {array} Person
// @Router /person [get]
func getPerson(c *gin.Context) {
	var people []*db.Person

	// find all people with db.FindAll
	people = db.FindAll()

	
	

	c.IndentedJSON(http.StatusOK, people)
}

// postPerson godoc
// @Summary Agrega una nueva persona
// @Description Agrega una nueva persona al sistema
// @Tags personas
// @Accept  json
// @Produce  json
// @Param person body Person true "Persona a agregar"
// @Success 201 {object} Person
// @Router /person [post]
func postPerson(c *gin.Context) {
	collection := db.GetCollections()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var newPerson Person
	if err := c.BindJSON(&newPerson); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	_, err := collection.InsertOne(ctx, newPerson)
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusCreated, newPerson)
}

// getPersonByID godoc
// @Summary Obtiene una persona por su ID
// @Description Obtiene detalles de una persona específica por su ID
// @Tags personas
// @Accept  json
// @Produce  json
// @Param id path string true "ID de la Persona"
// @Success 200 {object} Person
// @Failure 404 {object} map[string]interface{}
// @Router /person/{id} [get]
func getPersonByID(c *gin.Context) {
	
	id := c.Param("id")
	var person *db.Person
	// find person by id with db.FindByID
	person = db.FindByID(id)

	c.IndentedJSON(http.StatusOK, person)
}

func main() {

	// Cargar el archivo .env
	if err := godotenv.Load(); err != nil {
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


	router.GET("/person", getPerson)
	router.POST("/person", postPerson)
	router.GET("/person/:id", getPersonByID)
	router.Run(":3001")
}

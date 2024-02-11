package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "github.com/api-rest-go/docs" // Usa el nombre de tu módulo aquí


)

type Person struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Age  int    `json:"age" bson:"age"`
}

var client *mongo.Client

func connectToMongoDB() {
	var err error
	// Asegúrate de reemplazar "root" y "example" con tu usuario y contraseña reales.
	// Además, reemplaza "mongo" con "localhost" si estás ejecutando este código fuera de Docker,
	// o mantenlo como "mongo" si este código se ejecuta dentro de otro contenedor de Docker
	// que forma parte de la misma red que tu contenedor de MongoDB.
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://root:example@mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
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
	collection := client.Database("apirestgo").Collection("people")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var people []Person
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &people); err != nil {
		log.Fatal(err)
	}

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
	collection := client.Database("apirestgo").Collection("people")
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
	collection := client.Database("apirestgo").Collection("people")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	var person Person
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&person)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "person not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, person)
}

func main() {
	connectToMongoDB()
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()


	

	router := gin.Default()


	// Swagger endpoint
    url := ginSwagger.URL("http://localhost:3001/swagger/doc.json") // The url pointing to API definition
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))


	router.GET("/person", getPerson)
	router.POST("/person", postPerson)
	router.GET("/person/:id", getPersonByID)
	router.Run(":3001")
}

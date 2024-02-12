package handlers

import (
	"net/http"
	// importar el paquete de models para usar el tipo Person
	"github.com/api-rest-go/internal/models"
	"github.com/api-rest-go/pkg/db"

	"github.com/gin-gonic/gin"
	"time"
	"context"
	"log"

)


// getPerson godoc
// @Summary Lista todas las personas
// @Description Obtiene un listado de todas las personas
// @Tags personas
// @Accept  json
// @Produce  json
// @Success 200 {array} Person
// @Router /person [get]
func GetPerson(c *gin.Context) {
	var people []*models.Person

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
func PostPerson(c *gin.Context) {
	collection := db.GetCollections()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var newPerson models.Person
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
func GetPersonByID(c *gin.Context) {
	
	id := c.Param("id")
	var person *models.Person
	// find person by id with db.FindByID
	person = db.FindByID(id)

	c.IndentedJSON(http.StatusOK, person)
}

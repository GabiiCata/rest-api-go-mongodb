
package models

type Person struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Age  int    `json:"age" bson:"age"`
}


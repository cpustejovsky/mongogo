package models

type User struct {
	ID     string `json:"id, omitempty" bson:"_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Age    int    `json:"age"`
	Active bool   `json:"active, omitempty"`
}

type FormUser struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Age   *int    `json:"age"`
}

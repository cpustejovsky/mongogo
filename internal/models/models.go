package models

type User struct {
	ID     string `json:"id, omitempty"`
	Name   string `json:"name, omitempty"`
	Email  string `json:"email, omitempty"`
	Age    int    `json:"age, omitempty"`
	Active bool   `json:"active, omitempty"`
}

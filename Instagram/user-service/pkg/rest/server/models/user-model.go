package models

type User struct {
	Id int64 `json:"id,omitempty"`

	Age int `json:"age,omitempty"`

	Name string `json:"name,omitempty"`
}

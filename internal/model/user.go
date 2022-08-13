package model

type User struct {
	UserID string
	Nick   string
	Email  string
	Secret string
	Kind   string
	Person
}

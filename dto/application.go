package dto

type Application struct {
	Name     string     `json:"name"`
	Instance []Instance `json:"instance"`
}

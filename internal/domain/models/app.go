package models

type App struct {
	Id   int
	Name string
	// Very bad!
	Secret string
}

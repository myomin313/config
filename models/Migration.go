package models

var allModels []interface{}

func RegisterModel(m interface{}) {
	allModels = append(allModels, m)
}

func AllModels() []interface{} {
	return allModels
}

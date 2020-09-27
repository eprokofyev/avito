package models

type Response struct {
	Status int `json:"status"`
	Body map[string]interface{} `json:"body"`
}

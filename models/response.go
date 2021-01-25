package models

// Response data api
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Count   int         `json:"count"`
}

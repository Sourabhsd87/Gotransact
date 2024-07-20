package models

type Response struct {
	Status  int                    `json:"status"`
	Message string                 `json:"messsage"`
	Data    map[string]interface{} `json:"data"`
}

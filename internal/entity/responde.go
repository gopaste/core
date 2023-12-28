package entity

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Info    interface{} `json:"info,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

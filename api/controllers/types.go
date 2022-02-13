package controllers

type HTTPResponse struct {
	Code    int    `json:"code"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

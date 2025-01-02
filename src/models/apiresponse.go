package models

import "net/http"

type ApiResponse struct {
	Response   bool
	StatusCode int
	Headers    http.Header
	Message    string
	Body       []byte
}

func NewApiResponse(response bool, statusCode int, headers http.Header, message string, body []byte) ApiResponse {
	return ApiResponse{
		Response:   response,
		StatusCode: statusCode,
		Headers:    headers,
		Message:    message,
		Body:       body,
	}
}

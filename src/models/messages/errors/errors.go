package errors

type ErrorResponse struct {
    Response bool        `json:"response"`
    Code     string      `json:"code"`
    Message  string      `json:"message"`
    Data     interface{} `json:"data"`
}

func NewErrorResponse(response bool, code string, customMessage string) ErrorResponse {
    return ErrorResponse{
        Response: response,
        Code:     code,
        Message:  customMessage,
        Data:     struct{}{}, // An empty struct results in {} in JSON
    }
}

package response

type Response struct {
    Response bool        `json:"response"`
    Code     string      `json:"code"`
    Message  string      `json:"message"`
    Data     interface{} `json:"data"`
}

func NewResponse(response bool, code, message string, data interface{}) Response {
    return Response{
        Response: response,
        Code:     code,
        Message:  message,
        Data:     data,
    }
}
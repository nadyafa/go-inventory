package models

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func NewSuccessResponse(message string, data any) Response {
	return Response{
		Success: true,
		Data:    data,
		Message: message,
	}
}

func NewFailedResponse(message string) Response {
	return Response{
		Message: message,
	}
}

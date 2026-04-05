package errors

type FieldErrors map[string]string

type AppError struct {
	Status  int
	Code    string
	Message string
	Fields  FieldErrors
}

type ErrorBody struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Fields    FieldErrors `json:"fields,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

func (e *AppError) Response(requestID string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorBody{
			Code:      e.Code,
			Message:   e.Message,
			Fields:    e.Fields,
			RequestID: requestID,
		},
	}
}

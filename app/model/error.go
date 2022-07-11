package model

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type FieldErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ToFieldErrorsResponse(fieldErrors *[]FieldError) []FieldErrorResponse {
	res := make([]FieldErrorResponse, 0)
	for _, fieldError := range *fieldErrors {
		res = append(res, FieldErrorResponse{
			Field:   fieldError.Field,
			Message: fieldError.Message,
		})
	}

	return res
}

package errors

type ErrorResponseObject struct {
	Error    string
	Messages []string
}

func New() ErrorResponse {
	return &ErrorResponseObject{}
}

func NewWithDetails(error string, messages ...string) ErrorResponse {
	return &ErrorResponseObject{
		Error:    error,
		Messages: messages,
	}
}
func NewWithError(err error, messages ...string) ErrorResponse {
	return &ErrorResponseObject{
		Error:    err.Error(),
		Messages: messages,
	}
}

type ErrorResponse interface {
	AddError(error string)
	AddErrorIfMessagesExists(error string)
	AddMessage(message string)
	Result() ErrorResponse
}

func (e *ErrorResponseObject) AddMessage(message string) {
	e.Messages = append(e.Messages, message)
}

func (e *ErrorResponseObject) AddError(error string) {
	e.Error = error
}

func (e *ErrorResponseObject) AddErrorIfMessagesExists(error string) {
	if len(e.Messages) != 0 {
		e.Error = error
	}
}

func (e *ErrorResponseObject) Result() ErrorResponse {
	if e.Error == "" && len(e.Messages) == 0 {
		return nil
	}
	return e
}

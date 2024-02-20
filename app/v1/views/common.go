package views

type HTTPResponseVM struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func NewHTTPResponseVM(message string, data interface{}) HTTPResponseVM {
	return HTTPResponseVM{
		Message: message,
		Data:    data,
	}
}

func NewHTTPResponseVMFromError(err error) HTTPResponseVM {
	return HTTPResponseVM{
		Message: err.Error(),
		Data:    nil,
	}
}

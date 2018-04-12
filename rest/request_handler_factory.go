package rest

type RequestHandlerFactory struct{
	errorHandlers []ErrorHandler
}

func (rhf *RequestHandlerFactory) AddErrorHandler(handler ErrorHandler){
	rhf.errorHandlers = append(rhf.errorHandlers, handler)
}

func (rhf *RequestHandlerFactory) New(handlerFunction RequestHandlerFunc) RequestHandler{
	handler := NewRequestHandler(handlerFunction)
	handler.SetErrorHandlers(rhf.errorHandlers)
	return handler
}
package rest

import (
	"net/http"
	"authentication/service"
	"strings"
	"github.com/satori/go.uuid"
)

type RequestHandlerFunc func(userId uuid.UUID, w http.ResponseWriter, r *http.Request) error
type ErrorHandler func(w http.ResponseWriter, err error) bool

type RequestHandler struct
{
	handlerInternal RequestHandlerFunc
	errorHandlers []ErrorHandler
	auth *authentication.Authentication
}

func NewRequestHandler(handlerFunction RequestHandlerFunc) RequestHandler{
	var handler RequestHandler
	handler.handlerInternal = handlerFunction
	handler.auth = authentication.GetServiceInstance()
	return handler
}

func (rh RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	userID, _ := rh.Authenticate(r)
	err:= rh.handlerInternal(userID, w, r)
	rh.HandleError(w, err)
}

func (rh RequestHandler) Authenticate(r *http.Request) (uuid.UUID, error){
	if val, ok := r.Header["Authorization"]; ok {
		if strings.HasPrefix(val[0], "le v1="){
			authentication := strings.Split(val[0], "=")[1]
			userID := rh.auth.CheckTicket(authentication)
			if userID != uuid.Nil {
				return userID, nil
			}
		}
	}
	return uuid.Nil, authentication.AuthError{}
}

func (rh *RequestHandler) AddErrorHandler(handler ErrorHandler){
	rh.errorHandlers = append(rh.errorHandlers, handler)
}

func (rh *RequestHandler) SetErrorHandlers(handlers []ErrorHandler){
	rh.errorHandlers = handlers
}

func (rh *RequestHandler) HandleError(w http.ResponseWriter, err error){
	if err != nil {
		for _, handler := range rh.errorHandlers {
			if handler(w, err){
				return
			}
		}
		
		switch e := err.(type) {
		case authentication.AuthError:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			http.Error(w, e.Error(), 403)
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}
package main
import "promise/service"
import "github.com/gorilla/mux"
import "log"
import "net/http"
import "strconv"
import "authentication/service"
import "rest"
import "github.com/satori/go.uuid"
import "encoding/json"
import "io/ioutil"

type RestHandler struct{
	service *promises.Promise
	auth *authentication.Authentication
}

func CreateRestHandler() *RestHandler{
	handler := new(RestHandler)
	handler.service = promises.CreatePromiseService()

	return handler
}

func (f *RestHandler) StartListening(port int){
	requestHandlerFactory := rest.RequestHandlerFactory{}
	requestHandlerFactory.AddErrorHandler(f.ErrorHandler)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/promise", requestHandlerFactory.New(f.AddPromise).ServeHTTP).Methods("POST")
	router.HandleFunc("/promise/{id:[0-9]+}", requestHandlerFactory.New(f.PromiseInteract).ServeHTTP).Methods("GET", "PUT", "DELETE")
	
    log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), router))
}

func (f *RestHandler) PromiseInteract(userID uuid.UUID, w http.ResponseWriter, r *http.Request) error {
	switch r.Method{
	case http.MethodGet:
		return f.GetPromise(userID, w, r)
	case http.MethodPut:
		return f.EditPromise(userID, w, r)
	case http.MethodDelete:
		return f.DeletePromise(userID, w, r)
	}
	return nil
}

func (f *RestHandler) AddPromise(userID uuid.UUID, w http.ResponseWriter, r *http.Request) error {
	var promise PromiseAdd
	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, &promise)
	if err != nil{
		return err
	}
	f.service.AddPromise(userID, promise.Title, promise.Description)
	return nil
}

func (f *RestHandler) GetPromise(userID uuid.UUID, w http.ResponseWriter, r *http.Request) error {
	id,_ := strconv.Atoi(mux.Vars(r)["id"])
	promise := f.service.GetPromise(id)
	jsonBody, _ := json.Marshal(promise)
	w.Write(jsonBody)
	return nil
}

func (f *RestHandler) EditPromise(userID uuid.UUID, w http.ResponseWriter, r *http.Request) error {
	id,_ := strconv.Atoi(mux.Vars(r)["id"])
	var promise promises.LEPromise
	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, &promise)
	if err != nil{
		return err
	}
	if userID != promise.UserId{
		return promises.InvalidUserError{}
	}
	f.service.EditPromise(id, promise.Title, promise.Description, promise.PromisedTo, promise.Status, promise.Privacy)
	return nil
}

func (f *RestHandler) DeletePromise(userID uuid.UUID, w http.ResponseWriter, r *http.Request) error {
	id,_ := strconv.Atoi(mux.Vars(r)["id"])
	promise := f.service.GetPromise(id)
	jsonBody, _ := json.Marshal(promise)
	if userID != promise.UserId {
		return promises.InvalidUserError{}
	}
	f.service.DeletePromise(promise.PromiseId)
	w.Write(jsonBody)
	return nil
}

func (f *RestHandler) ErrorHandler(w http.ResponseWriter, err error) bool{
	switch e := err.(type) {
	case promises.InvalidUserError:
		// We can retrieve the status here and write out a specific
		// HTTP status code.
		http.Error(w, e.Error(), 400)
		return true
	}

	return false
}
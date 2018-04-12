package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	authentication "github.com/donydony2009/Promise/scripts/authentication/service"
	"github.com/gorilla/mux"
)

type RestHandler struct {
	service *authentication.Authentication
}

func CreateRestHandler() *RestHandler {
	handler := new(RestHandler)
	handler.service = authentication.GetServiceInstance()
	return handler
}

func (f *RestHandler) StartListening(port int) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/auth/login", f.Login).
		Methods("GET")
	router.HandleFunc("/user/createAccount", f.CreateAccount).
		Methods("POST")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}

func (f *RestHandler) Login(w http.ResponseWriter, r *http.Request) {
	if val, ok := r.Header["Authorization"]; ok {
		if strings.HasPrefix(val[0], "ad v1=") {
			authBase64 := strings.Split(val[0], "=")[1]
			decodedAuth, _ := base64.StdEncoding.DecodeString(authBase64)
			authInfo := strings.Split(string(decodedAuth[:len(decodedAuth)]), ":")
			user := authInfo[0]
			password := authInfo[1]
			ticket := f.service.Login(user, password)
			fmt.Fprintf(w, `{"ticket":"%s"}`, ticket)
		}
	}

}

func (f *RestHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var createAccount CreateAccount
	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, &createAccount)
	if err != nil {
		panic(err)
	}
	f.service.CreateAccount(createAccount.User, createAccount.Password, createAccount.Email)
	fmt.Fprintf(w, `{"user":"%s", "password":"%s", "email":"%s"}`, createAccount.User, createAccount.Password, createAccount.Email)
}

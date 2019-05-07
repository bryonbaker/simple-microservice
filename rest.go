package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Request struct {
	Text string `json:"text,omitempty"`
}

type Response struct {
	Text string `json:"text,omitempty"`
}

var version string = "1.0"

func HomeHandlerWithKey(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	json.NewEncoder(w).Encode(version)
}

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(version)
}

/*
func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    var person Person
    _ = json.NewDecoder(req.Body).Decode(&person)
    person.ID = params["id"]
    people = append(people, person)
    json.NewEncoder(w).Encode(people)
}

func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    for index, item := range people {
        if item.ID == params["id"] {
            people = append(people[:index], people[index+1:]...)
            break
        }
    }
    json.NewEncoder(w).Encode(people)
}*/

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/{key}", HomeHandlerWithKey).Methods("GET")
	// router.HandleFunc("/simple", SimpleRequest).Methods("GET")
	// router.HandleFunc("/simple/{id}", SimpleRequestWithId).Methods("GET")
	// router.HandleFunc("/simple/{id}", CreatePersonEndpoint).Methods("POST")
	// router.HandleFunc("/simple/{id}", DeletePersonEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))
}

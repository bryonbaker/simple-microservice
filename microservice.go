package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Request struct {
	ClientId string `json:"clientId,omitempty"`
}

type Response struct {
	ServiceVersion string `json:"serviceVersion,omitempty"`
}

var version string = "1.2"

/// This is the simple request handler that takes no onput parameters.
func HomeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("HomeHandler()")
	json.NewEncoder(w).Encode(version)
}

/// Handle the request that has an input parameter of the client ID.
func HomeHandlerWithKey(w http.ResponseWriter, req *http.Request) {
	var unknownList string = ""
	params := mux.Vars(req)
	// If clientId is in the parameters, process the request.
	if parm, ok := params["clientId"]; ok {
		fmt.Println("Request received with key: ", params[parm])
	} else {
		unknownList += parm + " "
	}

	// Dump out the list of unknown input params
	if unknownList != "" {
		fmt.Println("Unknown input parameter: ", unknownList)
		fmt.Fprintf(os.Stderr, "Unknown input parameter: ", unknownList)
	}

	// Build the response.
	json.NewEncoder(w).Encode(version)
}

func main() {
	fmt.Println("Simple microservice starting. Version: ", version)

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/{key}", HomeHandlerWithKey).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", router))
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

/*
type Request struct {
	ClientId string `json:"clientId,omitempty"`
}
*/

// The response to retutn for the HTML GET
type Response struct {
	ServiceVersion string `json:"serviceVersion,omitempty"`
}

// The version of this release.
var version string = "1.3c"

// Builds the json response string.
func BuildResponse() Response {
	var resp Response

	resp.ServiceVersion = version

	return resp
}

/// This is the simple request handler that takes no onput parameters.
func HomeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("HomeHandler()")
	json.NewEncoder(w).Encode(BuildResponse())
}

/// Handle the request that has an input parameter of the client ID.
func HomeHandlerWithKey(w http.ResponseWriter, req *http.Request) {
	var unknownList string = ""
	params := mux.Vars(req)
	// If clientId is in the parameters, process the request.
	if parm, ok := params["key"]; ok {
		fmt.Println("Request received with key: ", parm)
	} else {
		unknownList += parm + " "
	}

	// Dump out the list of unknown input params
	if unknownList != "" {
		fmt.Println("Unknown input parameter: ", unknownList)
		fmt.Fprintf(os.Stderr, "Unknown input parameter: ", unknownList)
	}

	// Build the response.
	json.NewEncoder(w).Encode(BuildResponse())
}

func main() {
	fmt.Println("Simple microservice starting. Version: ", version)

	router := mux.NewRouter()
	corsMw := mux.CORSMethodMiddleware(router)

	router.HandleFunc("/", HomeHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/{key}", HomeHandlerWithKey).Methods("GET", "OPTIONS")

	router.Use(corsMw)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe(":10000", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
	// log.Fatal(http.ListenAndServe(":10000", router))
}

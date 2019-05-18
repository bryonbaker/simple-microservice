package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

/*
type Request struct {
	ClientId string `json:"clientId,omitempty"`
}
*/

// Configuration holds the configuration loaded at startup
type Configuration struct {
	Version string `json:"version,omitempty"`
}

// Response defines the format of the HTML json response.
type Response struct {
	ServiceVersion string `json:"version,omitempty"`
}

// ConfigurationFile contains the name of the file containing the applicaiton configuration
var ConfigurationFile = "config.json"

// ServiceVersion keeps the service version handy
var ServiceVersion string = "1.0(default)"

// Setup is called on startup to initialise all valiables etc.
func setup() Configuration {
	// Read a config file from storage.
	config, err := loadConfig(ConfigurationFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Cannot load configuration.")
		panic(err)
	}

	// If the file is not there then use the internal default version.

	return config
}

// BuildResponse is used to generate theresponse.
func buildResponse() Response {
	var resp Response

	resp.ServiceVersion = ServiceVersion

	return resp
}

// HomeHandler is the simple request handler that takes no onput parameters.
func homeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("HomeHandler()")
	json.NewEncoder(w).Encode(buildResponse())
}

// HomeHandlerWithKey is a request that has an input parameter of the client ID.
func homeHandlerWithKey(w http.ResponseWriter, req *http.Request) {
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
		msg := "ERROR: Unknown input parameter: " + unknownList
		fmt.Println(msg)
		fmt.Fprintf(os.Stderr, msg)
	}

	// Build the response.
	json.NewEncoder(w).Encode(buildResponse())
}

// LoadConfig loads the config file that the program uses. at startup. Returns error is the file cannot be loaded
func loadConfig(file string) (Configuration, error) {
	var config Configuration

	// Find the directory that this app is running in and use it as the base directory for the config file.
	configFileAbsPath := getExePath() + "/" + file
	fmt.Println("Loading configuration from: ", configFileAbsPath)

	jsonFile, err := os.Open(configFileAbsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Cannot find configuration file in: ", configFileAbsPath)
		return config, err
	}

	defer jsonFile.Close()
	byteVal, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Cannot read the json file contents.\n")
		return config, err
	}

	must(json.Unmarshal(byteVal, &config))

	fmt.Println("Loaded configuration file version: ", config.Version)

	return config, nil
}

func getExePath() string {
	var exePath string
	var err error

	exePath, err = os.Executable()
	if err != nil {
		panic(err)
	}
	fmt.Println("Exe path = ", exePath)

	// Get the file info to see if it is a symlink.
	fileInfo, err := os.Lstat(exePath)
	if err != nil {
		log.Fatal(err)
	}

	// Double check the path is not a symlink. If so, get the target.
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		// Get the path the symlink points to.

		exePath, err = filepath.EvalSymlinks(exePath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Resolved symlink to be: ", exePath)
	}

	// Get the directory from the full path - removing any trailing / etc
	dir := filepath.Dir(exePath)

	return dir
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Starting service...")
	config := setup()
	ServiceVersion = config.Version

	router := mux.NewRouter()
	corsMw := mux.CORSMethodMiddleware(router)

	router.HandleFunc("/", homeHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/{key}", homeHandlerWithKey).Methods("GET", "OPTIONS")

	router.Use(corsMw)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe(":10000", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

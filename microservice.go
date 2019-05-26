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

const (
	// DefaultServiceVersion keeps the service version handy
	DefaultServiceVersion string = "unknown"

	// BootConfigurationFile contains the name of the file containing the bootstrap configuratiopn.
	// It must exist in the same directory as the executable.
	BootConfigurationFile = "boot-config.json"
)

// BootConfig holds the startup config used to bootstrap all other config.
type BootConfig struct {
	ConfigPath     string `json:"config-path,omitempty"`
	VersionVar     string `json:"version-env-var-name"`
	DefaultVersion string `json:"default-version,omitempty"`
}

// ApplicationConfiguration holds all config info for the application.
type ApplicationConfiguration struct {
	ServiceVersion string `json:"service-version,omitempty"`
}

// Response defines the format of the HTML json response.
type Response struct {
	ServiceVersion string `json:"id,omitempty"`
}

// Global variable definitions.
var appConfig ApplicationConfiguration

// Setup is called on startup to initialise all valiables etc. The idea of this is that it is a cascading load of
// configuration that enables testing Kubernetes persistent volumes and config.
func setup() ApplicationConfiguration {
	var config = ApplicationConfiguration{DefaultServiceVersion}

	// Read the boot configuration file. The file must exist of the app terminates.
	bootconfig, err := LoadBootConfig(BootConfigurationFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL ERROR: Cannot load boot configuration.\n")
		panic(err)
	}
	// Load the application config from the file specified in boot config
	appconfig, err := LoadAppConfig(bootconfig.ConfigPath)
	if err == nil {
		config.ServiceVersion = appconfig.ServiceVersion
	} else {
		fmt.Fprintf(os.Stderr, "WARNING: Application config file does not exist. Trying environment variable.\n")

		// If the applicaiton config file does not exist load the environment variable specified int he boot config.
		env := os.Getenv(bootconfig.VersionVar)
		if env != "" {
			config.ServiceVersion = env
		} else {
			// If the environment variable does not exist use the version loaded in the boot config.
			fmt.Fprintf(os.Stderr, "WARNING: %s is not defined. Defaulting to boot configuration version.\n", bootconfig.VersionVar)
			config.ServiceVersion = bootconfig.DefaultVersion
		}
	}

	return config
}

// LoadBootConfig loads the config file that the program uses. at startup. Returns error is the file cannot be loaded
func LoadBootConfig(file string) (BootConfig, error) {
	var bootConfig BootConfig

	// Find the directory that this app is running in and use it as the base directory for the config file.
	configFileAbsPath := getExePath() + "/" + file
	fmt.Println("Loading configuration from: ", configFileAbsPath)

	jsonFile, err := os.Open(configFileAbsPath)
	if err != nil {
		var msg = "WARNING: Configuration file: " + configFileAbsPath + " does not exist\n"
		fmt.Fprintf(os.Stderr, msg)
		return bootConfig, err
	}

	defer jsonFile.Close()
	byteVal, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Cannot read the json file contents.\n")
		return bootConfig, err
	}

	must(json.Unmarshal(byteVal, &bootConfig))

	fmt.Println("Loaded boot configuration file version: ", bootConfig.DefaultVersion)

	return bootConfig, nil
}

// LoadAppConfig loads the applicaon's config from the specified path.
func LoadAppConfig(path string) (ApplicationConfiguration, error) {
	var config ApplicationConfiguration

	fmt.Println("Loading application configuration from: ", path)

	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Configuration file: %s does not exist\n", path)
		return config, err
	}

	defer jsonFile.Close()
	byteVal, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Cannot read the json file contents. Check JSON format.\n")
		return config, err
	}

	must(json.Unmarshal(byteVal, &config))

	return config, nil
}

// GetExePath retrieves the fully qualified of the executable. This needs to consider symlinks - so they
// are traversed to make sure we find the real executable.
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

// Must tests if error is nil and if not then the program terminates.s
func must(err error) {
	if err != nil {
		panic(err)
	}
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

// BuildResponse is used to generate theresponse.
func buildResponse() Response {
	var resp Response

	resp.ServiceVersion = appConfig.ServiceVersion

	return resp
}

func main() {
	fmt.Println("Initialising service...")
	appConfig = setup()
	fmt.Println("Starting service version: ", appConfig.ServiceVersion)

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

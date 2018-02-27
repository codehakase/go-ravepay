// Package types defines structs and constants used by the sdk
package types

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"runtime"
)

var (
	Envrionment string = "staging"
	Version     string = "1"
)

// RaveErr describes errors which are usually associated with the 400 http status code
type RaveErr struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

// ValidationErr a pointer to RaveErr, are returned when one or more validation rules fail
// Examples include not passing required parameters e.g.
// not passing the transaction / provider ref during a re-query call will result in the error below:
// {
//	"status":"error",
//	"message":"Cardno is required",
//	"data": null
// }
type ValidationErr *RaveErr

// Resources desbribes endpoints which is used in the sdk
type Resources struct {
	CardResources *CardResources
}

// CardResources is a struct holding the endpoints used by the CardService object
type CardResources struct {
	V1 struct {
		Staging    VersionEnv `json:"staging"`
		Production VersionEnv `json:"production"`
	} `json:"v1"`
	V2 struct {
		Staging    VersionEnv `json:"staging"`
		Production VersionEnv `json:"production"`
	} `json:"v2"`
}
type VersionEnv struct {
	Tokenize string `json:"tokenize"`
	Charge   string `json:"charge"`
	Validate string `json:"validate"`
	Preauth  string `json:"preauth"`
	Capture  string `json:"capture"`
	Refund   string `json:"refund"`
	AVS      string `json:"avs"`
	Status   string `json:"status"`
}

func LoadConfigs() *Resources {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	dir := path.Dir(filename)

	// load resource configs for service objects
	resourcesFile, err := os.Open(dir + "/resources.json")
	defer resourcesFile.Close()
	if err != nil {
		log.Fatalln(err)
	}

	parser := json.NewDecoder(resourcesFile)
	var cr CardResources
	if err = parser.Decode(&cr); err != nil {
		log.Fatalln(err)
	}

	d := &Resources{
		CardResources: &cr,
	}
	return d
}

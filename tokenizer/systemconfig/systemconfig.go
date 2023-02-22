package systemconfig

import (
	"encoding/json"
	"os"
	"strings"
)

var servicePathName = "tokenizerService"

type Configuration struct {
	TokenizerServiceUrl     string // tokenizerService
	TokenizerServiceApiMode string // tokenizerService
	CORSAllowOrigin         string // tokenizerService
	MongoUri                string // datastore
	MongoDatabase           string // datastore
	PageRecordCount         int64  // tokenizer
	EncryptionKey           string // tokenizer
}

func Load() (Configuration, error) {
	var configuration Configuration
	path, err := os.Getwd()
	if err != nil {
		return configuration, err
	}

	pathparts := strings.Split(strings.TrimSpace(path), string(os.PathSeparator))
	var configpath string = ""
	// remove the last part of the path, which is "tokenizerService"
	for _, v := range pathparts {
		if v == servicePathName {
			break
		}
		configpath += v + string(os.PathSeparator)
	}
	configuration, err = loadFromFile(configpath + "config.json")
	return configuration, err
}

func loadFromFile(filepath string) (Configuration, error) {
	var configuration Configuration
	file, err := os.Open(filepath)
	if err != nil {
		return configuration, err
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	return configuration, err
}

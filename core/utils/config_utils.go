package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(filepath.Clean(path))
}

func ReadYamlConfig(path string, o interface{}) error {
	b, err := ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, o)
	if err != nil {
		log.Fatalf("unable to parse file with error %#v", err)
	}

	return nil
}

func ReadJsonConfig(path string, o interface{}) error {
	b, err := ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, o)
	if err != nil {
		log.Fatalf("unable to parse file with error %#v", err)
	}

	return nil
}

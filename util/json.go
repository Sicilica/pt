package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ReadJSONFile is a convenience function to write a JSON file.
func ReadJSONFile(filename string, v interface{}) error {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, v)
}

// WriteJSONFile is a convenience function to write a JSON file.
func WriteJSONFile(filename string, v interface{}, perm os.FileMode) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, raw, perm)
}

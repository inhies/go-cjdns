// Package config allows easy loading, manipulation, and saving of cjdns configuration files.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
)

// Loads and parses the input file and returns a Config structure with the minimal cjdroute.conf file requirements.
func LoadMinConfig(filein string) (*Config, error) {

	//Load the raw JSON data from the file
	raw, err := loadJson(filein)
	if err != nil {
		return nil, err
	}

	//Parse the JSON in to our struct which supports all requried fields for cjdns
	structured, err := parseJSONStruct(raw)
	if err != nil {
		return nil, err
	}

	//Parse the JSON in to an object to preserve non-standard fields
	object, err := parseJSONObject(raw)
	if err != nil {
		return nil, err
	}

	//Parse the odd security section of the config
	for _, value := range object["security"].([]interface{}) {
		v := reflect.ValueOf(value)
		if value == "nofiles" {
			structured.Security.NoFiles = true
		} else if v.Kind() == reflect.Map {
			user := value.(map[string]interface{})
			structured.Security.SetUser = user["setuser"].(string)
		}
	}
	return &structured, nil
}

// Loads and parses the input file and returns a map with all data found in the config file, including non-standard fields.
func LoadExtConfig(filein string) (map[string]interface{}, error) {

	//Load the raw JSON data from the file
	raw, err := loadJson(filein)
	if err != nil {
		return nil, err
	}

	//Parse the JSON in to an object to preserve non-standard fields
	object, err := parseJSONObject(raw)
	if err != nil {
		return nil, err
	}

	return object, nil
}

// Saves either of the two config types to the specified file with the specified permissions.
func SaveConfig(fileout string, config interface{}, perms os.FileMode) error {

	//check to see if we got a struct or a map (minimal or extended config, respectively)
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Struct {
		config := config.(Config)

		//Parse the nicely formatted security section, and set the raw values for JSON marshalling
		newSecurity := make([]interface{}, 0)
		if config.Security.NoFiles {
			newSecurity = append(newSecurity, "nofiles")
		}
		setuser := make(map[string]interface{})
		setuser["setuser"] = config.Security.SetUser
		newSecurity = append(newSecurity, setuser)
		config.RawSecurity = newSecurity

		jsonout, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			return err
		}
		return ioutil.WriteFile(fileout, jsonout, perms)
	} else if v.Kind() == reflect.Map {
		jsonout, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			return err
		}
		return ioutil.WriteFile(fileout, jsonout, perms)
	}
	return fmt.Errorf("Something very bad happened")
}

// Returns a []byte of raw JSON with comments removed.
func loadJson(filein string) ([]byte, error) {
	file, err := ioutil.ReadFile(filein)
	if err != nil {
		return nil, err
	}

	raw, err := stripComments(file)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// Returns a Config structure with the JSON unmarshalled in to it.
func parseJSONStruct(jsonIn []byte) (Config, error) {
	var structured Config
	err := json.Unmarshal(jsonIn, &structured)
	if err != nil {
		return Config{}, err
	}
	return structured, nil
}

// Returns a map with the JSON unmarshalled in to it.
func parseJSONObject(jsonIn []byte) (map[string]interface{}, error) {
	var object map[string]interface{}
	err := json.Unmarshal(jsonIn, &object)
	if err != nil {
		return nil, err
	}
	return object, nil
}

// Replaces all C-style comments (prefixed with "//" and inside "/* */") with empty strings. This is necessary in parsing JSON files that contain them.
// Returns b without comments. Credit to SashaCrofter, thanks!
func stripComments(b []byte) ([]byte, error) {
	regComment, err := regexp.Compile("(?s)//.*?\n|/\\*.*?\\*/")
	if err != nil {
		return nil, err
	}
	out := regComment.ReplaceAllLiteral(b, nil)
	return out, nil
}

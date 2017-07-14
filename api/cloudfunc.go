package api

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"fmt"
)

// GetInput accepts an interface and unmarshalls the function input using json.Unmarshal
func GetInput(i interface{}) error {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, i)
}

// Send marshalls given interface using json.Marshal and sends it back as a function output
// If the interface can't be serialized, it will be returned as a pure string
func Send(i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	_, err = fmt.Println(string(b))

	return err
}

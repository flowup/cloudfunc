package api

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"fmt"
)

// GetInput unarshalls the function input from JSON
func GetInput(i interface{}) error {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, i)
}

// Send marshalls given structure into the outgoing JSON
func Send(i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	_, err = fmt.Println(string(b))

	return err
}

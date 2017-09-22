package api

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// CloudFunc represents communication interface
type CloudFunc struct {
	input  io.Reader
	output io.Writer
}

func NewCloudFunc() *CloudFunc {
	return &CloudFunc{input: os.Stdin, output: os.Stdout}
}

// GetRequest accepts an request and unmarshals it to Request
func (c *CloudFunc) GetRequest() (*Request, error) {
	b, err := ioutil.ReadAll(c.input)
	if err != nil {
		return nil, err
	}

	request := Request{}

	err = json.Unmarshal(b, &request)
	return &request, err
}

// SendResponse marshalls given interface using json.Marshal and sends it back as a function output
// If the interface can't be serialized, it will be returned as a pure string
func (c *CloudFunc) SendResponse(i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(c.output, string(b))

	return err
}

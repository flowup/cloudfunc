package api

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"encoding/json"
)

type RequestSuite struct {
	suite.Suite
}

func (s *RequestSuite) SetupSuite() {
}

func (s *RequestSuite) SetupTest() {

}

// Test
func (s *RequestSuite) TestRequest() {
	type testStruct struct {
		Test string `json:"test"`
	}

	request := Request{
		Body: json.RawMessage(`{"test":"param"}`),
	}

	var body testStruct

	err := request.BindBody(&body)
	s.Nil(err)
	s.Equal(body.Test, "param")
}

func TestRequestSuite(t *testing.T) {
	suite.Run(t, new(RequestSuite))
}

func (s *RequestSuite) TearDownSuite() {

}

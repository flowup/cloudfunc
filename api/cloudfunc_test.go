package api

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"bytes"
	"io/ioutil"
	"bufio"
)

type CloudFuncSuite struct {
	suite.Suite
	OutputBuff bytes.Buffer
	Function   CloudFunc
}

func (s *CloudFuncSuite) SetupSuite() {
	inputFile, err := ioutil.ReadFile("test-mock/cloudfunc.json")
	s.Nil(err)

	input := bytes.NewReader(inputFile)
	output := bufio.NewWriter(&s.OutputBuff)

	s.Function = CloudFunc{input: input, output: output}
}

func (s *CloudFuncSuite) SetupTest() {

}

// TestGetInput
func (s *CloudFuncSuite) TestGetInput() {
	req, err := s.Function.GetRequest()
	s.Nil(err)

	s.Equal(req.BaseURL, "something.com")
	s.True(bytes.Equal([]byte(req.Body), []byte(`{"glel":"asdd"}`)), "request bodies don't match")
	s.Equal(req.Hostname, "something.com")
	s.Equal(req.IP, "1.2.3.4")
	s.Equal(req.Method, "POST")
	s.Equal(req.OriginalURL, "url")
	s.Equal(req.Hostname, "something.com")

	key, ok := req.Query["param"]
	s.True(ok)
	s.Equal(key, "value")

	key, ok = req.Query["key"]
	s.True(ok)
	s.Equal(key, "nextvalue")
}

// TestSend
func (s *CloudFuncSuite) TestSend() {

}

func TestCloudFuncSuite(t *testing.T) {
	suite.Run(t, new(CloudFuncSuite))
}

func (s *CloudFuncSuite) TearDownSuite() {

}

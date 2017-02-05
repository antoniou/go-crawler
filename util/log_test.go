package util

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LogTestSuite struct {
	suite.Suite
}

func helper(r io.Reader, outC chan string) {
	var buf bytes.Buffer
	io.Copy(&buf, r)
	outC <- buf.String()
}

func (suite *LogTestSuite) TestVerboseLogger() {
	r, w, _ := os.Pipe()

	outC := make(chan string)
	go helper(r, outC)

	_ = Logger(true)
	log.SetOutput(w)
	Println("This should be printed")

	w.Close()
	out := <-outC
	assert.Contains(suite.T(), out, "This should be printed")
}

func (suite *LogTestSuite) TestNonVerboseLogger() {
	r, w, _ := os.Pipe()

	outC := make(chan string)
	go helper(r, outC)

	_ = Logger(false)
	log.SetOutput(w)
	Println("This should NOT be printed")

	w.Close()
	out := <-outC
	assert.NotContains(suite.T(), out, "This should NOT be printed")
}

func TestLogTestSuite(t *testing.T) {
	suite.Run(t, new(LogTestSuite))
}

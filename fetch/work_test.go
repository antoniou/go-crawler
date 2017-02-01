package fetch

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WorkTestSuite struct {
	suite.Suite
	e *Encapsulator
}

type Encapsulator struct {
	*AsyncWorker
	mock.Mock
}

func (e *Encapsulator) mockedMethod() error {
	e.Called()
	return nil
}

func (e *Encapsulator) mockedMethodReturnsError() error {
	return fmt.Errorf("Mock Error")
}

func (e *Encapsulator) mockedMethodStopChannel() error {
	select {
	case <-e.AsyncWorker.quit:
		fmt.Println("HERE")
		e.Called()
		return nil
	default:
		return nil
	}
	// return nil
}

func NewEncapsulator() *Encapsulator {
	e := &Encapsulator{
		AsyncWorker: &AsyncWorker{
			Name: "Encapsulator",
		},
	}
	e.AsyncWorker.RunFunc = e.mockedMethod
	return e
}

func (suite *WorkTestSuite) TestWorkerCallsEncapsulatorsRunMethod() {
	e := NewEncapsulator()
	e.On("mockedMethod").Return(nil)
	e.AsyncWorker.Run()
	e.AssertExpectations(suite.T())

	// When mocked method returns error, Run should return error
	e.AsyncWorker.RunFunc = e.mockedMethodReturnsError
	assert.Error(suite.T(), e.AsyncWorker.Run())

}

func (suite *WorkTestSuite) TestWorkerNameIsCorrect() {
	e := NewEncapsulator()
	assert.Equal(suite.T(), "Encapsulator", e.AsyncWorker.Type())
}

func (suite *WorkTestSuite) TestWorkerStateIsRight() {
	e := NewEncapsulator()

	// Test Initial state
	assert.Equal(suite.T(), WAITING, e.State())

	// Test after setting state
	e.SetState(RUNNING)
	assert.Equal(suite.T(), RUNNING, e.State())

}

// func (suite *WorkTestSuite) TestWorkerStopping() {
// 	e := NewEncapsulator()
// 	e.AsyncWorker.RunFunc = e.mockedMethodStopChannel
// 	e.AsyncWorker.Run()
// 	e.On("mockedMethodStopChannel").Return(nil)
//
// 	e.AsyncWorker.Stop()
// 	e.AssertExpectations(suite.T())
// }

func TestWorkTestSuite(t *testing.T) {
	suite.Run(t, new(WorkTestSuite))
}

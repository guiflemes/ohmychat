package core

import (
	"context"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"oh-my-chat/src/connector/telegram"
	"oh-my-chat/src/models"
)

func TestGetMultiChannelConnector(t *testing.T) {
	assert := assert.New(t)

	m := &multiChannelConnector{}
	m.connectors = Connectors
	type testCase struct {
		desc           string
		connName       models.MessageConnector
		expectedConnFn NewConnector
		expectedError  error
	}

	for _, c := range []testCase{
		{
			desc:           "get telegram connector",
			connName:       models.Telegram,
			expectedConnFn: telegram.NewTelegramConnector,
			expectedError:  nil,
		},
		{
			desc:           "get invalidConnector connector",
			connName:       models.Test,
			expectedConnFn: nil,
			expectedError:  &NotSupportConnectorError{},
		},
	} {
		t.Run(c.desc, func(t *testing.T) {

			fn, err := m.getConnector(c.connName)

			if fn != nil {
				funcName1 := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
				funcName2 := runtime.FuncForPC(reflect.ValueOf(c.expectedConnFn).Pointer()).Name()
				assert.Equal(funcName1, funcName2)
			}

			assert.Equal(c.expectedError, err)

		})
	}
}

type FakeConnector struct {
	dispatchCh chan models.Message
	output     *models.Message
	done       chan struct{}
}

func (f *FakeConnector) Acquire(ctx context.Context, input chan<- models.Message) {
	input <- models.Message{Input: "Hello test"}
}
func (f *FakeConnector) Dispatch(output models.Message) {
	f.output.Output = output.Input + " Output"
	f.done <- struct{}{}
}

type MultiChannelConnectorSuite struct {
	suite.Suite
	multiChannelConnector *multiChannelConnector
	output                *models.Message
	done                  chan struct{}
}

func (m *MultiChannelConnectorSuite) SetupTest() {
	m.output = &models.Message{}
	m.done = make(chan struct{})
	m.multiChannelConnector = &multiChannelConnector{
		connector: &FakeConnector{output: m.output, done: m.done},
	}
}

func (m *MultiChannelConnectorSuite) TestRequest() {
	context, cancel := context.WithCancel(context.Background())
	defer cancel()
	input := make(chan models.Message, 1)
	go m.multiChannelConnector.Request(context, input)

	r := <-input
	m.Equal(r.Input, "Hello test")

	close(input)

}

func (m *MultiChannelConnectorSuite) TestResponse() {
	output := make(chan models.Message, 1)
	defer close(output)
	ctx, cancel := context.WithCancel(context.Background())

	go m.multiChannelConnector.Response(ctx, output)
	output <- models.Message{Input: "Hello"}
	<-m.done
	m.Equal(m.output.Output, "Hello Output")
	cancel()
}

func TestMultiChannelConnectorSuite(t *testing.T) {
	suite.Run(t, new(MultiChannelConnectorSuite))
}

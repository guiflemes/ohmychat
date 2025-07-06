package core

// import (
// 	"context"
// 	"testing"

// 	"github.com/stretchr/testify/suite"

// 	"oh-my-chat/src/message"
// )

// type FakeConnector struct {
// 	dispatchCh chan message.Message
// 	output     *message.Message
// 	done       chan struct{}
// }

// func (f *FakeConnector) Acquire(ctx context.Context, input chan<- message.Message) {
// 	input <- message.Message{Input: "Hello test"}
// }
// func (f *FakeConnector) Dispatch(output message.Message) {
// 	f.output.Output = output.Input + " Output"
// 	f.done <- struct{}{}
// }

// type MultiChannelConnectorSuite struct {
// 	suite.Suite
// 	multiChannelConnector *multiChannelConnector
// 	output                *message.Message
// 	done                  chan struct{}
// }

// func (m *MultiChannelConnectorSuite) SetupTest() {
// 	m.output = &message.Message{}
// 	m.done = make(chan struct{})
// 	m.multiChannelConnector = &multiChannelConnector{
// 		connector: &FakeConnector{output: m.output, done: m.done},
// 	}
// }

// func (m *MultiChannelConnectorSuite) TestRequest() {
// 	context, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	input := make(chan message.Message, 1)
// 	go m.multiChannelConnector.Request(context, input)

// 	r := <-input
// 	m.Equal(r.Input, "Hello test")

// 	close(input)

// }

// func (m *MultiChannelConnectorSuite) TestResponse() {
// 	output := make(chan message.Message, 1)
// 	defer close(output)
// 	ctx, cancel := context.WithCancel(context.Background())

// 	go m.multiChannelConnector.Response(ctx, output)
// 	output <- message.Message{Input: "Hello"}
// 	<-m.done
// 	m.Equal(m.output.Output, "Hello Output")
// 	cancel()
// }

// func TestMultiChannelConnectorSuite(t *testing.T) {
// 	suite.Run(t, new(MultiChannelConnectorSuite))
// }

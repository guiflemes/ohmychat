package rule_engine

import (
	"github.com/guiflemes/ohmychat"
	"github.com/guiflemes/ohmychat/utils"
)

// Flow Execution Overview
// The Build() method links all the defined steps into a singly linked list, but no steps are executed during build time.

// Initialization:
// The flow starts when a user input matches a rule and triggers the root step (root.OnReply).
// This root step handles the input by sending a response to the user and then setting the session state to the next step.

// Session-based Execution:
// The next step is now stored in the session state and waits for the next user input.
// Once the user responds, the session state is resolved, and its associated handler is executed.

// Step Transition:
// Inside this handler, the step processes the current input and checks if a Next step is defined.
// If there is a next step, it is responsible for replying to the user (e.g., with a new question or options) and updating the session state again to point to the following step.

// Repeat:
// This pattern repeats: each step handles the input, responds, and schedules the next one.
// As a result, the flow progresses linearly, step by step, with user input driving each transition.

type FlowBuilder struct {
	root *FlowStep
	tail *FlowStep
}

type FlowStep struct {
	Prompt    string
	OnReply   ohmychat.ActionFunc
	Next      *FlowStep
	NextState ohmychat.SessionState
}

func NewFlow() *FlowBuilder {
	return &FlowBuilder{}
}

func (f *FlowBuilder) linkSteps(step *FlowStep) {
	if f.root == nil {
		f.root = step
		f.tail = step
		return
	}

	f.tail.Next = step
	f.tail = step
}

func (f *FlowBuilder) AskChoice(prompt string, options []string, handler ohmychat.ActionFunc) *FlowBuilder {
	step := &FlowStep{}
	step.OnReply = func(ctx *ohmychat.Context, msg *ohmychat.Message) {
		msg.ResponseType = ohmychat.OptionResponse
		msg.Output = prompt
		msg.Options = utils.OptionsFromList(options)
		ctx.SendOutput(msg)

		step.NextState = ohmychat.WaitingChoiceState{
			Choices: ohmychat.Choices{}.BindMany(func(c *ohmychat.Context, m *ohmychat.Message) {
				handler(c, m)
				if step.Next != nil {
					step.Next.OnReply(c, m)
					c.SetSessionState(step.Next.NextState)
				}
			}, options...),
		}
		ctx.SetSessionState(step.NextState)
	}

	f.linkSteps(step)
	return f
}

func (f *FlowBuilder) ThenAsk(prompt string, handler ohmychat.ActionFunc) *FlowBuilder {
	step := &FlowStep{}
	step.OnReply = func(ctx *ohmychat.Context, msg *ohmychat.Message) {
		msg.ResponseType = ohmychat.TextResponse
		msg.Output = prompt
		ctx.SendOutput(msg)

		step.NextState = ohmychat.WaitingInputState{
			Action: func(c *ohmychat.Context, m *ohmychat.Message) {
				handler(c, m)
				if step.Next != nil {
					step.Next.OnReply(c, m)
					c.SetSessionState(step.Next.NextState)
				}
			},
		}
		ctx.SetSessionState(step.NextState)

	}
	f.linkSteps(step)
	return f
}

func (f *FlowBuilder) ThenFinal(handler ohmychat.ActionFunc) *FlowBuilder {
	step := &FlowStep{OnReply: handler, NextState: ohmychat.IdleState{}}
	f.linkSteps(step)
	return f
}

func (f *FlowBuilder) ThenSayAndWait(prompt string) *FlowBuilder {
	step := &FlowStep{}

	step.OnReply = func(ctx *ohmychat.Context, msg *ohmychat.Message) {
		msg.Output = prompt
		ctx.SendOutput(msg)
		step.NextState = ohmychat.WaitingInputState{
			Action: func(c *ohmychat.Context, m *ohmychat.Message) {
				if step.Next != nil {
					step.Next.OnReply(c, m)
					c.SetSessionState(step.Next.NextState)
				}
			},
		}
		ctx.SetSessionState(step.NextState)
	}

	f.linkSteps(step)
	return f
}

func (f *FlowBuilder) ThenSayAndContinue(prompt string) *FlowBuilder {
	step := &FlowStep{}

	step.OnReply = func(ctx *ohmychat.Context, msg *ohmychat.Message) {
		msg.Output = prompt
		msg.BotMode = true
		ctx.SendOutput(msg)
		step.NextState = ohmychat.WaitingBotResponseState{
			OnDone: func(c *ohmychat.Context, m *ohmychat.Message) {
				if step.Next != nil {
					step.Next.OnReply(c, m)
					c.SetSessionState(step.Next.NextState)
				}
			},
		}
		ctx.SetSessionState(step.NextState)
	}

	f.linkSteps(step)
	return f
}

func (f *FlowBuilder) Start() ohmychat.ActionFunc {
	return func(ctx *ohmychat.Context, msg *ohmychat.Message) {
		if f.root == nil {
			msg.Output = "invalid flow, no steps"
			ctx.SendOutput(msg)
			return
		}
		f.root.OnReply(ctx, msg)
		ctx.SetSessionState(f.root.NextState)
	}
}

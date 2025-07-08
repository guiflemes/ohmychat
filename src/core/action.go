package core

import (
	"oh-my-chat/src/message"
)

type ActionFunc func(ctx *Context, msg *message.Message)

func WithValidation(
	validate func(input string) bool,
	errorMsg string,
	action ActionFunc,
) ActionFunc {
	return func(ctx *Context, msg *message.Message) {
		if !validate(msg.Input) {
			msg.Output = errorMsg
			ctx.SendOutput(msg)
			return
		}
		ctx.Session().State = IdleState{}
		action(ctx, msg)
	}
}

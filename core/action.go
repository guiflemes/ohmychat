package core

import (
	"github.com/guiflemes/ohmychat/message"
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
		ctx.SetSessionState(IdleState{})
		action(ctx, msg)
	}
}

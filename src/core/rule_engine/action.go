package rule_engine

import "context"

func WithValidation(
	validate func(input string) bool,
	errorMsg string,
	action ActionFunc,
) ActionFunc {
	return func(ctx context.Context, input ActionInput) {
		if !validate(input.Message.Input) {
			input.Message.Output = errorMsg
			input.Output <- *input.Message
			return
		}
		input.Session.State = IdleState{}
		action(ctx, input)
	}
}

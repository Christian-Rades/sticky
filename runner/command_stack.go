package runner

import "context"


type CommandStack struct {
	Commands []Command
}

type Command interface {
	Run(context.Context, *CommandStack) error
}

func (stack *CommandStack) Add(c Command) {
	stack.Commands = append(stack.Commands, c)
}

func (stack *CommandStack) Run(ctx context.Context) error {
	var err error
	for err == nil {
		c := stack.pop()
		if c == nil {
			return nil
		}
		err = c.Run(ctx, stack)
	}
	return err
}

func (stack *CommandStack) pop() Command {
	if len(stack.Commands) < 1 {
		return nil
	}
	top := len(stack.Commands) - 1
	c := stack.Commands[top]
	stack.Commands = stack.Commands[:top]
	return c
}

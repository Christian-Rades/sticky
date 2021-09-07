package runner

import (
	"context"
	"errors"
	"testing"
)


type commandWrapper struct {
	f func () error
}

func (c *commandWrapper) Run(_ context.Context, _ *CommandStack) error {
	return c.f()
}


func TestRunSingleCommand(t *testing.T) {
	ran := false
	stack := &CommandStack{[]Command{}}
	stack.Add(&commandWrapper{func() error {
		ran = true
		return nil
	}})
	stack.Run(context.Background())
	if !ran {
		t.Fail()
	}
}

func TestSelfQueueing(t *testing.T) {
	log := ""
	stack := &CommandStack{[]Command{}}
	stack.Add(&commandWrapper{func() error {
		log = log + "a"
		stack.Add(&commandWrapper{func() error {
			log = log + "b"
			return nil
		}})
		return nil
	}})
	stack.Run(context.Background())
	if log != "ab" {
		t.Fatalf("commands ran in wrong order %q", log)
	}
}

func TestStopOnError(t *testing.T) {
	log := ""
	stack := &CommandStack{[]Command{}}
	stack.Add(&commandWrapper{func() error {
		log = log + "a"
		stack.Add(&commandWrapper{func() error {
			log = log + "b"
			return nil
		}})
		return errors.New("test err")
	}})
	err := stack.Run(context.Background())
	if log != "a" {
		t.Errorf("commands ran in wrong order %q", log)
	}
	if err.Error() != "test err" {
		t.Errorf("wrong error %q", err)
	}
}

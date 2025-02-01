package types

import "context"

type ShutdownCallback struct {
	Name     string
	Callback func(ctx context.Context) error
	Priority int
}

func NewShutdownCallback(name string, callback func(context.Context) error, priority int) *ShutdownCallback {
	return &ShutdownCallback{
		Name:     name,
		Callback: callback,
		Priority: priority,
	}
}

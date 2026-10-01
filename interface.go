package buttonmonitor

import (
	"io"
)

type ButtonMonitor interface {
	io.Closer

	Subscribe() (Subscription[string], error)
}

type Subscription[T any] interface {
	io.Closer

	Ch() <-chan T
	ErrCh() <-chan error
}

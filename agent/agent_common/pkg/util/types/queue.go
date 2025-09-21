package types

import "io"

type Queue[T any] interface {
	io.Closer
	Deque[T]
	Pusher[T]
	Counter

	Max() int
}

type Counter interface {
	Count() int
}
type Pusher[T any] interface {
	Push(data T) error
}

type Deque[T any] interface {
	Pop() (T,error)
}
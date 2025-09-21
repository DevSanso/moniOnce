package collection


type StdQueue[T any] struct {
	ch chan T
	max int
}

func NewStdQueue[T any](size int) *StdQueue[T] {
	return &StdQueue[T]{
		ch : make(chan T, size),
		max : size,
	}
}

func (ns *StdQueue[T])Push(data T) error {
	ns.ch <- data
	return nil
}

func (ns *StdQueue[T])Pop() (T, error) {
	return <-ns.ch, nil
}

func (ns *StdQueue[T])Close() error {
	close(ns.ch)
	return nil
}

func (ns *StdQueue[T])Count() int {
	return len(ns.ch)
}

func (ns *StdQueue[T])Max() int {
	return ns.max
}
package logger

type DataLogger[T any] interface {
	Log(data *T) error
}
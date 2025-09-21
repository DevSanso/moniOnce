package types

type StrGetter interface {
	Get(key string) (string, error)
}

type StrSetter interface {
	Set(key string, val string) error
}

type StrGetterAndSetter interface {
	StrGetter
	StrSetter
}

type StrCloneToGetter interface {
	CloneFromGetter(StrGetter) 
}
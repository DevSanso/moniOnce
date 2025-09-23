package types

type StrGetter interface {
	Get(key string) (string, error)
}

type StrSetter interface {
	Set(key string, val string) error
}

type StrGetterKeyList interface {
	StrGetter
	Keys() []string
}

type StrGetterAndSetter interface {
	StrGetter
	StrSetter
}

type StrGetterKeyListAndSetter interface {
	StrGetterKeyList
	StrSetter
}

type StrCloneToGetter interface {
	CloneFromGetter(StrGetter) 
}

type SetterInter[T any] interface {
	*T
	StrSetter
}

type GetterKeysetterInter[T any] interface {
	*T
	StrGetterKeyListAndSetter
}
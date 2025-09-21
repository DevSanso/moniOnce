package config

import (
	"agent_common/pkg/util/types"
)

type Configure[T types.StrSetter] interface {
	Load() (T, error)
}

package dbpool

import (

)

type DbPoolData interface {
	Count() int
	GetData(idx int,  args...any) any
}

type DBPool interface {
	Exec(query string, args ...any) (DbPoolData, error)
}
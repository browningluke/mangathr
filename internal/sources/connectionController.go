package sources

import (
	"github.com/browningluke/mangathr/v2/internal/sources/mangadex"
)

type Connection interface {
	GetConnectionName() string
}

func NewConnection(name string) Connection {
	m := map[string]func() Connection{
		"mangadex": func() Connection { return mangadex.NewConnection() },
	}

	connection, ok := m[name]
	if !ok {
		panic("Passed connection name not in map")
	}
	return connection()
}

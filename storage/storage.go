package storage

import "github.com/RajNykDhulapkar/go-ww-sse/types"

type Storage interface {
	Get(int) *types.User
}

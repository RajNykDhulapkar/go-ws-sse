package storage

import "github.com/RajNykDhulapkar/go-ww-sse/types"

type MemoryStorage struct{}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) Get(id int) *types.User {
	return &types.User{
		ID:   id,
		Name: "Memory User",
	}
}

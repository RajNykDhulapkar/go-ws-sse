package storage

import "github.com/RajNykDhulapkar/go-ww-sse/types"

type MongoStorage struct{}

func (s *MongoStorage) Get(id int) *types.User {
	return &types.User{
		ID:   id,
		Name: "Mongo User",
	}
}

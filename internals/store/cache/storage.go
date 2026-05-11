package cache

import (
	"context"

	"github.com/Davidmuthee12/eazymarket/internals/store"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Users interface {
		Get(context.Context, string) (*store.User, error)
		Set(context.Context, *store.User) error
		Delete(context.Context, string) error
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rdb},
	}
}

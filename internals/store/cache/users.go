package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Davidmuthee12/eazymarket/internals/store"
	"github.com/redis/go-redis/v9"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute

func (s *UserStore) Get(ctx context.Context, userUUID string) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%s", userUUID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%s", user.UUID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, cacheKey, json, UserExpTime).Err()
}

func (s *UserStore) Delete(ctx context.Context, userUUID string) error {
	cacheKey := fmt.Sprintf("user-%s", userUUID)
	return s.rdb.Del(ctx, cacheKey).Err()
}

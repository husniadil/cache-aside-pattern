package repository

import (
	"fmt"
	"time"

	"github.com/bluele/gcache"
)

// RedisRepository simulates redis repository.
type RedisRepository struct {
	redis gcache.Cache
	mysql Repository
	ttl   time.Duration
}

// NewRedisRepository instantiates a new RedisRepository.
func NewRedisRepository(ttl time.Duration, mysql Repository) Repository {
	return &RedisRepository{
		redis: gcache.New(100).LRU().Build(),
		mysql: mysql,
		ttl:   ttl,
	}
}

// DoAnExpensiveQuery simulates mysql call
func (r *RedisRepository) DoAnExpensiveQuery(id string) (*string, error) {
	start := time.Now()
	defer func() {
		fmt.Printf("redis.DoAnExpensiveQuery took %s\n\n", time.Since(start))
	}()
	fmt.Printf("redis.DoAnExpensiveQuery: %s\n", id)

	// simulating network roundtrip for redis
	time.Sleep(5 * time.Millisecond)
	rawData, err := r.redis.Get(id)

	if err == nil && rawData != nil {
		// if the data is in redis, return it
		fmt.Printf("Cache hit for id: %s\n", id)
		data := rawData.(string)
		return &data, nil
	}

	if err != nil {
		// in case of error, do not return
		// have a try reading from database
		fmt.Printf("Cache miss for id: %s\n", id)
	}

	// if the data is not in the cache yet,
	// get it from database
	result, err := r.mysql.DoAnExpensiveQuery(id)
	if err != nil {
		return nil, err
	}
	// and eventually store the value to cache
	go func(result string) {
		r.redis.SetWithExpire(id, result, r.ttl)
	}(*result)

	return result, nil
}

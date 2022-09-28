package repository

import (
	"fmt"
	"time"

	"github.com/bluele/gcache"
	"golang.org/x/sync/singleflight"
)

// RedisRepository simulates redis repository.
type RedisRepository struct {
	redis gcache.Cache
	mysql Repository
	ttl   time.Duration
	sf    singleflight.Group
}

// NewRedisRepository instantiates a new RedisRepository.
func NewRedisRepository(ttl time.Duration, mysql Repository) Repository {
	return &RedisRepository{
		redis: gcache.New(100).LRU().Build(),
		mysql: mysql,
		ttl:   ttl,
		sf:    singleflight.Group{},
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
	mysqlResult, err, _ := r.sf.Do(id, func() (interface{}, error) {
		res, errMySQL := r.mysql.DoAnExpensiveQuery(id)
		if errMySQL != nil {
			return nil, errMySQL
		}
		return res, nil
	})
	if err != nil {
		return nil, err
	}

	var result *string
	switch res := mysqlResult.(type) {
	case string:
		result = &res
	case *string:
		result = res
	}

	// and eventually store the value to cache
	go func(result string) {
		r.redis.SetWithExpire(id, result, r.ttl)
	}(*result)

	return result, nil
}

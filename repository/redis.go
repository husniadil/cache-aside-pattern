package repository

import (
	"log"
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

// ResolvePersonNameByID simulates mysql call
func (r *RedisRepository) ResolvePersonNameByID(id string) (*string, error) {
	start := time.Now()
	defer func() {
		log.Printf("redis.ResolvePersonNameByID took %s\n\n", time.Since(start))
	}()
	log.Printf("redis.ResolvePersonNameByID: %s\n", id)
	// if the data is in redis, return it

	// simulating network roundtrip for redis
	time.Sleep(5 * time.Millisecond)
	rawData, err := r.redis.Get(id)

	if err == nil && rawData != nil {
		log.Printf("Cache hit for id: %s\n", id)
		data := rawData.(string)
		return &data, nil
	}

	if err != nil {
		// in case of error, do not return
		// have a try reading from database
		log.Printf("Cache miss for id: %s\n", id)
	}

	// if the data is not in the cache yet
	// get it from database
	// and eventually store the value to cache
	result, err := r.mysql.ResolvePersonNameByID(id)
	if err != nil {
		return nil, err
	}
	defer r.redis.SetWithExpire(id, *result, r.ttl)

	return result, nil
}

// Store simulates store db.
func (r RedisRepository) Store(id, name string) error {
	start := time.Now()
	defer func() {
		log.Printf("redis.Store took %s\n\n", time.Since(start))
	}()
	log.Printf("redis.Store: %s\n", id)

	// store it on mysql
	r.mysql.Store(id, name)

	// invalidate cache
	go func(id string) {
		// simulates latency
		time.Sleep(time.Millisecond * 5)
		r.redis.Remove(id)
	}(id)

	return nil
}

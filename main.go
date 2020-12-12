package main

import (
	"fmt"
	"time"

	"github.com/husniadil/cache-aside-pattern/repository"
)

func main() {
	mysqlRepository := repository.NewMySQLRepository()

	ttl := time.Second * 3
	redisRepository := repository.NewRedisRepository(ttl, mysqlRepository)

	id := "2c1b7cd2-0420-4b73-a3f9-734504842fb9"

	var names []string
	for i := 0; i < 5; i++ {
		name, err := redisRepository.ResolvePersonNameByID(id)
		if err != nil {
			fmt.Printf("Error loading [%s]: %s", id, err.Error())
			continue
		}
		names = append(names, *name)
	}

	// wait for cache expiration and we'll see a cache miss
	fmt.Printf("-------------- waiting cache expiration --------------\n\n")

	time.Sleep(ttl)

	for i := 0; i < 5; i++ {
		name, err := redisRepository.ResolvePersonNameByID(id)
		if err != nil {
			fmt.Printf("Error loading [%s]: %s", id, err.Error())
			continue
		}
		names = append(names, *name)
	}

	fmt.Println("Result:", names)
}

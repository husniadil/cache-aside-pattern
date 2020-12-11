package main

import (
	"fmt"
	"log"
	"time"

	"github.com/husniadil/cache-aside-pattern/repository"
)

func main() {
	mysqlRepository := repository.NewMySQLRepository()

	ttl := time.Second * 3
	redisRepository := repository.NewRedisRepository(ttl, mysqlRepository)

	var names []string
	for i := 0; i < 5; i++ {
		id := "2c1b7cd2-0420-4b73-a3f9-734504842fb9"
		name, err := redisRepository.ResolvePersonNameByID(id)
		if err != nil {
			log.Printf("Error loading [%s]: %s", id, err.Error())
			continue
		}
		names = append(names, *name)
	}

	// wait for cache expiration and we'll see a cache miss
	log.Printf("-------------- waiting cache expiration --------------\n\n")

	time.Sleep(ttl)

	for i := 0; i < 5; i++ {
		id := "2c1b7cd2-0420-4b73-a3f9-734504842fb9"
		name, err := redisRepository.ResolvePersonNameByID(id)
		if err != nil {
			log.Printf("Error loading [%s]: %s", id, err.Error())
			continue
		}
		names = append(names, *name)
	}

	fmt.Println("Result:", names)
}

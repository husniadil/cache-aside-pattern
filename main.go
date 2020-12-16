package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/husniadil/cache-aside-pattern/repository"
)

func main() {
	ttl := time.Second * 3

	repo := repository.NewMySQLRepository()
	repo = repository.NewRedisRepository(ttl, repo)

	id := "2c1b7cd2-0420-4b73-a3f9-734504842fb9"

	var wg sync.WaitGroup
	wg.Add(5)
	names := make([]string, 5)
	for i := 0; i < 5; i++ {
		go func(i int, id string) {
			defer wg.Done()
			name, err := repo.DoAnExpensiveQuery(id)
			if err != nil {
				fmt.Printf("Error loading [%s]: %s", id, err.Error())
				return
			}
			names[i] = *name
		}(i, id)
	}
	wg.Wait()

	fmt.Println("Result:", names)
}

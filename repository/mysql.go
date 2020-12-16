package repository

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// MySQLRepository simulates mysql repository.
type MySQLRepository struct {
	mysqldb sync.Map
}

// NewMySQLRepository instantiates a new MySQLRepository.
func NewMySQLRepository() Repository {
	repo := &MySQLRepository{
		mysqldb: sync.Map{},
	}

	// populate data for simulation
	repo.mysqldb.Store("2c1b7cd2-0420-4b73-a3f9-734504842fb9", "Husni")
	repo.mysqldb.Store("6e341b0b-dc78-4c59-91dc-d6251124e4b4", "Adil")
	repo.mysqldb.Store("ea5e9f28-46d8-4160-af68-6e0f71efd62d", "Makmur")
	return repo
}

// DoAnExpensiveQuery simulates mysql query.
func (m *MySQLRepository) DoAnExpensiveQuery(id string) (*string, error) {
	start := time.Now()
	defer func() {
		fmt.Printf("mysql.DoAnExpensiveQuery took %s\n", time.Since(start))
	}()
	fmt.Printf("mysql.DoAnExpensiveQuery: %s\n", id)

	// simulates latency
	time.Sleep(time.Millisecond * 100)

	// simulates get real data
	if rawData, ok := m.mysqldb.Load(id); ok {
		data := rawData.(string)
		return &data, nil
	}
	return nil, errors.New("record not found")
}

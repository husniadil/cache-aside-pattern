package repository

// Repository is a repository interface.
type Repository interface {
	DoAnExpensiveQuery(id string) (*string, error)
}

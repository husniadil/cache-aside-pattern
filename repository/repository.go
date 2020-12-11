package repository

// Repository is a repository interface.
type Repository interface {
	ResolvePersonNameByID(id string) (*string, error)
	Store(id, name string) error
}

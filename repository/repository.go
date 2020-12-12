package repository

// Repository is a repository interface.
type Repository interface {
	ResolvePersonNameByID(id string) (*string, error)
}

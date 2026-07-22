package storage

// Store is the persistence boundary for registrations, CDR and config.
// MVP keeps data in memory via other packages; implementations come later.
type Store interface {
	Ping() error
}

// MemoryStore is a no-op placeholder.
type MemoryStore struct{}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Ping() error {
	return nil
}

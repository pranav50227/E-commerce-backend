package cache

// Cache defines standard memory/distributed cache interfaces.
type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Delete(key string) error
}

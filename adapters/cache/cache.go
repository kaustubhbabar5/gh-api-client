package cache

import "time"

type Cache interface {
	Ping() error
	ReadString(key string) (string, error)
	WriteString(key string, value string) error
	Increment(key string) (int64, error)
	AddExpiry(key string, expiryTime time.Duration) error
	Delete(key string) error
	Close()
	WriteWithExpiry(key, value string, expiryTime time.Duration) error
}

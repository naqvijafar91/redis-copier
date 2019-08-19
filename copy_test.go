package rediscopy

import (
	"testing"

	"github.com/go-redis/redis"
)

func getLocalConnection() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func TestShouldCopy(t *testing.T) {
	copier := NewCopier(getLocalConnection(), getLocalConnection())
	err := copier.CopySortedSet("deltas", "deltas-copy")
	if err != nil {
		t.Error(err)
	}
}

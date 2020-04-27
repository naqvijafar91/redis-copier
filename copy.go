package rediscopy

import (
	"fmt"

	"github.com/go-redis/redis"
)

/*
Purpose of this project is to safely copy all keys from Redis running in production to your local one. This would be useful
when you cannot take a dump from production for some reason. Even if you are able to get that dump, this tool can help you
do it in the background in an automatic way.

Version 0.1 : Copy a sorted set, given a key in the command
*/

// Copier - Main struct which copies data
type Copier struct {
	server *redis.Client
	local  *redis.Client
}

/*
CopySortedSet - Copy a sorted set from production to local instance. Right now copy everything in one instance,
later it should be done in batches so that we do not put a lot of load on our Redis instance
*/
func (cp *Copier) CopySortedSet(sortedSetSource, sortedSetTarget string) error {
	vals, err := cp.server.ZRangeByScoreWithScores(sortedSetSource, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  0,
	}).Result()
	if err != nil {
		return err
	}
	valPtr := make([]*redis.Z, len(vals))
	for index, val := range vals {
		valPtr[index] = &val
		_, err := cp.local.ZAdd(sortedSetTarget, valPtr[0]).Result()
		if err != nil {
			return err
		}
	}

	return nil
}

/*
CopyKeyValue - Copy a key value from production to local instance. Right now copy everything in one instance,
later it should be done in batches so that we do not put a lot of load on our Redis instance
*/
func (cp *Copier) CopyKeyValue(sourceKey, targetKey string) error {
	val, err := cp.server.Get(sourceKey).Result()
	if err != nil {
		return err
	}
	fmt.Printf("Copying %s .....", sourceKey)
	err = cp.local.Set(targetKey, val, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func NewCopier(server, client *redis.Client) *Copier {
	return &Copier{server, client}
}

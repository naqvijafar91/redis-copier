package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/go-redis/redis"
	rediscopy "github.com/naqvijafar91/redis-copy"
)

func main() {
	sourceAddress := flag.String("s", "localhost:6379", "source address")
	sourcePassword := flag.String("sp", "", "source password")
	sourceDb := flag.Int("sdb", 0, " Source Database number")

	destinationAddress := flag.String("d", "localhost:6379", "destination address")
	destinationPassword := flag.String("dp", "", "source password")
	destinationDb := flag.Int("ddb", 0, " Source Database number")

	sourceSetName := flag.String("skey", "", "Source Key")
	destinationKeyName := flag.String("dkey", "", "Destination key")
	flag.Parse()
	// Now we need to make sure that we do not copy anything into production, we want a safety hook
	// Therefore, we would have localhost and 127.0.0.1 as the only 2 allowed destination addresses
	if !shouldAllowDestinationAddress(*destinationAddress) {
		panic(errors.New(fmt.Sprint("Destination cannot be ", *destinationAddress, ". Make sure you are not copying anything to prod")))
	}
	copier := rediscopy.NewCopier(getConnection(*sourceAddress, *sourcePassword, *sourceDb),
		getConnection(*destinationAddress, *destinationPassword, *destinationDb))
	// @Todo : Print and ask for confirmation before starting the work
	err := copier.CopySortedSet(*sourceSetName, *destinationKeyName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Key", *sourceSetName, "copied sucessfully")
}

func shouldAllowDestinationAddress(destinationAddress string) bool {
	address := strings.Split(destinationAddress, ":")
	if address[0] == "localhost" || address[0] == "127.0.0.1" {
		return true
	}
	return false
}

func getConnection(address, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:        address,
		Password:    password, // no password set
		DB:          db,       // use default DB
		ReadTimeout: 9999999999999999,
	})
}

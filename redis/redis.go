package redis

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/go-redis/redis"
)

// config contains all data related to establishing Redis connection.
type config struct {
	Address  string `json:"redis_address"`
	Password string `json:"redis_password"`
	Database int    `json:"redis_db"`
}

var (
	conf config
	// Client represents Redis connection to work with.
	Client *redis.Client
)

func init() {
	getConnection()
}

// getConnection establishes Redis connection.
func getConnection() {
	var err error

	file, err := filepath.Abs("redis/config.json")
	if err != nil {
		log.Fatal(err)
	}

	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(fileContent, &conf)
	if err != nil {
		log.Fatal(err)
	}

	Client = redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		DB:       conf.Database,
	})

	_, err = Client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
}

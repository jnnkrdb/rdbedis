package rdbedis

import (
	"encoding/base64"
	"time"

	// go get github.com/go-redis/redis
	"github.com/go-redis/redis"

	"github.com/jnnkrdb/corerdb/prtcl"
)

// Struct of the redis-configuration as json.
// Ensure the "passwd" value is base64 encrypted.
//
//	{
//	  "address" : "redis:port",
//	  "passwd" : "<password>",
//	  "dbindex" : 0
//	}
type Redis struct {

	// collection of vars which will be collected from the json

	URI  string `json:"uri"`
	Port string `json:"port"`
	// paswword must be base64
	Password string `json:"password"`
	DBIndex  int    `json:"dbindex"`

	// undefined vars
	client *redis.Client
}

// decode password from base64
func (rds Redis) getUnencodedPW() string {

	if str, err := base64.StdEncoding.DecodeString(rds.Password); err != nil {

		return ""

	} else {

		return string(str)
	}
}

// connect to the redis from the config
func (rds *Redis) Connect() {

	prtcl.Log.Println("connecting to redis")

	rds.client = redis.NewClient(&redis.Options{
		Addr:     rds.URI + ":" + rds.Port,
		Password: rds.getUnencodedPW(),
		DB:       rds.DBIndex,
	})

	rds.CheckConnection()
}

// check the connection to redis
func (rds *Redis) CheckConnection() error {

	prtcl.Log.Println("checking connection")

	if pong, err := rds.client.Ping().Result(); err != nil {

		prtcl.PrintObject(rds, pong, err)

		return err

	} else {

		prtcl.Log.Println("connection established")

		return nil
	}
}

// -----------------------------------------------------------------
// functions for the actual redis part, like adding, reading and deleting keys

// add a key value (plus duration) pair to the redis instance
//
// Parameters:
//   - `key` : string > used as the address of the value
//   - `value` : string > contains the value, of the key
//   - `expirationtime` : int > expiration time of the key-value-pair, 0 means it does not expire
func (rds Redis) Add(key, value string, expirationtime int) error {

	prtcl.Log.Println("adding key:value [duration] to redis", rds.URI, rds.DBIndex)

	if err := rds.CheckConnection(); err != nil {

		prtcl.Log.Println("not connected to", rds.URI)

		return err

	} else {

		if err := rds.client.Set(key, value, time.Duration(expirationtime)); err != nil {

			prtcl.Log.Println("error while adding key:value [duration]")

			prtcl.PrintObject(rds, err)
		}

		return nil
	}
}

// read a specific key from the redis instance
//
// Parameters:
//   - `key` : string > used to address the value
func (rds Redis) Read(key string) (string, error) {

	prtcl.Log.Println("reading key from redis", rds.URI, rds.DBIndex)

	if err := rds.CheckConnection(); err != nil {

		prtcl.Log.Println("not connected to", rds.URI)

		return "", err

	} else {

		if result, err := rds.client.Get(key).Result(); err != nil {

			prtcl.Log.Println("error while reading key", key)

			prtcl.PrintObject(rds, result, err)

			return "", err

		} else {

			return result, nil
		}
	}
}

// delete a specific key from the redis instance
// DOES NOT WORK NOW
//
// Parameters:
//   - `key` : string > used to address the value
func (rds Redis) Delete(key string) error {

	return nil
}

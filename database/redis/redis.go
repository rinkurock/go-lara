package redis

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	c "app/config"
	h "app/helpers"
)

var redisClient *redis.Client
var once sync.Once

func Initialize() {
	once.Do(func() {
		_conf := c.GetConfig().Redis

		log.Infoln("connecting to redis at :", _conf.Host, _conf.Port, " ....")

		//TODO POOLING AND RELATED
		redisClient = redis.NewClient(&redis.Options{
			Addr:        _conf.Host + ":" + h.ToString(_conf.Port),
			Password:    _conf.Password,
			DB:          _conf.Db,
			DialTimeout: time.Duration(_conf.DialTimeout) * time.Second,
		})

		if _, err := redisClient.Ping().Result(); err != nil {
			log.Println("failed to connect redis: ", err)
			panic(err)
		}
		log.Infoln("redis connection successful !")
	})
}

func Set(key string, value interface{}, ttl time.Duration) error {
	byteData, err := json.Marshal(value)
	if err != nil {
		log.Error(err)
		return err
	}
	err = redisClient.Set(key, string(byteData), ttl).Err()
	//err := redisClient.Set(key, value, ttl).Err()
	if err != nil {
		log.Error(err)
	}
	return err
}

func Get(key string, out interface{}) error {
	value, err := redisClient.Get(key).Result()
	if err != nil {
		log.Error(err)
		return err
	}
	err = json.Unmarshal([]byte(value), &out)
	return err

}

func Del(key string) error {
	_, error := redisClient.Del(key).Result()
	return error
}
func GetByteArrayToMap(key string, out interface{}) error {
	value, err := redisClient.Get(key).Result()
	if err != nil {
		log.Error(err)
		return err
	}
	m, err := h.ByteArrayToMap([]byte(value))
	if err != nil {
		log.Error(err)
		return err
	}
	err = h.MapToStructWeak(m, &out, true)

	return err
}

func Remove(key string) error {
	_, error := redisClient.Del(key).Result()
	return error
}

func MultiRemove(multiKey []string) error {
	keys := strings.Join(multiKey, " ")
	_, err := redisClient.Del(keys).Result()
	if err != nil {
		log.Errorln("Multi Remove ERROR: " + err.Error())
	}
	return err
}
func MultipleGet(keys []string, out ...interface{}) error {
	//check keys and out length
	if len(keys) != len(out) {
		return errors.New("number of keys and number of interface should match")
	}

	//Get data from redis
	redisMGetResult := redisClient.MGet(keys...)
	interfaces, err := redisMGetResult.Result()
	if err != nil {
		return err
	}

	//Loop over interfaces found in redis MGet
	for i, inFs := range interfaces {
		if inFs == nil {
			continue
		}

		//try to decode single interface to 'out'
		m, err := h.ByteArrayToMap([]byte(inFs.(string)))
		if err != nil {
			return err
		}
		err = h.MapToStructWeak(m, &out[i], true)
		if err != nil {
			return err
		}
	}
	return nil
}

func PingRedis() error {
	if redisClient != nil {
		_, err := redisClient.Ping().Result()
		return err
	}
	return errors.New("redis client is not initialized")
}

func GetClient() *redis.Client {
	if redisClient == nil {
		Initialize()
	}
	return redisClient
}
func CloseRedis() {
	if redisClient != nil {
		_ = redisClient.Close()
	}
}
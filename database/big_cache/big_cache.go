package bigcache

import (
	c "app/config"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/allegro/bigcache"
	log "github.com/sirupsen/logrus"
)

var once sync.Once
var _cache *bigcache.BigCache

func Initialize() {
	once.Do(func() {
		ttl := time.Duration(c.GetConfig().Cache.InMemory.TTLInSec) * time.Second
		bConf := bigcache.DefaultConfig(ttl)
		bConf.CleanWindow = ttl
		c, err := bigcache.NewBigCache(bConf)
		if err != nil {
			log.Errorln(err)
			panic(err)
		}
		_cache = c
	})
}

func Set(key string, value interface{}) (err error) {
	// convert interface to byte
	var b []byte
	b, err = json.Marshal(value)
	if err == nil {
		err = _cache.Set(key, b)
	}

	if err != nil {
		log.Errorln(err)
	}

	return
}

func Get(key string, out interface{}) (err error) {
	var b []byte

	b, err = _cache.Get(key)
	if err == nil {
		if len(b) > 0 {
			err = json.Unmarshal(b, &out)
		} else {
			errorText := "no data found in cache"
			err = errors.New(errorText)
		}
	}

	if err != nil {
		log.Errorln(err)
	}

	return
}

func Delete(key string) error {
	return _cache.Delete(key)
}

func Close() {
	_ = _cache.Close()
}

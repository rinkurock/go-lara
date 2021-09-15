package config

import (
	"app/helpers"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var once sync.Once
var conf = &_Config{}
var runtimeViper = viper.New()

type _Config struct {
	IsLoadedFromConsul bool
	IsLoadedFromFile   bool
	IsLoaded           bool
	Server             struct {
		Port         int `mapstructure:"port" json:"port"`
		ReadTimeout  int `mapstructure:"read_timeout" json:"read_time_out"`
		WriteTimeout int `mapstructure:"write_timeout" json:"write_time_out"`
	} `mapstructure:"server" json:"server"`
	Database struct {
		Host               string `mapstructure:"host" json:"host"`
		Port               int    `mapstructure:"port" json:"port"`
		UserName           string `mapstructure:"username" json:"username"`
		Password           string `mapstructure:"password" json:"password"`
		DatabaseName       string `mapstructure:"database_name" json:"database_name"`
		ConnMaxLifetimeMin int    `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime"`
		MaxIdleConnections int    `mapstructure:"max_idle_connections" json:"max_idle_connections"`
		MaxOpenConnections int    `mapstructure:"max_open_connections" json:"max_open_connections"`
		Debug              bool   `mapstructure:"debug" json:"debug"`
	} `mapstructure:"database" json:"database"`
	Redis struct {
		Host        string `mapstructure:"host" json:"host"`
		Port        int    `mapstructure:"port" json:"port"`
		Password    string `mapstructure:"password" json:"password"`
		Db          int    `mapstructure:"db" json:"db"`
		Prefix      string `mapstructure:"prefix" json:"prefix"`
		DialTimeout int    `mapstructure:"dial_timeout" json:"dial_timeout"`
		TTLInSec    int    `mapstructure:"ttl_in_sec" json:"ttl_in_sec"`
	} `mapstructure:"redis" json:"redis"`
	Cache struct {
		Redis struct {
			PutEnabled bool `mapstructure:"put_enabled" json:"put_enabled"`
			GetEnabled bool `mapstructure:"get_enabled" json:"get_enabled"`
		} `json:"redis"`
		InMemory struct {
			PutEnabled bool `mapstructure:"put_enabled" json:"put_enabled"`
			GetEnabled bool `mapstructure:"get_enabled" json:"get_enabled"`
			TTLInSec   int  `mapstructure:"ttl_in_sec" json:"ttl_in_sec"`
		} `mapstructure:"in_memory" json:"in_memory"`
	}
	Others struct {
		ConsulWatch        bool   `mapstructure:"consul_watch" json:"consul_watch"`
		ConsulWatchTimeout int    `mapstructure:"consul_watch_timeout" json:"consul_watch_timeout"`
		LogFormat          string `mapstructure:"log_format" json:"log_format"`
		LogLevel           string `mapstructure:"log_level" json:"log_level"`
		GoMaxProcess       int    `mapstructure:"go_max_process" json:"go_max_process"`
		ResponseLog        bool   `mapstructure:"response_log" json:"response_log"`
	} `mapstructure:"others" json:"others"`
}

func GetConfig() *_Config {
	once.Do(func() {
		//Once.Do will ensure that codes inside it will run only once in the whole application life
		setDefaultValues()
		if cast.ToBool(helpers.GetEnv("CONSUL", "false")) {
			loadConfigFromConsul()
			conf.IsLoadedFromConsul = true
		} else {
			log.Info("consul config is disabled !!")
			loadConfigJson()
			log.Info("config load from json file")
		}
		conf.IsLoaded = true
		if conf.Others.ConsulWatch {
			watchConfigChangForConsul()
		}
	})

	return conf
}

func loadConfigJson() {
	jsonFileConf := helpers.GetEnv("CONFIG_JSON", "config.json")
	configJson, err := os.Open(jsonFileConf)
	if err != nil {
		log.Error("ERROR ON LOAD JSON CONFIG: " + jsonFileConf)
		log.Panic(err)
	}

	byteValue, err := ioutil.ReadAll(configJson)
	if err != nil {
		log.Error("ERROR ON LOAD JSON CONFIG: " + jsonFileConf)
		log.Panic(err)
	}
	//json.Unmarshal(byteValue, &conf)
	if err = json.Unmarshal(byteValue, &conf); err != nil {

		log.Error("ERROR ON LOAD JSON CONFIG Unmarshal: " + jsonFileConf)
		log.Panic(err)
	}
	conf.IsLoadedFromFile = true
	return

}

func loadConfigFromConsul() {
	consulHost := helpers.GetEnv("CONSUL_URL", "127.0.0.1:8500")
	consulPath := helpers.GetEnv("CONSUL_PATH", "bid")

	if err := runtimeViper.AddRemoteProvider("consul", consulHost, consulPath); err == nil {
		runtimeViper.SetConfigType("json")
		if err = runtimeViper.ReadRemoteConfig(); err == nil {
			if err = runtimeViper.Unmarshal(&conf); err != nil {
				log.Error("ERROR ON  REMOTE CONFIG Unmarshal")
				log.Error(err)
				panic(err)
			}
			conf.IsLoadedFromConsul = true
			log.Infoln("got config from consul")
			b, _ := json.Marshal(conf)
			log.Infoln(string(b))
		} else {
			log.Error("ERROR ON LOAD REMOTE CONFIG FROM: " + consulHost + " AND KEY: " + consulPath)
			log.Error(err)
			panic(err)
		}
	} else {
		log.Error(err)
		panic(err)
	}
}

func watchConfigChangForConsul() {
	// open a goroutine to watch remote changes forever
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(conf.Others.ConsulWatchTimeout)) // delay after each request
			// currently, only tested with etcd support
			err := runtimeViper.WatchRemoteConfig()
			if err != nil {
				log.Errorf("unable to read remote config: %v", err)
				continue
			}
			newConf := &_Config{}
			if err = runtimeViper.Unmarshal(&newConf); err == nil {
				conf = newConf
				log.Debugln(fmt.Sprintf("%#v", newConf))
			} else {
				log.Error("error watching config from consul")
				log.Error(err)
			}
			timeOut := runtimeViper.GetInt("others.consul_watch_timeout")
			if timeOut > 0 {
				conf.Others.ConsulWatchTimeout = timeOut
			}
		}
	}()
}

func setDefaultValues() {
	conf.IsLoadedFromConsul = false
	conf.IsLoadedFromFile = false
	conf.IsLoaded = false
	conf.Server.Port = 8080
	conf.Server.ReadTimeout = 4
	conf.Server.WriteTimeout = 4

	conf.Database.Host = "127.0.0.1"
	conf.Database.Port = 3306
	conf.Database.DatabaseName = "bid"
	conf.Database.UserName = "root"
	conf.Database.Password = "c0mm0m"
	conf.Database.ConnMaxLifetimeMin = 15
	conf.Database.MaxIdleConnections = 5
	conf.Database.MaxOpenConnections = 10
	conf.Database.Debug = false

	conf.Redis.Host = "127.0.0.1"
	conf.Redis.Port = 6379
	conf.Redis.TTLInSec = 120

	conf.Cache.Redis.GetEnabled = true
	conf.Cache.Redis.PutEnabled = true

	conf.Cache.InMemory.GetEnabled = false
	conf.Cache.InMemory.PutEnabled = false
	conf.Cache.InMemory.TTLInSec = 10

	conf.Others.ConsulWatchTimeout = 5

	conf.Others.LogFormat = "json"

	conf.Others.GoMaxProcess = -1
	conf.Others.ResponseLog = false
}

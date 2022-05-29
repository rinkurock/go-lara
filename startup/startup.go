package startup

import (
	"app/config"
	"app/conn"
	bc "app/database/big_cache"
	m "app/database/mysql"
	r "app/database/redis"
	h "app/helpers"
	"app/middlewares"
	"app/routes"
	"app/services/logger"
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

func Initialize() {
	conf := config.GetConfig()
	println("Running Echo Server")
	log.Infoln("go max process :", runtime.GOMAXPROCS(conf.Others.GoMaxProcess))

	e := echo.New()
	routes.DefineRoutes(e)
	middlewares.ApplyMiddleware(e)

	//initialize logger with config
	logger.Initialize()

	//initialize mysql database
	m.InitDatabase()
	defer m.Close()

	//Redis
	r.Initialize()
	defer r.CloseRedis()

	//Configure in-memory cache
	if conf.Cache.InMemory.GetEnabled || conf.Cache.InMemory.PutEnabled {
		bc.Initialize()
		defer bc.Close()
	}
	// Http Service Connection
	conn.ServiceConnection()

	s := &http.Server{
		Addr:         ":" + h.ToString(conf.Server.Port),
		ReadTimeout:  time.Duration(conf.Server.ReadTimeout) * time.Minute,
		WriteTimeout: time.Duration(conf.Server.WriteTimeout) * time.Minute,
	}

	go func() {
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	} else {
		e.Logger.Info("gracefully shutting down server")
	}

	//e.Logger.Fatal(e.StartServer(s))
	//e.Logger.Fatal(e.StartAutoTLS(":443"))
}

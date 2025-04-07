package main

import (
	"net/http"
	"os"
	"time"

	"todo-app/internal/config"
	"todo-app/internal/db"
	_taskHTTPHndlr "todo-app/internal/delivery/http"
	_repo "todo-app/internal/repository"
	_taskUscase "todo-app/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// initialize logger configurations
func initLogger() {
	logLevel := logrus.ErrorLevel
	switch config.Env() {
	case "dev", "development":
		logLevel = logrus.InfoLevel
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		DisableSorting:  true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05 01-01-2025",
	})

	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	logrus.SetLevel(logLevel)
}

// run initLogger() before running main()
func init() {
	config.GetConf()
	initLogger()
}

func main() {
	e := echo.New()

	db.InitializePostgresConn()
	db.InitializeRedisConn()

	cacheRepo := _repo.NewCacheRepository(db.RedisClient)
	taskRepo := _repo.NewTaskRepository(db.PostgresDB, cacheRepo)
	taskUsecase := _taskUscase.NewTaskUsecase(taskRepo)
	_taskHTTPHndlr.NewTaskHTTPHandler(e, taskUsecase)

	s := &http.Server{
		Addr:         ":" + config.ServerPort(),
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 2 * time.Minute,
	}

	logrus.Fatal(e.StartServer(s))
}

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bluecover/qiniu_token/model"
	"github.com/bluecover/qiniu_token/object"
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

func initConfig(configPath string) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName("stash")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("stash")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func initDB() *gorm.DB {
	mysqlDSN := viper.GetString("mysql.dsn")
	db, err := gorm.Open("mysql", mysqlDSN)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(model.AllModels()...).Error
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	var logger log.Logger
	logger = log.NewJSONLogger(os.Stderr)
	logger = log.With(logger, "ts", log.TimestampFormat(time.Now, "2006-01-02 15:04:05.000000"))
	logger = log.With(logger, "caller", log.DefaultCaller)

	configPath := os.Getenv("STASH_CONFIG_PATH")
	if len(configPath) == 0 {
		configPath = "config"
		logger.Log("warning", "no STASH_CONFIG_PATH in env, use default")
	}

	logger.Log("configPath", configPath)

	initConfig(configPath)

	// Do not use database.
	// db := initDB()
	// defer db.Close()
	// fmt.Println("done: make database connection")
	var db *gorm.DB

	// Create service.
	service, err := object.NewService(db, logger, configPath)
	if err != nil {
		panic(err)
	}
	handler := object.MakeHTTPHandler(service, logger)

	fmt.Println("done: create service")

	// Create URL routing.
	mux := http.NewServeMux()
	mux.Handle("/v1/oss/", handler)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		httpAddr := fmt.Sprintf("%s:%d", viper.GetString("server.addr"), viper.GetInt("server.port"))
		logger.Log("transport", "HTTP", "addr", httpAddr)
		errs <- http.ListenAndServe(httpAddr, mux)
	}()

	logger.Log("exit", <-errs)
}

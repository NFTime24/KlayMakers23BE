package main

import (
	"flag"
	"github.com/nftime/api"
	"github.com/nftime/config"
	"github.com/nftime/db"
	"github.com/nftime/logger"
	"github.com/nftime/redis"
	"github.com/nftime/route"
	"log"
	"os"
)

// @title NFTime Sample Swagger API
// @version 1.0
// @host localhost:9200
// @BasePath /
func main() {
	environmentPtr := flag.String("env", "unknown", "env mode name")
	flag.Parse()

	logger.LogInit(os.Stdout, os.Stdout, os.Stderr)

	redis.Init()

	streamCfg, s3Info, err := config.GetStreamConfig(*environmentPtr)

	timeCfg, s3Info, err := config.GetTimeConfig(*environmentPtr)

	if err != nil {
		logger.Error.Printf("failed to get env setting", err)
		log.Fatal()
	} else {
		logger.Debug.Printf("streamCfg: %v\n", streamCfg)
		logger.Debug.Printf("timeCfg: %v\n", timeCfg)
		logger.Debug.Printf("s3Info: %v\n", s3Info)
	}

	//db.InitStreamPlatformDb()
	db.InitTimeStorageDb()

	e := route.Init()

	//logger.Info.Printf("Home: %v\n", os.Getenv("APP_ENV"))
	api.KlipRequestMap = make(map[uint64]string)

	logger.Info.Printf("Starting %v server...", "NFTIME")

	e.Run(":9200")
	//e.RunTLS(":443", "server.crt", "server.key")
}

package db

import (
	"fmt"
	"github.com/nftime/config"
	"github.com/nftime/model"
	"gorm.io/driver/postgres"
	"time"

	//"github.com/nftime/logger"
	"gorm.io/gorm/logger"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

var err error
var streamPlatformDb *gorm.DB
var timeStorageDb *gorm.DB

func InitStreamPlatformDb() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=stream port=%s sslmode=%s", config.Cfg.DB.PostgresHost, config.Cfg.DB.PostgresUser, config.Cfg.DB.PostgresPassword, config.Cfg.DB.PostgresPort, config.Cfg.DB.PostgresSSLMode)
	streamPlatformDb, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NowFunc: func() time.Time {
				utc, _ := time.LoadLocation("Asia/Seoul")
				return time.Now().In(utc)
			},
		})
	if err != nil {
		panic("DB Connection Error")
	}

	err := streamPlatformDb.AutoMigrate(&model.Artist{}, &model.User{}, &model.Work{}, &model.Nft{}, &model.Like{}, &model.Fantalk{}, &model.Log{}, &model.Ticket{}, &model.Playlist{}, &model.PlaylistWork{}, &model.SocialUser{})
	if err != nil {
		return
	}
	//db.AutoMigrate(&model.Log{})
}

func InitTimeStorageDb() {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=time port=%s sslmode=%s", config.Cfg.DB.PostgresHost, config.Cfg.DB.PostgresUser, config.Cfg.DB.PostgresPassword, config.Cfg.DB.PostgresPort, config.Cfg.DB.PostgresSSLMode)
	timeStorageDb, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NowFunc: func() time.Time {
				utc, _ := time.LoadLocation("Asia/Seoul")
				return time.Now().In(utc)
			},
		})
	if err != nil {
		panic("DB Connection Error")
	}

	timeStorageDb.AutoMigrate(&model.CertiUser{}, &model.Company{}, &model.Certificate{}, &model.CertificateUser{})
	//db.AutoMigrate(&model.Log{})
}

func StreamPlatformDbManager() *gorm.DB {
	if streamPlatformDb != nil {
		return streamPlatformDb
	} else {
		panic(err)
	}
}

func TimeStorageDbManager() *gorm.DB {
	if timeStorageDb != nil {
		return timeStorageDb
	} else {
		panic(err)
	}
}

package config

import (
	"17live_wso_be/util"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port     int
	Mode     string
	Database struct {
		Username string
		Password string
		Host     string
		Port     string
		Name     string
	}
	Jwt struct {
		Hmac        string
		Issure      string
		ExpiredHour time.Duration
	}
	User struct {
		Admin   string
		Domains []string
	}
	Google struct {
		Endpoint   string
		TokenQuery string
	}
	Email struct {
		Sender     string
		SenderName string
		ApiKey     string
		Content    struct {
			NoRegion       EmailContent
			SyncDataFinish EmailContent
		}
	}
	Media struct {
		Authentication   MediaConfig
		Campaign         MediaConfig
		Leaderboard      MediaConfig
		Streamer         MediaConfig
		StreamerContract MediaConfig
	}
}

type EmailContent struct {
	Subject   string
	PlainText string
	HtmlText  string
}

type MediaConfig struct {
	Server         string
	Path           string
	ClientID       string
	Secret         string
	TokenHeaderKey string
	Token          string
	StatusQuery    string
	StatusValue    string
	LimitQuery     string
	LimitValue     string
	RegionQuery    string
	RegionValue    string
	CursorQuery    string
	CursorValue    string
	CountQuery     string
	CountValue     string
}

var (
	once            sync.Once
	instance        *Config
	log             = util.GetLogger()
	configDirectory = "config"
)

func Load(filename string, config interface{}) {
	viper.AddConfigPath(configDirectory)
	viper.SetConfigName(filename)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("fail to load config file: %s", err.Error())
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("fail to decode config: %s", err.Error())
	}
}

func New() *Config {
	once.Do(func() {
		viper.SetConfigType("yaml")
		instance = &Config{}
		Load("config", &instance)
		log.Info("config initialized")
	})
	return instance
}

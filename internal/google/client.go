package google

import (
	"17live_wso_be/config"
	"17live_wso_be/util"
	"sync"
)

type Client struct {
	Endpoint   string
	TokenQuery string
}

var (
	once   sync.Once
	log    = util.GetLogger()
	client *Client
)

func New() *Client {
	cfg := config.New().Google

	once.Do(func() {
		client = &Client{
			Endpoint:   cfg.Endpoint,
			TokenQuery: cfg.TokenQuery,
		}

		log.Info("google client initialized")
	})

	return client
}

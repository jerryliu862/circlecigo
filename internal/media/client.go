package media

import (
	"17live_wso_be/config"
	"17live_wso_be/util"
	"sync"
)

type Client struct {
	Authentication   config.MediaConfig
	Campaign         config.MediaConfig
	Leaderboard      config.MediaConfig
	Streamer         config.MediaConfig
	StreamerContract config.MediaConfig
}

var (
	once   sync.Once
	log    = util.GetLogger()
	client *Client
)

func New() *Client {
	cfg := config.New().Media

	once.Do(func() {
		client = &Client{
			Authentication:   cfg.Authentication,
			Campaign:         cfg.Campaign,
			Leaderboard:      cfg.Leaderboard,
			Streamer:         cfg.Streamer,
			StreamerContract: cfg.StreamerContract,
		}

		log.Info("mediaApi client initialized")
	})

	return client
}

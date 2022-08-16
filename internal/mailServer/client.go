package mailServer

import (
	"17live_wso_be/config"
	"17live_wso_be/util"
	"sync"
)

type Client struct {
	Sender     string
	SenderName string
	ApiKey     string
	Content    struct {
		NoRegion       config.EmailContent
		SyncDataFinish config.EmailContent
	}
}

var (
	once   sync.Once
	log    = util.GetLogger()
	client *Client
)

func New() *Client {
	cfg := config.New().Email

	once.Do(func() {
		client = &Client{
			Sender:     cfg.Sender,
			SenderName: cfg.SenderName,
			ApiKey:     cfg.ApiKey,
			Content:    cfg.Content,
		}

		log.Info("mail server initialized")
	})

	return client
}

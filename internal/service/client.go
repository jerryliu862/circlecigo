package service

import (
	"17live_wso_be/internal/google"
	"17live_wso_be/internal/mailServer"
	"17live_wso_be/internal/media"
	"17live_wso_be/internal/repository"
	"17live_wso_be/util"
	"sync"
)

type Client struct {
	RepositoryClient *repository.Client
	GoogleClient     *google.Client
	MediaClient      *media.Client
	MailClient       *mailServer.Client
}

var (
	once   sync.Once
	log    = util.GetLogger()
	client *Client
)

func New() *Client {
	once.Do(func() {
		client = &Client{
			RepositoryClient: repository.New(),
			GoogleClient:     google.New(),
			MediaClient:      media.New(),
			MailClient:       mailServer.New(),
		}

		log.Info("service client initialized")
	})
	return client
}

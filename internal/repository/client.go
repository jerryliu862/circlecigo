package repository

import (
	"17live_wso_be/config"
	"17live_wso_be/util"

	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Client struct {
	Database *gorm.DB
}

const (
	Currency            string = "Currency"
	Region              string = "Region"
	TaxRate             string = "TaxRate"
	User                string = "User"
	UserAuth            string = "UserAuth"
	Bank                string = "Bank"
	BankBranch          string = "BankBranch"
	BankAccount         string = "BankAccount"
	StreamerAgency      string = "StreamerAgency"
	Streamer            string = "Streamer"
	Campaign            string = "Campaign"
	CampaignLeaderboard string = "CampaignLeaderboard"
	CampaignRank        string = "CampaignRank"
	CampaignBonus       string = "CampaignBonus"
	Payout              string = "Payout"
	PayoutDetail        string = "PayoutDetail"
	PayoutRecord        string = "PayoutRecord"

	ViewCampaignList                string = "v_campaignList"
	ViewRegionList                  string = "v_regionList"
	ViewRegionTaxList               string = "v_regionTaxList"
	ViewUngroupedPayout             string = "v_ungroupedPayout"
	ViewUnpaidBonusList             string = "v_unpaidBonusList"
	ViewPayoutList                  string = "v_payoutList"
	ViewGroupedPayoutDetail         string = "v_payoutDetail"
	ViewUngroupedPayoutDetail       string = "v_payoutDetail2"
	ViewPayoutRecord                string = "v_payoutRecord"
	ViewPayoutReportList            string = "v_payoutReportList"
	ViewPayoutReportDetail          string = "v_payoutReportDetail"
	ViewPayoutReportDetailByAgency  string = "v_payoutReportDetailByAgency"
	ViewPayoutReportDetailByAgency2 string = "v_payoutReportDetailByAgency2"
)

var (
	once   sync.Once
	client *Client
	log    = util.GetLogger()
)

func New() *Client {
	once.Do(func() {
		cfg := config.New().Database

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Panicf("fail to connect database: %s", err.Error())
		}

		client = &Client{
			Database: db,
		}

		log.Info("repository client initialized")
	})

	return client
}

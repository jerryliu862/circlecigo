package model

import (
	"github.com/google/uuid"
)

func (d *CampaignSet) SplitData() *CampaignSet {
	for len(d.Campaigns) > 1000 {
		d.CampaignsDataSplit = append(d.CampaignsDataSplit, d.Campaigns[:1000])
		d.Campaigns = d.Campaigns[1000:]
	}
	d.CampaignsDataSplit = append(d.CampaignsDataSplit, d.Campaigns)

	for len(d.Leaderboards) > 1000 {
		d.LeaderboardsDataSplit = append(d.LeaderboardsDataSplit, d.Leaderboards[:1000])
		d.Leaderboards = d.Leaderboards[1000:]
	}
	d.LeaderboardsDataSplit = append(d.LeaderboardsDataSplit, d.Leaderboards)

	for len(d.Ranks) > 1000 {
		d.RanksDataSplit = append(d.RanksDataSplit, d.Ranks[:1000])
		d.Ranks = d.Ranks[1000:]
	}
	d.RanksDataSplit = append(d.RanksDataSplit, d.Ranks)

	for len(d.Streamers) > 1000 {
		d.StreamersDataSplit = append(d.StreamersDataSplit, d.Streamers[:1000])
		d.Streamers = d.Streamers[1000:]
	}
	d.StreamersDataSplit = append(d.StreamersDataSplit, d.Streamers)

	for len(d.StreamerAgencies) > 1000 {
		d.StreamerAgenciesDataSplit = append(d.StreamerAgenciesDataSplit, d.StreamerAgencies[:1000])
		d.StreamerAgencies = d.StreamerAgencies[1000:]
	}
	d.StreamerAgenciesDataSplit = append(d.StreamerAgenciesDataSplit, d.StreamerAgencies)

	return d
}

type CampaignSet struct {
	Campaigns                 []CampaignInsert
	CampaignsDataSplit        [][]CampaignInsert
	Leaderboards              []CampaignLeaderboard
	LeaderboardsDataSplit     [][]CampaignLeaderboard
	Ranks                     []CampaignRank
	RanksDataSplit            [][]CampaignRank
	Streamers                 []Streamer
	StreamersDataSplit        [][]Streamer
	StreamerAgencies          []StreamerAgency
	StreamerAgenciesDataSplit [][]StreamerAgency
	Regions                   []Region
}

type MediaApiAuthentication struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type MediaApiCampaignSet struct {
	Campaign     MediaApiCampaign      `json:"campaign"`
	Leaderboards []MediaApiLeaderboard `json:"leaderboards"`
}

type MediaApiCampaign struct {
	Id          int                   `json:"ID"`
	Regions     []string              `json:"regions"`
	DefaultLang string                `json:"defaultLang"`
	Names       MediaApiMultiLanguage `json:"names"`
	OpenName    string                `json:"openName"`
	StartTime   int64                 `json:"startTime"`
	EndTime     int64                 `json:"endTime"`
}

type MediaApiLeaderboard struct {
	Id        uuid.UUID `json:"ID"`
	GroupName string    `json:"groupName"`
}

type MediaApiLeaderboardDetail struct {
	Ranks []MediaApiRank `json:"data"`
}

type MediaApiRank struct {
	Streamer uuid.UUID `json:"key"`
	Score    int       `json:"value"`
	Rank     int       `json:"rank"`
}

type MediaApiStreamer struct {
	Name   string `json:"name"`
	OpenID string `json:"openID"`
}

type MediaApiStreamerContract struct {
	Region         string `json:"region"`
	AgencyID       *int   `json:"agencyID"`
	AgencyName     string `json:"agencyName"`
	PayoutAccounts struct {
		Third MediaApiStreamerAccount `json:"3"`
	} `json:"payoutAccounts"`
}

type MediaApiStreamerAccount struct {
	SwiftCode    string  `json:"swiftCode"`
	OwnerType    *int    `json:"ownerType"`
	AccountType  *int    `json:"accountType"`
	BankName     string  `json:"bankName"`
	BankCode     string  `json:"bankCode"`
	BranchName   string  `json:"branchName"`
	BranchCode   string  `json:"branchCode"`
	BankAccount  string  `json:"account"`
	PayeeName    *string `json:"payeeName"`
	PayeeID      *string `json:"payeeID"`
	PayeeType    *int    `json:"payeeType"`
	FeeOwnerType int     `json:"feeOwnerType"`
}

type MediaApiMultiLanguage struct {
	AR string `json:"ar"`
	EN string `json:"en"`
	US string `json:"en_US"`
	JP string `json:"ja"`
	CN string `json:"zh_CN"`
	HK string `json:"zh_HK"`
	TW string `json:"zh_TW"`
}

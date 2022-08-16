package model

import (
	"time"

	"github.com/google/uuid"
)

func (d *CampaignLeaderboard) UUIDToBinary() *CampaignLeaderboard {
	d.BinaryId, _ = d.Id.MarshalBinary()
	return d
}

func (d *CampaignLeaderboard) BinaryToUUID() *CampaignLeaderboard {
	d.Id, _ = uuid.FromBytes(d.BinaryId)
	return d
}

func (d *CampaignRank) UUIDToBinary() *CampaignRank {
	d.BinaryLeaderboardID, _ = d.LeaderboardID.MarshalBinary()
	d.BinaryStreamerUserID, _ = d.StreamerUserID.MarshalBinary()
	return d
}

func (d *CampaignRank) BinaryToUUID() *CampaignRank {
	d.LeaderboardID, _ = uuid.FromBytes(d.BinaryLeaderboardID)
	d.StreamerUserID, _ = uuid.FromBytes(d.BinaryStreamerUserID)
	return d
}

func (d *LeaderboardDeatail) UUIDToBinary() *LeaderboardDeatail {
	d.BinaryLeaderboardID, _ = d.LeaderboardID.MarshalBinary()
	return d
}

func (d *LeaderboardDeatail) BinaryToUUID() *LeaderboardDeatail {
	d.LeaderboardID, _ = uuid.FromBytes(d.BinaryLeaderboardID)
	return d
}

func (d *CampaignBonus) ParseDateFormat() *CampaignBonus {
	date, _ := time.Parse(time.RFC3339, *d.PayDate)
	*d.PayDate = date.Format("2006-01-02")
	return d
}

type Campaign struct {
	Id             int        `json:"id" gorm:"column:id"`
	Title          string     `json:"title" gorm:"column:title"`
	RegionCode     *string    `json:"region" gorm:"column:regionCode"`
	RegionList     *string    `json:"-" gorm:"column:regionList"`
	StartTime      int64      `json:"startTime" gorm:"column:startTime"`
	EndTime        int64      `json:"endTime" gorm:"column:endTime"`
	Budget         float64    `json:"budget" gorm:"column:budget"`
	TotalBonus     float64    `json:"totalBonus" gorm:"column:totalBonus"`
	BonusDiff      float64    `json:"bonusDiff" gorm:"column:bonusDiff"`
	Remark         *string    `json:"remark" gorm:"column:remark"`
	SyncTime       time.Time  `json:"syncTime" gorm:"column:syncTime"`
	ModifyTime     *time.Time `json:"-" gorm:"column:modifyTime"`
	ModifierUID    *int       `json:"-" gorm:"column:modifierUID"`
	ModifierName   *string    `json:"-" gorm:"column:modifierName"`
	ApprovalStatus bool       `json:"approvalStatus" gorm:"column:approvalStatus"`
	ApprovalTime   *time.Time `json:"approvalTime" gorm:"column:approvalTime"`
	ApproverUID    *int       `json:"-" gorm:"column:approverUID"`
	ApproverName   *string    `json:"approverName" gorm:"column:approverName"`
	PayDate        string     `json:"-" gorm:"column:payDate"`
	TotalCount     int        `json:"-" gorm:"column:totalCount"`
}

type CampaignInsert struct {
	Id             int        `json:"id" gorm:"column:id"`
	Title          string     `json:"title" gorm:"column:title"`
	RegionCode     *string    `json:"region" gorm:"column:regionCode"`
	RegionList     *string    `json:"-" gorm:"column:regionList"`
	StartTime      int64      `json:"startTime" gorm:"column:startTime"`
	EndTime        int64      `json:"endTime" gorm:"column:endTime"`
	Budget         float64    `json:"budget" gorm:"column:budget"`
	TotalBonus     float64    `json:"totalBonus" gorm:"column:totalBonus"`
	Remark         *string    `json:"remark" gorm:"column:remark"`
	SyncTime       time.Time  `json:"syncTime" gorm:"column:syncTime"`
	ModifyTime     *time.Time `json:"-" gorm:"column:modifyTime"`
	ModifierUID    *int       `json:"-" gorm:"column:modifierUID"`
	ApprovalStatus bool       `json:"approvalStatus" gorm:"column:approvalStatus"`
	ApprovalTime   *time.Time `json:"approvalTime" gorm:"column:approvalTime"`
	ApproverUID    *int       `json:"-" gorm:"column:approverUID"`
}

type CampaignLeaderboard struct {
	Id         uuid.UUID `json:"id" gorm:"-"`
	BinaryId   []byte    `json:"-" gorm:"column:id"`
	CampaignID int       `json:"campaignID" gorm:"column:campaignID"`
	Title      string    `json:"title" gorm:"column:title"`
	SyncTime   time.Time `json:"syncTime" gorm:"column:syncTime"`
}

type CampaignRank struct {
	Id                   int       `json:"id" gorm:"column:id"`
	LeaderboardID        uuid.UUID `json:"leaderboardID" gorm:"-"`
	BinaryLeaderboardID  []byte    `json:"-" gorm:"column:leaderboardID"`
	Rank                 int       `json:"rank" gorm:"column:rank"`
	Score                int       `json:"score" gorm:"column:score"`
	StreamerUserID       uuid.UUID `json:"userID" gorm:"-"`
	BinaryStreamerUserID []byte    `json:"-" gorm:"column:userID"`
	SyncTime             time.Time `json:"syncTime" gorm:"column:syncTime"`
}

type CampaignBonus struct {
	Id          int       `json:"bonusID" gorm:"column:id"`
	RankID      int       `json:"rankID" gorm:"column:rankID"`
	BonusType   int       `json:"bonusType" gorm:"column:bonusType" validate:"oneof=0 1 2 3 4"`
	Amount      float64   `json:"amount" gorm:"column:amount"`
	PayDate     *string   `json:"payDate" gorm:"column:payDate"`
	Remark      *string   `json:"remark" gorm:"column:remark"`
	CreateTime  time.Time `json:"createTime" gorm:"column:createTime"`
	CreatorUID  int       `json:"creatorUID" gorm:"column:creatorUID"`
	ModifyTime  time.Time `json:"modifyTime" gorm:"column:modifyTime"`
	ModifierUID int       `json:"modifierUID" gorm:"column:modifierUID"`
}

type CampaignFilter struct {
	PageFilter
	Date     *time.Time
	Regions  []string
	Approval *bool
	IsZero   *bool
	Keyword  string
}

// struct for campaign detail

type CampaignDetail struct {
	Campaign        `json:",inline"`
	LeaderboardList []LeaderboardDeatail `json:"leaderboardList" gorm:"-"`
	BonusList       []CampaignBonus      `json:"-" gorm:"-"`
}

type LeaderboardDeatail struct {
	LeaderboardID       uuid.UUID    `json:"leaderboardID" gorm:"-"`
	BinaryLeaderboardID []byte       `json:"-" gorm:"column:leaderboardID"`
	Name                string       `json:"name" gorm:"column:title"`
	FixedBonus          float64      `json:"fixedBonus" gorm:"column:totalFixedBonus"`
	VariableBonus       float64      `json:"variableBonus" gorm:"column:totalVariableBonus"`
	TotalBonus          float64      `json:"totalBonus" gorm:"column:totalBonus"`
	RankList            []RankDetail `json:"rankList" gorm:"-"`
}

type RankDetail struct {
	RankID         int       `json:"rankID" gorm:"column:rankID"`
	Rank           int       `json:"rank" gorm:"column:rank"`
	Score          int       `json:"score" gorm:"column:score"`
	StreamerOpenID string    `json:"openID" gorm:"column:openID"`
	StreamerUserID uuid.UUID `json:"userID" gorm:"column:userID"`
	StreamerRegion string    `json:"userRegion" gorm:"column:streamerRegion"`
	FixedBonus     float64   `json:"fixedBonus" gorm:"column:fixedBonus"`
	VariableBonus  float64   `json:"variableBonus" gorm:"column:variableBonus"`
	TotalBonus     float64   `json:"totalBonus" gorm:"column:totalBonus"`
}

// struct for no region campaign

type CampaignBasic struct {
	Id         int    `json:"id" gorm:"column:id"`
	Title      string `json:"title" gorm:"column:title"`
	StartTime  int64  `json:"startTime" gorm:"column:startTime"`
	EndTime    int64  `json:"endTime" gorm:"column:endTime"`
	RegionList string `json:"regionList" gorm:"column:regionList"`
	TotalCount int    `json:"-" gorm:"column:totalCount"`
}

// struct for unpaid campaign bonus

type CampaignBonusUnpaid struct {
	RankID           int       `json:"rankID" gorm:"column:rankID"`
	UserID           uuid.UUID `json:"userID" gorm:"column:userID"`
	OpenID           string    `json:"openID" gorm:"column:openID"`
	CampaignID       int       `json:"campaignID" gorm:"column:campaignID"`
	CampaignTitle    string    `json:"campaignTitle" gorm:"column:campaignTitle"`
	CampaignCurrency string    `json:"campaignCurrency" gorm:"column:campaignCurrency"`
	RegionCode       string    `json:"-" gorm:"column:regionCode"`
}

package model

import (
	"time"

	"github.com/google/uuid"
)

func (d *PayoutList) ParseDateFormat() *PayoutList {
	date, _ := time.Parse(time.RFC3339, d.PayDate)
	d.PayDate = date.Format("2006-01-02")
	return d
}

func (d *PayoutReport) ParseDateFormat() *PayoutReport {
	date, _ := time.Parse(time.RFC3339, d.PayDate)
	d.PayDate = date.Format("2006-01-02")
	return d
}

func (d *PayoutGroup) ParseBonusIdList() *PayoutGroup {
	for _, v := range d.IdList {
		d.Ids = append(d.Ids, v.Id)
	}

	return d
}

type Payout struct {
	Id          int       `json:"id" gorm:"column:id"`
	RegionCode  string    `json:"regionCode" gorm:"column:regionCode"`
	PayDate     string    `json:"payDate" gorm:"column:payDate"`
	PayStatus   bool      `json:"status" gorm:"column:payStatus"`
	CreateTime  time.Time `json:"createTime" gorm:"column:createTime"`
	CreatorUID  int       `json:"creatorUID" gorm:"column:creatorUID"`
	ModifyTime  time.Time `json:"modifyTime" gorm:"column:modifyTime"`
	ModifierUID int       `json:"modifierUID" gorm:"column:modifierUID"`
}

type PayoutDetail struct {
	PayoutID   int       `json:"payoutID" gorm:"column:payoutID"`
	BonusID    int       `json:"bonusID" gorm:"column:bonusID"`
	CreateTime time.Time `json:"-" gorm:"column:createTime"`
	CreatorUID int       `json:"-" gorm:"column:creatorUID"`
}

type PayoutRecord struct {
	PayoutID      int       `json:"payoutID" gorm:"column:payoutID"`
	UserID        []byte    `json:"userID" gorm:"column:userID"`
	AgencyID      *int      `json:"agencyID" gorm:"column:agencyID"`
	AccountID     *int      `json:"accountID" gorm:"column:accountID"`
	PayType       *string   `json:"payType" gorm:"column:payType"`
	PayAmount     float64   `json:"payAmount" gorm:"column:payAmount"`
	TaxableAmount float64   `json:"taxableAmount" gorm:"column:taxableAmount"`
	TaxRate       *float64  `json:"taxRate" gorm:"column:taxRate"`
	TaxAmount     *float64  `json:"taxAmount" gorm:"column:taxAmount"`
	TaxID         string    `json:"taxID" gorm:"column:taxID"`
	FeeOwnerType  *int      `json:"feeOwnerType" gorm:"column:feeOwnerType"`
	CreateTime    time.Time `json:"createTime" gorm:"column:createTime"`
	CreatorUID    int       `json:"creatorUID" gorm:"column:creatorUID"`
}

type PayoutList struct {
	PayoutID       int     `json:"id" gorm:"column:payoutID"`
	PayDate        string  `json:"payDate" gorm:"column:payDate"`
	RegionCode     string  `json:"region" gorm:"column:regionCode"`
	FixedBonus     float64 `json:"fixedBonus" gorm:"column:fixedBonus"`
	VariableBonus  float64 `json:"variableBonus" gorm:"column:variableBonus"`
	Transportation float64 `json:"transportation" gorm:"column:transportation"`
	Addon          float64 `json:"addon" gorm:"column:addon"`
	Deduction      float64 `json:"deduction" gorm:"column:deduction"`
	Total          float64 `json:"total" gorm:"column:total"`
	Budget         float64 `json:"budget" gorm:"column:totalBudget"`
	Difference     float64 `json:"difference" gorm:"column:difference"`
	IsPaid         bool    `json:"isPaid" gorm:"column:isPaid"`
	Status         int     `json:"status" gorm:"column:status"`
	TotalCount     int     `json:"-" gorm:"column:totalCount"`
}

type PayoutListDetail struct {
	OpenID         string    `json:"openID" gorm:"column:openID"`
	UserID         uuid.UUID `json:"userID" gorm:"column:userID"`
	RegionCode     string    `json:"-" gorm:"column:regionCode"`
	StreamerRegion string    `json:"region" gorm:"column:streamerRegion"`
	CampaignID     int       `json:"campaignID" gorm:"column:campaignID"`
	CampaignTitle  string    `json:"campaignTitle" gorm:"column:campaignTitle"`
	BonusID        int       `json:"bonusID" gorm:"column:bonusID"`
	BonusType      int       `json:"bonusType" gorm:"column:bonusType"`
	Amount         float64   `json:"amount" gorm:"column:amount"`
	Remark         string    `json:"remark" gorm:"column:remark"`
	TotalCount     int       `json:"-" gorm:"column:totalCount"`
}

type PayoutReport struct {
	PayoutID       int     `json:"id" gorm:"column:payoutID"`
	PayDate        string  `json:"payDate" gorm:"column:payDate"`
	RegionCode     string  `json:"region" gorm:"column:regionCode"`
	ForeignAmount  float64 `json:"foreignAmount" gorm:"column:foreignAmount"`
	LocalAmount    float64 `json:"localAmount" gorm:"column:localAmount"`
	AgencyAmount   float64 `json:"agencyAmount" gorm:"column:agencyAmount"`
	StreamerCount  int64   `json:"streamerCount" gorm:"column:streamerCount"`
	MissingCount   int64   `json:"missingCount" gorm:"column:missingCount"`
	TotalAmount    float64 `json:"totalAmount" gorm:"column:totalAmount"`
	TotalTaxAmount float64 `json:"totalTaxAmount" gorm:"column:totalTaxAmount"`
	AfterTaxAmount float64 `json:"afterTaxAmount" gorm:"column:afterTaxAmount"`
	PayStatus      bool    `json:"status" gorm:"column:payStatus"`
	TotalCount     int     `json:"-" gorm:"column:totalCount"`
}

type PayoutReportDetail struct {
	UserID          uuid.UUID `json:"userID" gorm:"column:userID"`
	OpenID          *string   `json:"openID" gorm:"column:openID"`
	StreamerRegion  *string   `json:"streamerRegion" gorm:"column:streamerRegion"`
	PayType         *string   `json:"payType" gorm:"column:payType"`
	PayAmount       float64   `json:"payAmount" gorm:"column:payAmount"`
	TaxAmount       float64   `json:"taxAmount" gorm:"column:taxAmount"`
	AfterTaxAmount  float64   `json:"afterTaxAmount" gorm:"column:afterTaxAmount"`
	BankCode        *string   `json:"bankCode" gorm:"column:bankCode"`
	BankName        *string   `json:"bankName" gorm:"column:bankName"`
	BranchCode      *string   `json:"branchCode" gorm:"column:branchCode"`
	BranchName      *string   `json:"branchName" gorm:"column:branchName"`
	BankAccount     *string   `json:"bankAccount" gorm:"column:accountNo"`
	BankAccountName *string   `json:"bankAccountName" gorm:"column:accountName"`
	FeeOwnerType    *int      `json:"feeOwnerType" gorm:"column:feeOwnerType"`
	TaxID           string    `json:"taxID" gorm:"column:taxID"`
	TotalCount      int       `json:"-" gorm:"column:totalCount"`
}

type PayoutReportDetailByAgency struct {
	AgencyName      *string                `json:"agencyName" gorm:"column:agencyName"`
	AgencyId        string                 `json:"agencyID" gorm:"column:agencyID"`
	TotalAmount     float64                `json:"totalAmount" gorm:"column:totalAmount"`
	TotalTaxAmount  float64                `json:"totalTaxAmount" gorm:"column:totalTaxAmount"`
	AfterTaxAmount  float64                `json:"afterTaxAmount" gorm:"column:afterTaxAmount"`
	BankCode        *string                `json:"bankCode" gorm:"column:bankCode"`
	BankName        *string                `json:"bankName" gorm:"column:bankName"`
	BranchCode      *string                `json:"branchCode" gorm:"column:branchCode"`
	BranchName      *string                `json:"branchName" gorm:"column:branchName"`
	BankAccount     *string                `json:"bankAccount" gorm:"column:accountNo"`
	BankAccountName *string                `json:"bankAccountName" gorm:"column:accountName"`
	FeeOwnerType    *int                   `json:"feeOwnerType" gorm:"column:feeOwnerType"`
	TaxID           string                 `json:"taxID" gorm:"column:taxID"`
	TotalCount      int                    `json:"-" gorm:"column:totalCount"`
	StreamerList    []PayoutReportStreamer `json:"streamerList" gorm:"-"`
}

type PayoutReportStreamer struct {
	UserID         uuid.UUID `json:"userID" gorm:"column:userID"`
	OpenID         *string   `json:"openID" gorm:"column:openID"`
	StreamerRegion *string   `json:"streamerRegion" gorm:"column:streamerRegion"`
	CampaignTitle  string    `json:"campaignTitle" gorm:"column:campaignTitle"`
	PayAmount      float64   `json:"payAmount" gorm:"column:payAmount"`
}

type PayoutFilter struct {
	PageFilter
	Regions   []string
	Date      *time.Time
	IsPaid    *bool
	IsGrouped *bool
	BonusType *int
	PayType   string
	Keyword   string
}

type PayoutGroup struct {
	RegionCode string `json:"region"`
	PayDate    string `json:"payDate"`
	IdList     []struct {
		Id int `json:"id"`
	} `json:"idList"`
	Ids []int
}

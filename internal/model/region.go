package model

import "time"

type Region struct {
	Code         string  `json:"code" gorm:"column:code"`
	Name         string  `json:"name" gorm:"column:name"`
	CurrencyCode *string `json:"-" gorm:"column:currencyCode"`
	Show         bool    `json:"-" gorm:"column:show"`
}

type RegionDetail struct {
	Code           string  `json:"code" gorm:"column:code"`
	Name           string  `json:"name" gorm:"column:name"`
	CurrencyCode   *string `json:"currencyCode" gorm:"column:currencyCode"`
	CurrencyName   *string `json:"currencyName" gorm:"column:currencyName"`
	CurrencyFormat *int    `json:"currencyFormat" gorm:"column:currencyFormat"`
	PayType        string  `json:"-" gorm:"column:payType"`
	TaxFrom        int     `json:"-" gorm:"column:taxFrom"`
	TaxRate        float64 `json:"-" gorm:"column:taxRate"`
}

type TaxRate struct {
	Id          int       `json:"-" gorm:"column:id"`
	RegionCode  string    `json:"-" gorm:"column:regionCode"`
	PayType     string    `json:"payType" gorm:"column:payType" validate:"oneof=F L A"`
	TaxFrom     int       `json:"taxFrom" gorm:"column:taxFrom"`
	TaxRate     float64   `json:"taxRate" gorm:"column:taxRate" validate:"lt=1"`
	CreateTime  time.Time `json:"-" gorm:"column:createTime"`
	CreatorUID  int       `json:"-" gorm:"column:creatorUID"`
	ModifyTime  time.Time `json:"-" gorm:"column:modifyTime"`
	ModifierUID int       `json:"-" gorm:"column:modifierUID"`
}

type RegionWithTaxList struct {
	RegionDetail `json:",inline"`
	TaxList      []TaxRate `json:"taxList" validate:"dive"`
}

package model

import "time"

type BankAccount struct {
	Id              int       `json:"id" gorm:"column:id"`
	SwiftCode       string    `json:"swiftCode" gorm:"column:swiftCode"`
	BankCode        string    `json:"bankCode" gorm:"column:bankCode"`
	BankName        string    `json:"bankName" gorm:"column:bankName"`
	BranchCode      string    `json:"branchCode" gorm:"column:branchCode"`
	BranchName      string    `json:"branchName" gorm:"column:branchName"`
	BankAccount     string    `json:"bankAccount" gorm:"column:accountNo"`
	BankAccountName string    `json:"bankAccountName" gorm:"column:accountName"`
	BankAccountType int       `json:"bankAccountType" gorm:"column:accountType"`
	PayeeType       int       `json:"payeeType" gorm:"column:payeeType"`
	PayeeID         *string   `json:"payeeID" gorm:"column:payeeID"`
	OwnerType       *int      `json:"ownerType" gorm:"column:ownerType"`
	FeeOwnerType    *int      `json:"feeOwnerType" gorm:"column:feeOwnerType"`
	SyncTime        time.Time `json:"syncTime" gorm:"column:syncTime"`
}

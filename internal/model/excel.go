package model

import (
	"github.com/google/uuid"
)

type PayoutReportExcel struct {
	Region                  string
	PayMonth                string
	SheetLocalTop           []map[string]string
	SheetLocalData          []map[string]interface{}
	SheetForeignTop         []map[string]string
	SheetForeignData        []map[string]interface{}
	SheetAgencyTop          []map[string]string
	SheetAgencyData         []map[string]interface{}
	SheetTransportationTop  []map[string]string
	SheetTransportationData []map[string]interface{}
	SheetMissingTop         []map[string]string
	SheetMissingData        []map[string]interface{}
}

type PayoutReportDetailForExcel struct {
	UserID          uuid.UUID `json:"userID" gorm:"column:userID"`
	OpenID          string    `json:"openID" gorm:"column:openID"`
	StreamerRegion  string    `json:"streamerRegion" gorm:"column:streamerRegion"`
	PayType         string    `json:"payType" gorm:"column:payType"`
	PayAmount       float64   `json:"payAmount" gorm:"column:payAmount"`
	TaxAmount       float64   `json:"taxAmount" gorm:"column:taxAmount"`
	AfterTaxAmount  float64   `json:"afterTaxAmount" gorm:"column:afterTaxAmount"`
	BankCode        string    `json:"bankCode" gorm:"column:bankCode"`
	BankName        string    `json:"bankName" gorm:"column:bankName"`
	BranchCode      string    `json:"branchCode" gorm:"column:branchCode"`
	BranchName      string    `json:"branchName" gorm:"column:branchName"`
	BankAccount     string    `json:"bankAccount" gorm:"column:accountNo"`
	BankAccountName string    `json:"bankAccountName" gorm:"column:accountName"`
	FeeOwnerType    int       `json:"feeOwnerType" gorm:"column:feeOwnerType"`
	TaxID           string    `json:"taxID" gorm:"column:taxID"`
}

type PayoutReportDetailByAgencyForExcel struct {
	AgencyName      string  `json:"agencyName" gorm:"column:agencyName"`
	AgencyId        string  `json:"agencyID" gorm:"column:agencyID"`
	TotalAmount     float64 `json:"totalAmount" gorm:"column:totalAmount"`
	TotalTaxAmount  float64 `json:"totalTaxAmount" gorm:"column:totalTaxAmount"`
	AfterTaxAmount  float64 `json:"afterTaxAmount" gorm:"column:afterTaxAmount"`
	BankCode        string  `json:"bankCode" gorm:"column:bankCode"`
	BankName        string  `json:"bankName" gorm:"column:bankName"`
	BranchCode      string  `json:"branchCode" gorm:"column:branchCode"`
	BranchName      string  `json:"branchName" gorm:"column:branchName"`
	BankAccount     string  `json:"bankAccount" gorm:"column:accountNo"`
	BankAccountName string  `json:"bankAccountName" gorm:"column:accountName"`
	FeeOwnerType    int     `json:"feeOwnerType" gorm:"column:feeOwnerType"`
	TaxID           string  `json:"taxID" gorm:"column:taxID"`
}

func (d *PayoutReportExcel) GenerateExcelTop() *PayoutReportExcel {
	agencyTop := make([]map[string]string, 0)
	agencyTop = append(agencyTop, map[string]string{"key": "agencyName", "title": "Agency Name", "width": "30", "is_num": "0"})
	agencyTop = append(agencyTop, map[string]string{"key": "agencyID", "title": "Agency ID", "width": "20", "is_num": "0"})
	agencyTop = append(agencyTop, map[string]string{"key": "amount", "title": "Amount", "width": "20", "is_num": "0"})
	agencyTop = append(agencyTop, map[string]string{"key": "payAmount", "title": "Pay Amount", "width": "20", "is_num": "0"})
	agencyTop = append(agencyTop, map[string]string{"key": "payoutType", "title": "Payout Type", "width": "20", "is_num": "0"})
	agencyTop = append(agencyTop, map[string]string{"key": "bankName", "title": "Bank Name", "width": "30", "is_num": "0"})
	agencyTop = append(agencyTop, map[string]string{"key": "bankCode", "title": "Bank Code", "width": "20", "is_num": "0"})
	agencyTop = append(agencyTop, map[string]string{"key": "branchName", "title": "Branch Name", "width": "30", "is_num": "0"})
	agencyTop = append(agencyTop, map[string]string{"key": "branchCode", "title": "Branch Code", "width": "20", "is_num": "0"})
	agencyTop = append(agencyTop, map[string]string{"key": "bankAccount", "title": "Bank Account", "width": "30", "is_num": "0"})
	agencyTop = append(agencyTop, map[string]string{"key": "payeeName", "title": "Payee Name", "width": "30", "is_num": "0"})
	agencyTop = append(agencyTop, map[string]string{"key": "feeOwner", "title": "Fee Owner", "width": "20", "is_num": "0"})
	agencyTop = append(agencyTop, map[string]string{"key": "taxID", "title": "Tax ID", "width": "20", "is_num": "0"})

	detailTop := make([]map[string]string, 0)
	detailTop = append(detailTop, map[string]string{"key": "userID", "title": "User ID", "width": "40", "is_num": "0"})
	detailTop = append(detailTop, map[string]string{"key": "openID", "title": "ID", "width": "30", "is_num": "0"})
	detailTop = append(detailTop, map[string]string{"key": "amount", "title": "Amount", "width": "20", "is_num": "0"})
	detailTop = append(detailTop, map[string]string{"key": "taxAmount", "title": "Tax Amount", "width": "20", "is_num": "0"})
	detailTop = append(detailTop, map[string]string{"key": "payAmount", "title": "Pay Amount", "width": "20", "is_num": "0"})
	detailTop = append(detailTop, map[string]string{"key": "payoutType", "title": "Payout Type", "width": "20", "is_num": "0"})
	detailTop = append(detailTop, map[string]string{"key": "bankName", "title": "Bank Name", "width": "30", "is_num": "0"})
	detailTop = append(detailTop, map[string]string{"key": "bankCode", "title": "Bank Code", "width": "20", "is_num": "0"})
	detailTop = append(detailTop, map[string]string{"key": "branchName", "title": "Branch Name", "width": "30", "is_num": "0"})
	detailTop = append(detailTop, map[string]string{"key": "branchCode", "title": "Branch Code", "width": "20", "is_num": "0"})
	detailTop = append(detailTop, map[string]string{"key": "bankAccount", "title": "Bank Account", "width": "30", "is_num": "0"})
	detailTop = append(detailTop, map[string]string{"key": "payeeName", "title": "Payee Name", "width": "30", "is_num": "0"})
	detailTop = append(detailTop, map[string]string{"key": "feeOwner", "title": "Fee Owner", "width": "20", "is_num": "0"})
	detailTop = append(detailTop, map[string]string{"key": "taxID", "title": "Tax ID", "width": "20", "is_num": "0"})

	d.SheetLocalTop = detailTop
	d.SheetForeignTop = detailTop
	d.SheetAgencyTop = agencyTop
	d.SheetTransportationTop = detailTop
	d.SheetMissingTop = detailTop

	return d
}

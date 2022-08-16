package repository

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"context"
)

func (c *Client) GetBankAccount(ctx context.Context, swiftCode, bankCode, branchCode, bankAccount string) ([]model.BankAccount, error) {
	var data []model.BankAccount

	if err := c.Database.Table(BankAccount).Find(&data, "swiftCode = ? AND bankCode = ? AND branchCode = ? AND accountNo = ?", swiftCode, bankCode, branchCode, bankAccount).Error; err != nil {
		log.Errorf("fail to get bank account: %s", err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) CreateBankAccount(ctx context.Context, data model.BankAccount) (model.BankAccount, error) {
	if err := c.Database.Table(BankAccount).Create(&data).Error; err != nil {
		log.Errorf("fail to create bank account: %v. %s", data, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	// log.Infof("bank account created: %d", data.Id)

	return data, nil
}

func (c *Client) UpdateBankAccount(ctx context.Context, data model.BankAccount) error {
	if err := c.Database.Table(BankAccount).Save(&data).Error; err != nil {
		log.Errorf("fail to create bank account: %v. %s", data, err.Error())
		return customError.New(customError.DatabaseError)
	}

	// log.Infof("bank account updated: %d", data.Id)

	return nil
}

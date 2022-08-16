package service

import (
	"17live_wso_be/internal/model"
	"context"
)

func (c *Client) GetBankAccount(ctx context.Context, swiftCode, bankCode, branchCode, bankAccount string) ([]model.BankAccount, error) {
	return c.RepositoryClient.GetBankAccount(ctx, swiftCode, bankCode, branchCode, bankAccount)
}

func (c *Client) CreateBankAccount(ctx context.Context, data model.BankAccount) (model.BankAccount, error) {
	return c.RepositoryClient.CreateBankAccount(ctx, data)
}

func (c *Client) UpdateBankAccount(ctx context.Context, data model.BankAccount) error {
	return c.RepositoryClient.UpdateBankAccount(ctx, data)
}

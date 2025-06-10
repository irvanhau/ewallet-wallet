package repository

import (
	"context"
	"ewallet-wallet/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type WalletRepository struct {
	DB *gorm.DB
}

func (r *WalletRepository) CreateWallet(ctx context.Context, wallet *models.Wallet) error {
	return r.DB.Create(wallet).Error
}

func (r *WalletRepository) UpdateBalance(ctx context.Context, userID uint, amount float64) (models.Wallet, error) {
	var wallet models.Wallet
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Raw("SELECT id,user_id,balance FROM wallets WHERE user_id = ? FOR UPDATE", userID).Scan(&wallet).Error
		if err != nil {
			return err
		}

		if (wallet.Balance + amount) < 0 {
			return fmt.Errorf("current balance is not enough to perform the transaction: %f - %f", wallet.Balance, amount)
		}

		err = tx.Exec("UPDATE wallets SET balance = balance + ? WHERE user_id = ?", int(amount), userID).Error
		if err != nil {
			return err
		}
		return nil
	})

	return wallet, err
}

func (r *WalletRepository) CreateWalletTrx(ctx context.Context, walletTrx *models.WalletTransaction) error {
	return r.DB.Create(walletTrx).Error
}

func (r *WalletRepository) GetWalletTransactionByReference(ctx context.Context, reference string) (models.WalletTransaction, error) {
	var (
		resp models.WalletTransaction
	)

	err := r.DB.Where("reference = ?", reference).Last(&resp).Error

	return resp, err
}

func (r *WalletRepository) GetWalletByUserID(ctx context.Context, userID uint) (models.Wallet, error) {
	var (
		resp models.Wallet
	)

	err := r.DB.Where("user_id =  ?", userID).Last(&resp).Error

	return resp, err
}

func (r *WalletRepository) GetWalletHistory(ctx context.Context, walletID uint, offset int, limit int, transactionType string) ([]models.WalletTransaction, error) {
	var (
		resp []models.WalletTransaction
	)

	sql := r.DB
	if transactionType != "" {
		sql = sql.Where("wallet_transaction_type = ?", transactionType)
	}
	err := sql.Limit(limit).Offset(offset).Order("id DESC").Find(&resp).Error

	return resp, err
}

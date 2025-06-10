package services

import (
	"context"
	"ewallet-wallet/internal/interfaces"
	"ewallet-wallet/internal/models"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type WalletService struct {
	WalletRepository interfaces.IWalletRepo
}

func (s *WalletService) Create(ctx context.Context, wallet *models.Wallet) error {
	return s.WalletRepository.CreateWallet(ctx, wallet)
}

func (s *WalletService) CreditBalance(ctx context.Context, userID uint, req models.TransactionRequest) (models.BalanceResponse, error) {
	var (
		resp models.BalanceResponse
	)

	history, err := s.WalletRepository.GetWalletTransactionByReference(ctx, req.Reference)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return resp, errors.Wrap(err, "failed to check reference")
		}
	}

	if history.ID > 0 {
		return resp, errors.New("reference is duplicated")
	}

	wallet, err := s.WalletRepository.UpdateBalance(ctx, userID, req.Amount)
	if err != nil {
		return resp, errors.Wrap(err, "failed to update balance")
	}

	walletHistory := &models.WalletTransaction{
		WalletID:              wallet.ID,
		Amount:                req.Amount,
		WalletTransactionTyoe: "CREDIT",
		Reference:             req.Reference,
	}

	err = s.WalletRepository.CreateWalletTrx(ctx, walletHistory)
	if err != nil {
		return resp, errors.Wrap(err, "failed to insert wallet transaction")
	}

	resp.Balance = wallet.Balance + req.Amount

	return resp, nil
}

func (s *WalletService) DebitBalance(ctx context.Context, userID uint, req models.TransactionRequest) (models.BalanceResponse, error) {
	var (
		resp models.BalanceResponse
	)

	history, err := s.WalletRepository.GetWalletTransactionByReference(ctx, req.Reference)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return resp, errors.Wrap(err, "failed to check reference")
		}
	}

	if history.ID > 0 {
		return resp, errors.New("reference is duplicated")
	}

	wallet, err := s.WalletRepository.UpdateBalance(ctx, userID, -req.Amount)
	if err != nil {
		return resp, errors.Wrap(err, "failed to update balance")
	}

	walletHistory := &models.WalletTransaction{
		WalletID:              wallet.ID,
		Amount:                req.Amount,
		WalletTransactionTyoe: "DEBIT",
		Reference:             req.Reference,
	}

	err = s.WalletRepository.CreateWalletTrx(ctx, walletHistory)
	if err != nil {
		return resp, errors.Wrap(err, "failed to insert wallet transaction")
	}

	resp.Balance = wallet.Balance + req.Amount

	return resp, nil
}

func (s *WalletService) GetBalance(ctx context.Context, userID uint) (models.BalanceResponse, error) {
	var (
		resp models.BalanceResponse
	)
	wallet, err := s.WalletRepository.GetWalletByUserID(ctx, userID)
	if err != nil {
		return resp, errors.Wrap(err, "failed to get wallet")
	}

	resp.Balance = wallet.Balance

	return resp, nil
}

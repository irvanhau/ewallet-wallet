package services

import (
	"context"
	"ewallet-wallet/internal/interfaces"
	"ewallet-wallet/internal/models"
)

type WalletService struct {
	WalletRepository interfaces.IWalletRepo
}

func (s *WalletService) Create(ctx context.Context, wallet *models.Wallet) error {
	return s.WalletRepository.CreateWallet(ctx, wallet)
}

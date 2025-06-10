package cmd

import (
	"ewallet-wallet/external"
	"ewallet-wallet/helpers"
	"ewallet-wallet/internal/api"
	"ewallet-wallet/internal/interfaces"
	"ewallet-wallet/internal/repository"
	"ewallet-wallet/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func ServeHTTP() {
	d := dependencyInject()

	r := gin.Default()

	r.GET("/health", d.HealthCheckAPI.HealthCheckHandlerHTTP)

	walletV1 := r.Group("/wallet/v1")
	walletV1.POST("/", d.WalletAPI.Create)
	walletV1.PUT("/balance/credit", d.MiddlewareValidateToken, d.WalletAPI.CreditBalance)
	walletV1.PUT("/balance/debit", d.MiddlewareValidateToken, d.WalletAPI.DebitBalance)
	walletV1.GET("/balance", d.MiddlewareValidateToken, d.WalletAPI.GetBalance)
	walletV1.GET("/history", d.MiddlewareValidateToken, d.WalletAPI.GetWalletHistory)

	err := r.Run(":" + helpers.GetEnv("PORT", "8080"))
	if err != nil {
		log.Fatal(err)
	}
}

type Dependency struct {
	HealthCheckAPI interfaces.IHealthCheckAPI
	WalletAPI      interfaces.IWalletAPI
	External       interfaces.IExternal
}

func dependencyInject() Dependency {
	healthCheckSvc := &services.HealthCheck{}
	healthCheckAPI := &api.HealthCheck{
		HealthCheckServices: healthCheckSvc,
	}

	walletRepo := &repository.WalletRepository{
		DB: helpers.DB,
	}
	walletSvc := &services.WalletService{
		WalletRepository: walletRepo,
	}
	walletAPI := &api.WalletAPI{
		WalletService: walletSvc,
	}

	external := &external.External{}

	return Dependency{
		HealthCheckAPI: healthCheckAPI,
		WalletAPI:      walletAPI,
		External:       external,
	}
}

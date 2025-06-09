package cmd

import (
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

	err := r.Run(":" + helpers.GetEnv("PORT", "8080"))
	if err != nil {
		log.Fatal(err)
	}
}

type Dependency struct {
	HealthCheckAPI interfaces.IHealthCheckAPI
	WalletAPI      interfaces.IWalletAPI
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

	return Dependency{
		HealthCheckAPI: healthCheckAPI,
		WalletAPI:      walletAPI,
	}
}

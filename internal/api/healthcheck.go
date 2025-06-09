package api

import (
	"ewallet-wallet/helpers"
	"ewallet-wallet/internal/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheck struct {
	HealthCheckServices interfaces.IHealthCheckServices
}

func (api *HealthCheck) HealthCheckHandlerHTTP(c *gin.Context) {
	msg, err := api.HealthCheckServices.HealthCheckServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, msg, nil)
}

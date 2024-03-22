package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/wonwooseo/panawa-api/build"
)

type AdminController struct {
	logger zerolog.Logger
}

func NewAdminController(baseLogger zerolog.Logger) *AdminController {
	return &AdminController{
		logger: baseLogger.With().Str("caller", "controller/admin").Logger(),
	}
}

func (c *AdminController) HealthCheckEndpoint(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
}

func (c *AdminController) VersionEndpoint(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Version{
		Version:   build.Version,
		BuildTime: build.BuildTime,
	})
}

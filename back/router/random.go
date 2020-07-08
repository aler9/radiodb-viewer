package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"rdbviewer/shared"
)

func (h *router) onDataRandom(ctx *gin.Context) {
	res, err := h.dbClient.BootlegRand(context.Background(), &shared.BootlegRandReq{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"path": "/bootlegs/" + res.Id,
	})
}

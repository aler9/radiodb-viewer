package main

import (
	"context"

	"github.com/gin-gonic/gin"

	"rdbviewer/shared"
)

func (h *router) onPageRandom(c *gin.Context) {
	res, err := h.dbClient.BootlegRand(context.Background(), &shared.BootlegRandReq{})
	if err != nil {
		GinServerErrorText(c)
		return
	}
	c.Redirect(302, "/bootlegs/"+res.Id)
}

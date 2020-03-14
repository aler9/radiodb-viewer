package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"rdbviewer/shared"
)

func (h *router) onPageRandom(c *gin.Context) {
	res, err := h.dbClient.BootlegRand(context.Background(), &shared.BootlegRandReq{})
	if err != nil {
		http.Error(c.Writer, "500 internal server error", http.StatusInternalServerError)
		return
	}

	c.Redirect(302, "/bootlegs/"+res.Id)
}

package main

import (
    "context"
    "rdbviewer/back/shared"
    "github.com/gin-gonic/gin"
)

func (h *Router) onPageRandom(c *gin.Context) {
    res,err := h.dbClient.BootlegRand(context.Background(), &shared.BootlegRandReq{})
    if err != nil {
        GinServerErrorText(c)
        return
    }
    c.Redirect(302, "/bootlegs/" + res.Item.Id)
}

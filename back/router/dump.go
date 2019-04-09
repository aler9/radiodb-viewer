package main

import (
    "github.com/gin-gonic/gin"
)

func (h *Router) onPageDump(c *gin.Context) {
    GinTpl(c, h.frameWrapper(c, FrameConf{
        Title: "Dump",
        Content: TplExecute(h.templates["dump"], gin.H{}),
    }))
}

func(h *Router) onPageDumpGet(c *gin.Context) {
    c.Header("Content-Type", "application/json; charset=utf-8")
    c.Header("Content-Disposition", "attachment; filename=\"radiodb.json\"")
    c.File("/data/radiodb.json")
}

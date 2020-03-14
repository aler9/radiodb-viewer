package main

import (
	"github.com/gin-gonic/gin"
)

func (h *router) onPageDump(c *gin.Context) {
	GinTpl(c, h.frameWrapper(c, FrameConf{
		Title:   "Dump",
		Content: TplRender(h.templates["dump"], gin.H{}),
	}))
}

func (h *router) onPageDumpGet(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=\"radiodb.json\"")
	c.File("/data/radiodb.json")
}

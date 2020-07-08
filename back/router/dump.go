package main

import (
	"github.com/gin-gonic/gin"
)

func (h *router) onPageDumpGet(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.Header("Content-Disposition", "attachment; filename=\"radiodb.json\"")
	ctx.File("/data/radiodb.json")
}

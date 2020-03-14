package main

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

type FrameConf struct {
	Title   string
	Class   string
	Content string
}

func (h *router) frameWrapper(c *gin.Context, conf FrameConf) string {
	title := "RadioDB"
	if conf.Title != "" {
		title += " - " + conf.Title
	} else {
		title += " | The complete Radiohead bootlegs database"
	}

	return TplRender(h.templates["frame"], gin.H{
		"CurPath": c.Request.URL.Path,
		"Title":   title,
		"Class":   conf.Class,
		"Content": template.HTML(conf.Content),
	})
}

package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
)

type FrameConf struct {
	Title   string
	Class   string
	Content string
}

func (h *Router) frameWrapper(c *gin.Context, conf FrameConf) string {
	title := "RadioDB"
	if conf.Title != "" {
		title += " - " + conf.Title
	} else {
		title += " | The complete Radiohead bootlegs database"
	}

	return TplExecute(h.templates["frame"], gin.H{
		"CurPath": c.Request.URL.Path,
		"Title":   title,
		"Class":   conf.Class,
		"Content": template.HTML(conf.Content),
	})
}

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"rdbviewer/shared"
)

func (h *Router) onDataSearch(c *gin.Context) {
	var in struct {
		Query string
	}
	if err := GinPostBody(c, &in); err != nil {
		GinServerErrorJson(c)
		return
	}

	if len(in.Query) < 1 || len(in.Query) > 255 {
		GinJson(c, gin.H{"res": nil})
		return
	}

	res, err := h.dbClient.Search(context.Background(), &shared.SearchReq{Query: in.Query})
	if err != nil {
		GinServerErrorJson(c)
		return
	}

	GinJson(c, gin.H{
		"res": func() (ret []gin.H) {
			for _, item := range res.Items {
				ret = append(ret, gin.H{
					"url":   item.Url,
					"title": item.Title,
				})
			}
			return ret
		}(),
	})
}

package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"rdbviewer/shared"
)

func (h *router) onDataSearch(c *gin.Context) {
	var in struct {
		Query string
	}
	if err := json.NewDecoder(c.Request.Body).Decode(&in); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	if len(in.Query) < 1 || len(in.Query) > 255 {
		c.JSON(http.StatusOK, gin.H{"res": nil})
		return
	}

	res, err := h.dbClient.Search(context.Background(), &shared.SearchReq{Query: in.Query})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
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

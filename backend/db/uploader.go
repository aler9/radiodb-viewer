package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (db *database) onUpload(c *gin.Context) {
	byts, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = ioutil.WriteFile(DB_FILE, byts, 0644)
	if err != nil {
		panic(err)
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.load()
}

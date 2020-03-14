package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (db *database) onUpload(c *gin.Context) {
	mpf, err := c.MultipartForm()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	files, ok := mpf.File["file"]
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	if len(files) != 1 {
		c.Status(http.StatusBadRequest)
		return
	}
	file := files[0]

	err = c.SaveUploadedFile(file, DB_FILE)
	if err != nil {
		panic(err)
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.load()
	log.Printf("[db] data reloaded")
}

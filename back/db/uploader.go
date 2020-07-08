package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (db *database) onUpload(ctx *gin.Context) {
	mpf, err := ctx.MultipartForm()
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	files, ok := mpf.File["file"]
	if !ok {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if len(files) != 1 {
		ctx.Status(http.StatusBadRequest)
		return
	}
	file := files[0]

	err = ctx.SaveUploadedFile(file, DB_FILE)
	if err != nil {
		panic(err)
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.load()
	log.Printf("[db] data reloaded")
}

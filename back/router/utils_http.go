package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type TplMap map[string]*template.Template

func TplLoadAll(dpath string) TplMap {
	ret := make(TplMap)
	files, err := ioutil.ReadDir(dpath)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		name := f.Name()
		name = name[:len(name)-len(filepath.Ext(name))]
		fpath := filepath.Join(dpath, f.Name())
		ret[name] = template.Must(template.New(filepath.Base(fpath)).ParseFiles(fpath))
	}
	return ret
}

func TplRender(tpl *template.Template, data interface{}) string {
	var buf strings.Builder
	if err := tpl.Execute(&buf, data); err != nil {
		panic(err)
	}
	return buf.String()
}

func GinPostBody(c *gin.Context, target interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(target)
}

func GinNotFoundText(c *gin.Context) {
	c.Abort()
	http.NotFound(c.Writer, c.Request)
}

func GinNotFoundJson(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
}

func GinServerErrorText(c *gin.Context) {
	c.Abort()
	http.Error(c.Writer, "500 internal server error", http.StatusInternalServerError)
}

func GinServerErrorJson(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
}

func GinTpl(c *gin.Context, rendered string) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteString(rendered)
}

func GinJson(c *gin.Context, cnt interface{}) {
	c.JSON(http.StatusOK, cnt)
}

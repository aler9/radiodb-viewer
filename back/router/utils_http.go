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

func TplLoadAll(dpath string) map[string]*template.Template {
	ret := make(map[string]*template.Template)
	files, err := ioutil.ReadDir(dpath)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		name := f.Name()
		name = name[:len(name)-len(filepath.Ext(name))]
		ret[name] = TplLoad(filepath.Join(dpath, f.Name()))
	}
	return ret
}

func TplLoad(fpath string) *template.Template {
	return template.Must(template.New(filepath.Base(fpath)).ParseFiles(fpath))
}

func TplExecute(tpl *template.Template, data interface{}) string {
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

package main

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type TplMap map[string]*template.Template

func templateLoadAll(dpath string) TplMap {
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

func templateRender(tpl *template.Template, data interface{}) string {
	var buf strings.Builder
	if err := tpl.Execute(&buf, data); err != nil {
		panic(err)
	}
	return buf.String()
}

func ginTemplate(c *gin.Context, rendered string) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteString(rendered)
}

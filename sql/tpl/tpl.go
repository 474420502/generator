package tpl

import (
	"bytes"
	"io"
	"path/filepath"
	"runtime"
	"text/template"
)

var defaultTemplate *template.Template

func init() {
	var err error

	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)

	// 解析模板文件
	defaultTemplate, err = template.ParseGlob(filepath.Join(currentDir, "/*.tpl"))
	if err != nil {
		panic(err)
	}
}

func ExecuteTemplate(wr io.Writer, name string, data any) error {
	return defaultTemplate.ExecuteTemplate(wr, name, data)
}

func ExecuteTemplateBytes(name string, data any) ([]byte, error) {
	var err error
	var buf bytes.Buffer
	err = defaultTemplate.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

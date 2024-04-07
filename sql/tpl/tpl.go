package tpl

import (
	"bytes"
	"go/format"
	"io"
	"os"
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
	globpath := filepath.Join(currentDir, "/*.tpl")
	defaultTemplate, err = template.ParseGlob(globpath)
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

func ExecuteTemplateWithCreateGoFile(genGoFilePath string, templateName string, paramsMap any) {
	err := createDirIfNotExists(genGoFilePath)
	if err != nil {
		panic(err)
	}

	typesGenFile, err := os.OpenFile(genGoFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	// log.Println(typesGenFile)

	typesGenData, err := ExecuteTemplateBytes(templateName, paramsMap)
	if err != nil {
		panic(err)
	}
	typesGenData, err = format.Source(typesGenData)
	if err != nil {
		panic(err)
	}
	typesGenFile.Write(typesGenData)
}

func createDirIfNotExists(genFilePath string) error {
	// 获取文件路径的目录部分
	dirPath := filepath.Dir(genFilePath)

	// 如果目录不存在,则创建
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

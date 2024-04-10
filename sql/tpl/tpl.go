package tpl

import (
	"bytes"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
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
	// log.Println(string(typesGenData))
	typesGenDataFormat, err := format.Source(typesGenData)
	if err != nil {
		// regexp.MustCompile(``)
		lerr := locateError(string(typesGenData), err.Error())
		panic(err.Error() + "\n" + lerr)
	}
	typesGenFile.Write(typesGenDataFormat)
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

func locateError(source, errMsg string) (errline string) {
	// 匹配错误行号和列号
	re := regexp.MustCompile(`(\d+)[:](\d+)[:]`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) < 3 {
		return ""
	}

	line, _ := strconv.Atoi(matches[1])
	// col, _ := strconv.Atoi(matches[2])

	// 将源代码分割成行
	lines := strings.Split(source, "\n")
	if line > len(lines) {
		return ""
	}

	// 计算错误位置的偏移
	// offset := 0
	// for i := 0; i < line-1; i++ {
	// 	offset += len(lines[i]) + 1 // +1 for newline
	// }
	// offset += col

	// // 返回错误位置的开始和结束偏移
	// start := offset - 1
	// end := start + 1
	return lines[line]
}

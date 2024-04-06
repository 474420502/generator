package sql_generator

import (
	"bytes"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var caser = cases.Title(language.English)

func createFileWithTplIfNotExists(filename string, do func(f io.Writer) error) error {
	// 检测文件是否存在
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// 文件不存在，创建文件
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		var buf = bytes.NewBuffer(nil)
		err = do(buf)
		if err != nil {
			panic(err)
		}
		data, err := format.Source(buf.Bytes())
		if err != nil {
			_, err = file.Write(buf.Bytes())
		} else {
			_, err = file.Write(data)
		}

		if err != nil {
			panic(err)
		}

		log.Printf("%s 文件已创建并写入内容\n", filename)
	} else if err != nil {
		// 发生其他错误
		return err
	} else {
		// 文件已存在
		log.Printf("%s 文件已存在\n", filename)
	}

	return nil
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

func toPascalCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		words[i] = caser.String(strings.ToLower(word))
	}
	return strings.Join(words, "")
}

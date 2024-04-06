package yamlconfig

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

func Parse2ConfigFile(filePath string, packageName string, yamlFile []byte) {
	// 读取YAML文件内容

	// 解析YAML内容
	var data interface{}
	err := yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// 创建输出文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	// // 写入包名和结构体定义
	file.WriteString(fmt.Sprintf("package %s\n\n", packageName))

	var buf bytes.Buffer
	buf.WriteString("type Config ")
	// 递归处理YAML数据
	processYAML(&buf, data, "", false)
	file.WriteString(buf.String())
}
func processYAML(buf io.Writer, data interface{}, indent string, isArray bool) {

	if isArray {
		buf.Write([]byte("[]"))
	}

	switch dv := data.(type) {
	case map[interface{}]interface{}:

		buf.Write([]byte("struct {\n"))

		nextIndent := indent + "    "
		for k, v := range dv {
			kstr := fmt.Sprintf("%v", k)
			keystr := formatFieldName(kstr)
			yamlTag := kstr
			buf.Write([]byte(nextIndent + keystr + " "))
			processYAML(buf, v, nextIndent, false)
			buf.Write([]byte(fmt.Sprintf("`yaml:\"%s\"`\n", yamlTag)))
		}
		buf.Write([]byte(indent + "} "))
	case []interface{}:
		for _, v := range dv {
			processYAML(buf, v, indent, true)
			break
		}
	case string:
		buf.Write([]byte("string "))
	case int:
		buf.Write([]byte("int64 "))
	case int64:
		buf.Write([]byte("int64 "))
	case float64:
		buf.Write([]byte("float64 "))
	case bool:
		buf.Write([]byte("bool "))
	default:
		buf.Write([]byte("interface{} "))
	}
}

// 格式化字段名
func formatFieldName(name string) string {

	// 将字段名中的下划线转换为驼峰命名
	name = strings.ReplaceAll(name, "_", " ")
	words := strings.Fields(name)
	var result []string
	for _, word := range words {

		word = cases.Upper(language.English).String(word[0:1]) + word[1:]
		result = append(result, word)
	}
	return strings.Join(result, "")
}

package jsonconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ConfigSetting struct {
	StructName  string
	FilePath    string
	PackageName string
	IsLog       bool
	JsonData    []byte
}

func New() *ConfigSetting {
	cs := &ConfigSetting{
		StructName:  "Config",
		FilePath:    "./config_gen.go",
		PackageName: "config",
		IsLog:       true,
	}
	return cs
}

func (cs *ConfigSetting) WithStructName(StructName string) *ConfigSetting {
	cs.StructName = StructName
	return cs
}

func (cs *ConfigSetting) WithFilePath(FilePath string) *ConfigSetting {
	cs.FilePath = FilePath
	return cs
}

func (cs *ConfigSetting) WithPackageName(PackageName string) *ConfigSetting {
	cs.PackageName = PackageName
	return cs
}

func (cs *ConfigSetting) WithJsonPath(JsonFile string) *ConfigSetting {
	jsonData, err := os.ReadFile(JsonFile)
	if err != nil {
		panic(err)
	}
	cs.JsonData = jsonData
	return cs
}

func (cs *ConfigSetting) WithJsonBytes(JsonBytes []byte) *ConfigSetting {
	cs.JsonData = JsonBytes
	return cs
}

func (cs *ConfigSetting) WithJsonString(JsonStr string) *ConfigSetting {
	cs.JsonData = []byte(JsonStr)
	return cs
}

func (cs *ConfigSetting) WithIsLog(IsLog bool) *ConfigSetting {
	cs.IsLog = IsLog
	return cs
}

func (cs *ConfigSetting) Create() {
	if len(cs.JsonData) == 0 {
		panic("WithJsonPath | WithJsonBytes | WithJsonString must be called at least once")
	}
	cs.jsonBytesParse()
}

func (cs *ConfigSetting) jsonBytesParse() {
	// Read Json data
	var data interface{}
	err := json.Unmarshal(cs.JsonData, &data)
	if err != nil {
		log.Fatalf("Failed to unmarshal Json: %v", err)
	}

	// Create output file
	file, err := os.OpenFile(cs.FilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	// Write package name and struct definition
	file.WriteString(fmt.Sprintf("package %s\n\n", cs.PackageName))
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("type %s ", cs.StructName))

	// Process Json data recursively
	processJson(&buf, data, "", false)
	file.WriteString(buf.String())
	log.Printf("package %s -> %s created!\n %s", cs.PackageName, cs.FilePath, buf.String())
}

func processJson(buf io.Writer, data interface{}, indent string, isArray bool) {
	if isArray {
		buf.Write([]byte("[]"))
	}
	switch dv := data.(type) {
	case map[string]interface{}:
		buf.Write([]byte("struct {\n"))
		nextIndent := indent + " "
		for k, v := range dv {
			keystr := formatFieldName(k)
			buf.Write([]byte(nextIndent + keystr + "    "))
			processJson(buf, v, nextIndent, false)
			buf.Write([]byte(fmt.Sprintf("`json:\"%s\"`\n", k)))
		}
		buf.Write([]byte(indent + "} "))
	case []interface{}:
		for _, v := range dv {
			processJson(buf, v, indent, true)
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

// Format field name
func formatFieldName(name string) string {
	name = strings.ReplaceAll(name, "_", " ")
	words := strings.Fields(name)
	var result []string
	for _, word := range words {
		word = cases.Upper(language.English).String(word[:1]) + word[1:]
		result = append(result, word)
	}
	return strings.Join(result, "")
}

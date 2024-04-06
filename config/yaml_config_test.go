package yamlconfig

import (
	"os"
	"testing"
)

func TestCase1(t *testing.T) {
	ydata, err := os.ReadFile("./etc.yaml")
	if err != nil {
		panic(err)
	}
	Parse2ConfigFile("./etc_config.go", "yamlconfig", ydata)
}

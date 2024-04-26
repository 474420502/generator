package yamlconfig

import (
	"testing"
)

func TestCase1(t *testing.T) {
	cs := New()
	cs.WithYamlPath("./etc.yaml").WithPackageName("yamlconfig").Create()
}

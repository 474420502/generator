package jsonconfig

import "testing"

func TestCase1(t *testing.T) {
	cs := New()

	cs.WithJsonPath("./etc.json").WithPackageName("jsonconfig").Create()
}

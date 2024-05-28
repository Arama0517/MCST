package lib_test

import (
	"encoding/json"
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/pkg/lib"
)

func TestDetectJava(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过寻找Java测试")
	}
	java, err := lib.DetectJava()
	if err != nil {
		t.Fatal(err)
	}
	jsonJava, err := json.MarshalIndent(java, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonJava))
}

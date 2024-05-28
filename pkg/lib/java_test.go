package lib_test

import (
	"encoding/json"
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/pkg/lib"
)

func TestDetectJava(t *testing.T) {
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

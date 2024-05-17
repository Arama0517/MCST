package apis_test

import (
	"encoding/json"
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/apis"
)

func TestGetFastMirrorDatas(t *testing.T) {
	data, err := apis.GetFastMirrorDatas()
	if err != nil {
		t.Error(err)
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(jsonData))
}

func TestGetFastMirrorBuildsDatas(t *testing.T) {
	data, err := apis.GetFastMirrorBuildsDatas("Mohist", "1.20.1")
	if err != nil {
		t.Error(err)
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(jsonData))
}

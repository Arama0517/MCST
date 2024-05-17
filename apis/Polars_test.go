package apis_test

import (
	"encoding/json"
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/apis"
)

func TestGetPolarsDatas(t *testing.T) {
	data, err := apis.GetPolarsDatas()
	if err != nil {
		t.Error(err)
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(jsonData))
}

func TestGetPolarsCoresDatas(t *testing.T) {
	data, err := apis.GetPolarsCoresDatas(2)
	if err != nil {
		t.Error(err)
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(jsonData))
}

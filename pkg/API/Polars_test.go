package api_test

import (
	"encoding/json"
	"testing"

	api "github.com/Arama-Vanarana/MCSCS-Go/pkg/API"
	"github.com/Arama-Vanarana/MCSCS-Go/pkg/lib"
)

func init() {
	lib.Init()
	api.InitPolars()
}

func TestPolars(t *testing.T) {
	jsonData, err := json.MarshalIndent(api.Polars, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonData))
}

func TestPolarsCore(t *testing.T) {
	data, err := api.GetPolarsCoresDatas(16)
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(jsonData))
}

func TestPolarsCoreDownload(t *testing.T) {
	data, err := api.GetPolarsCoresDatas(16)
	if err != nil {
		t.Fatal(err)
	}
	var info api.PolarsCores
	for _, v := range data {
		info = v
		break
	}
	path, err := api.DownloadPolarsServer(info.DownloadURL, info.Name)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
}

package apis_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/apis"
)

func TestGetPolarsDatas(t *testing.T) {
	data, err := apis.GetPolarsDatas()
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(jsonData))
}

func TestGetPolarsCoresDatas(t *testing.T) {
	data, err := apis.GetPolarsCoresDatas(2)
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(jsonData))
}

func TestDownloadPolarsServer(t *testing.T) {
	data, err := apis.GetPolarsCoresDatas(2)
	if err != nil {
		t.Fatal(err)
	}
	var name string
	var downloadUrl string
	for _, value := range data {
		name = value.Name
		downloadUrl = value.DownloadURL
		break
	}
	path, err := apis.DownloadPolarsServer(downloadUrl, name)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(path)
	err = os.Remove(path)
	if err != nil {
		t.Fatal(err)
	}
}

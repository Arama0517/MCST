package apis_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/apis"
)

func TestGetFastMirrorDatas(t *testing.T) {
	data, err := apis.GetFastMirrorDatas()
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(jsonData))
}

func TestGetFastMirrorBuildsDatas(t *testing.T) {
	data, err := apis.GetFastMirrorBuildsDatas("Mohist", "1.20.1")
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(jsonData))
}

func TestDownloadFastMirrorServer(t *testing.T) {
	data, err := apis.GetFastMirrorDatas()
	if err != nil {
		t.Fatal(err)
	}
	version := data["Mohist"].MC_Versions[0]
	buildsData, err := apis.GetFastMirrorBuildsDatas("Mohist", version)
	if err != nil {
		t.Fatal(err)
	}
	var build string
	for _, v := range buildsData {
		build = v.Core_Version
		break
	}
	path, err := apis.DownloadFastMirrorServer("Mohist", version, build)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(path)
	err = os.Remove(path)
	if err != nil {
		t.Fatal(err)
	}
}

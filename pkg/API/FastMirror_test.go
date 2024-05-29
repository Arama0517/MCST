package api_test

import (
	"encoding/json"
	"testing"

	api "github.com/Arama-Vanarana/MCSCS-Go/pkg/API"
	"github.com/Arama-Vanarana/MCSCS-Go/pkg/lib"
)

func init() {
	lib.Init()
	api.InitFastMirror()
}

func TestFastMirror(t *testing.T) {
	jsonData, err := json.MarshalIndent(api.FastMirror, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonData))
}

func TestFastMirrorBuilds(t *testing.T) {
	MC_Version := api.FastMirror["Mohist"].MC_Versions[0]
	builds, err := api.GetFastMirrorBuildsDatas("Mohist", MC_Version)
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.MarshalIndent(builds, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonData))
}

func TestFastMirrorBuildsDownload(t *testing.T) {
	MC_Version := api.FastMirror["Mohist"].MC_Versions[0]
	builds, err := api.GetFastMirrorBuildsDatas("Mohist", MC_Version)
	if err != nil {
		t.Fatal(err)
	}
	var build string
	for k := range builds {
		build = k
		break
	}
	path, err := api.DownloadFastMirrorServer(lib.ServerInfo{
		ServerType:       "Mohist",
		MinecraftVersion: MC_Version,
		BuildVersion:     build,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
}

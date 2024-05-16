package apis

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/Arama-Vanarana/MCSCS-Go/lib"
)

type FastMirrorDatas struct {
	data []struct {
		name        string
		tag         string
		homepage    string
		recommanded bool
		mc_versions []string
	}
	code    string
	suceess bool
	message string
}

func GetFastMirrorDatas() (FastMirrorDatas, error) {
	url := url.URL{
		Scheme: "https",
		Host:   "download.fastmirror.net",
		Path:   "/api/v3",
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return FastMirrorDatas{}, err
	}
	req.Header.Set("User-Agent", "MCSCS-Golang/"+lib.VERSION)
	resp, err := client.Do(req)
	if err != nil {
		return FastMirrorDatas{}, err
	}
	defer resp.Body.Close()
	var data FastMirrorDatas
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return FastMirrorDatas{}, err
	}
	return data, nil
}

type FastMirrorBuildsDatas struct {
	data struct {
		builds []struct {
			name         string
			mc_version   string
			core_version string
			update_time  string
			sha1         string
		}
		offset int
		limit  int
		count  int
	}
	code    string
	suceess bool
	message string
}

func GetFastMirrorBuildsDatas(server_type string, minecraft_version string) (FastMirrorBuildsDatas, error) {
	url := url.URL{
		Scheme:   "https",
		Host:     "download.fastmirror.net",
		Path:     "/api/v3/" + server_type + "/" + minecraft_version,
		RawQuery: "?offset=0&limit=25",
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return FastMirrorBuildsDatas{}, err
	}
	req.Header.Set("User-Agent", "MCSCS-Golang/"+lib.VERSION)
	resp, err := client.Do(req)
	if err != nil {
		return FastMirrorBuildsDatas{}, err
	}
	defer resp.Body.Close()
	var data FastMirrorBuildsDatas
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return FastMirrorBuildsDatas{}, err
	}
	return data, nil
}

func DownloadFastMirrorServers(server_type string, minecraft_version string, build_version string) (string, error) {
	url := url.URL{
		Scheme: "https",
		Host:   "download.fastmirror.net",
		Path:   "/download/" + server_type + "/" + minecraft_version + "/" + build_version,
	}
	path, err := lib.Download(url.String(), server_type+"-"+minecraft_version+"-"+build_version+".jar")
	if err != nil {
		return "", err
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hasher := sha1.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	hash := hasher.Sum(nil)
	FastMirrorBuildsData, err := GetFastMirrorBuildsDatas(server_type, minecraft_version)
	if err != nil {
		return "", err
	}
	if string(hash) != FastMirrorBuildsData.data.builds[0].sha1 {
		return "", errors.New("Sha1不匹配")
	}
	return path, nil
}

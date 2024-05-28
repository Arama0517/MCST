/*
 * MCSCS can be used to easily create, launch, and configure a Minecraft server.
 * Copyright (C) 2024 Arama
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package apis_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/pkg/apis"
	"github.com/Arama-Vanarana/MCSCS-Go/pkg/lib"
)

func init() {
	lib.Init()
	apis.InitPolars()
}

func TestGetPolarsDatas(t *testing.T) {
	jsonData, err := json.MarshalIndent(apis.Polars, "", "    ")
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
	if testing.Short() {
		t.Skip("跳过下载测试")
	}
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

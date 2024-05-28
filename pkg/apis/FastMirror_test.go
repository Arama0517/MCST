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
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/pkg/apis"
	"github.com/Arama-Vanarana/MCSCS-Go/pkg/lib"
)

func init() {
	lib.Init()
	apis.InitFastMirror()
}

func TestGetFastMirrorDatas(t *testing.T) {
	jsonData, err := json.MarshalIndent(apis.FastMirror, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(jsonData))
}

func TestGetFastMirrorBuildsDatas(t *testing.T) {
	data, err := apis.GetFastMirrorBuildsDatas("Mohist", apis.FastMirror["Mohist"].MC_Versions[0])
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(jsonData))
}

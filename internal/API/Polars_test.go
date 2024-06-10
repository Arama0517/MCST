/*
 * Minecraft Server Tool(MCST) is a command-line utility making Minecraft server creation quick and easy for beginners.
 * Copyright (c) 2024-2024 Arama.
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

package api_test

import (
	"encoding/json"
	"testing"

	"github.com/Arama-Vanarana/MCServerTool/internal/API"
	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	goversion "github.com/caarlos0/go-version"
)

func TestPolars(t *testing.T) {
	if err := lib.Init(goversion.GetVersionInfo()); err != nil {
		t.Fatal(err)
	}
	data, err := api.GetPolarsData()
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonData))
}

func TestPolarsCore(t *testing.T) {
	if err := lib.Init(goversion.GetVersionInfo()); err != nil {
		t.Fatal(err)
	}
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

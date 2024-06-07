/*
 * Minecraft Server Tool(MCST) is a command-libe utility making Minecraft server creation quick and easy for beginners.
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

package lib_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/Arama-Vanarana/MCServerTool/internal/lib"
)

var URL = url.URL{ // https://ash-speed.hetzner.com/100MB.bin
	Scheme: "https",
	Host:   "ash-speed.hetzner.com",
	Path:   "/100MB.bin",
}

func TestDownload(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过下载")
	}
	if err := lib.InitAll("dev"); err != nil {
		t.Fatal(err)
	}
	lib.EnableAria2c = false
	path, err := lib.NewDownloader(URL).Download()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
}

func TestAria2Downlaod(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过下载")
	}
	if err := lib.InitAll("dev"); err != nil {
		t.Fatal(err)
	}
	lib.EnableAria2c = true
	path, err := lib.NewDownloader(URL).Download()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
}

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

package requests_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/Arama0517/MCST/internal/configs"
	"github.com/Arama0517/MCST/internal/requests"
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
	if err := configs.InitData(); err != nil {
		t.Fatal(err)
	}
	configs.Configs.Settings.Aria2.Enabled = false
	path, err := requests.NewDownloader(URL).Download()
	if err != nil {
		panic(err)
	}
	t.Log(path)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
}

func TestAria2Download(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过下载")
	}
	if err := configs.InitData(); err != nil {
		t.Fatal(err)
	}
	configs.Configs.Settings.Aria2.Enabled = true
	path, err := requests.NewDownloader(URL).Download()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
}

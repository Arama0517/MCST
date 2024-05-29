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

package lib_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/pkg/lib"
)

func init() {
	lib.Init()
}

func TestMultiDownload(t *testing.T) {
	path, err := lib.Downloader{
		URL: url.URL{
			Scheme: "https",
			Host:   "golang.org",
			Path:   "/dl/go1.22.3.src.tar.gz",
		},
		FileName: "go1.22.3.src.tar.gz",
		Multi:    2,
	}.Download()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
	os.Remove(path)
}

func TestSingleDownload(t *testing.T) {
	path, err := lib.Downloader{
		URL: url.URL{
			Scheme: "https",
			Host:   "golang.org",
			Path:   "/dl/go1.22.3.src.tar.gz",
		},
		FileName: "go1.22.3.src.tar.gz",
		Multi:    1,
	}.Download()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
	os.Remove(path)
}

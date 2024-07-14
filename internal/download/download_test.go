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

package download_test

import (
	"os"
	"testing"

	"github.com/Arama0517/MCST/internal/download"
)

var URL = "https://dlied4.myapp.com/myapp/1104466820/cos.release-40109/10040714_com.tencent.tmgp.sgame_a2480356_8.2.1.9_F0BvnI.apk"

func TestDownload(t *testing.T) {
	path, err := download.NewDownloader(URL).Download()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
	if err = os.Remove(path); err != nil {
		t.Fatal(err)
	}
}

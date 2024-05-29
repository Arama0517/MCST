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
	"encoding/json"
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/pkg/lib"
)

func TestDetectJava(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过寻找Java测试")
	}
	java, err := lib.DetectJava()
	if err != nil {
		t.Fatal(err)
	}
	jsonJava, err := json.MarshalIndent(java, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonJava))
}

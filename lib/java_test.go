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
	"fmt"
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/lib"
)

func TestDetectJava(t *testing.T) {
	java := lib.DetectJava()
	jsonJava, err := json.MarshalIndent(java, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(jsonJava))
}

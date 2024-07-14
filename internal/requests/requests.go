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

package requests

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Arama0517/MCST/internal/build"
)

// NewRequest 替代 [http.NewRequest]; 此函数的作用是在 [http.NewRequest] 函数的基础上默认添加 User-Agent
func NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body) //nolint:forbidigo
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("MCST/%s", build.Version.GitVersion))
	return req, nil
}

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
	"net/url"

	"github.com/Arama0517/MCST/internal/build"
)

// Request 请求URL, 返回响应; 运行成功后请添加`defer resp.Body.Close()`到你的代码内
func Request(url url.URL, method string, header map[string]string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", fmt.Sprintf("MCST/%s ", build.Version.GitVersion))
	for key, value := range header {
		request.Header.Set(key, value)
	}
	return http.DefaultClient.Do(request)
}

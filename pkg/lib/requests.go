/*
 * Minecraft Server Tool(MCST) is a command-line utility making Minecraft server creation quick and easy for beginners.
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

package lib

import (
    "fmt"
    "net/http"
    "net/url"
)

// Request 请求URL, 返回响应; 运行成功后请添加`defer resp.Body.Close()`到你的代码内
func Request(URL url.URL, Method string, Header map[string]string) (*http.Response, error) {
    client := http.Client{}
    req, err := http.NewRequest(Method, URL.String(), nil)
    req.Header.Set("User-Agent", fmt.Sprintf("MCSCS-Go/%s", Version))
    for k, v := range Header {
        req.Header.Set(k, v)
    }
    if err != nil {
        return nil, err
    }
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

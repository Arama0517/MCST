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

// 极星云镜像

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arama0517/MCST/internal/requests"
)

func GetPolarsData() (map[string]PolarsData, error) {
	req, err := requests.NewRequest(http.MethodGet, "https://mirror.polars.cc/api/query/minecraft/core", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var data []PolarsData
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	result := map[string]PolarsData{}
	for i := 0; i < len(data); i++ {
		data := data[i]
		result[data.Name] = data
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}
	return result, nil
}

func GetPolarsCoresData(id int) (map[int]PolarsCores, error) {
	req, err := requests.NewRequest(http.MethodGet, fmt.Sprintf("https://mirror.polars.cc/api/query/minecraft/core/%d", id), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var data []PolarsCores
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	parsedData := map[int]PolarsCores{}
	for i := 0; i < len(data); i++ {
		data := data[i]
		parsedData[data.ID] = data
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}
	return parsedData, nil
}

type PolarsData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type PolarsCores struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
	Type        int    `json:"type"`
}

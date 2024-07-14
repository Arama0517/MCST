//go:build windows

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

package download

import (
	"time"

	"github.com/Arama0517/MCST/internal/configs"
	"github.com/Navid2zp/idm"
)

func (d *Downloader) idmDownload() (string, error) {
	downloader, err := idm.NewDownload(d.URL)
	if err != nil {
		return "", err
	}
	downloader.SetFilePath(configs.DownloadsDir)
	if err = downloader.Start(); err != nil {
		return "", err
	}
	if err = downloader.VerifyDownload(time.Second * 10); err != nil {
		return "", err
	}
	return downloader.GetFullPath()
}

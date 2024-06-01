#!/bin/sh
#
# Minecraft Server Tool(MCST) is a command-line utility making Minecraft server creation quick and easy for beginners.
# Copyright (c) 2024-2024 Arama.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
#

set -e

# 设置Github仓库信息
OWNER="Arama-Vanarana"
REPO="MCServerTool"
GITHUB_API="https://api.github.com/repos/${OWNER}/${REPO}/releases/latest"
GITHUB_DOWNLOAD="https://github.com/${OWNER}/${REPO}/releases/download"

get_latest_version() {
  TAG=$(curl -H "Accept: application/vnd.github.v3+json" -s $GITHUB_API | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
  if [ -z "$TAG" ]; then
    echo "Failed to fetch the latest version."
    exit 1
  fi
}

get_latest_version
curl "$GITHUB_DOWNLOAD/$TAG/MCST-Linux.tar.gz"
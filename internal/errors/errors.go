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

package errors

import (
	"errors"
)

var (
	ErrServerNotFound = errors.New("服务器不存在, 你可以使用 'MCST list' 命令查看所有已创建的服务器")
	ErrServerExists   = errors.New("服务器已存在, 你可以使用 'MCST list' 命令查看所有已创建的服务器")
)

var (
	ErrXmsToLow                 = errors.New("java初始堆内存不得小于1048576B(1MiB)")
	ErrXmxTooLow                = errors.New("java最大堆内存不得小于1048576(1MiB)")
	ErrXmxLessThanXms           = errors.New("java最大堆内存不得小于java初始堆内存")
	ErrXmxExceedsPhysicalMemory = errors.New("java最大堆内存不得大于物理内存")
	ErrXmsExceedsPhysicalMemory = errors.New("java初始堆内存不得大于物理内存")
)

var ErrInvalidUnit = errors.New("这不是一个有效的单位")

var ErrEulaRequired = errors.New("微软要求必须同意EULA协议(https://aka.ms/MinecraftEULA/)")

var ErrCoreNotFound = errors.New("核心不存在")

const (
	InitConfigFail = iota + 1
	InitLocaleFail
	RunFail
)

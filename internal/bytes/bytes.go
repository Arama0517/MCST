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

package bytes

import (
	"strconv"
	"strings"

	MCSTErrors "github.com/Arama0517/MCST/internal/errors"
)

const (
	Bytes uint64 = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
)

func ToBytes(byteStr string) (uint64, error) {
	var num, unit string

	// 分离数字部分和单位部分
	for _, char := range byteStr {
		switch {
		case char >= '0' && char <= '9':
			num += string(char)
		case (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z'):
			unit += strings.ToUpper(string(char))
		default:
			return 0, MCSTErrors.ErrInvalidUnit
		}
	}

	parsedNum, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return 0, err
	}
	switch {
	case strings.Contains(unit, "T"):
		return parsedNum * TiB, nil
	case strings.Contains(unit, "G"):
		return parsedNum * GiB, nil
	case strings.Contains(unit, "M"):
		return parsedNum * MiB, nil
	case strings.Contains(unit, "K"):
		return parsedNum * KiB, nil
	case strings.Contains(unit, "B"), unit == "":
		return parsedNum * Bytes, nil
	default:
		return 0, MCSTErrors.ErrInvalidUnit
	}
}

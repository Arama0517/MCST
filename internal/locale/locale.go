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

package locale

import (
	"embed"
	"encoding/json"

	"github.com/Arama0517/MCST/internal/configs"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	Bundle    *i18n.Bundle
	Localizer *i18n.Localizer
)

//go:embed locale.*.json
var localeFS embed.FS

func InitLocale() {
	Bundle = i18n.NewBundle(configs.Configs.Language)
	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, _ = Bundle.LoadMessageFileFS(localeFS, "locale.en.json")
	_, _ = Bundle.LoadMessageFileFS(localeFS, "locale.zh.json")
	Localizer = i18n.NewLocalizer(Bundle)
}

func GetLocaleMessage(messageID string) string {
	str, err := Localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		panic(err)
	}
	return str
}

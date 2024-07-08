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

	"github.com/Arama0517/MCST/internal/configs"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var (
	Bundle    *i18n.Bundle
	Localizer *i18n.Localizer
)

//go:embed locale.*.yaml
var localeFS embed.FS

func InitLocale() error {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	files, err := localeFS.ReadDir(".")
	if err != nil {
		return err
	}
	for _, f := range files {
		if _, err = Bundle.LoadMessageFileFS(localeFS, f.Name()); err != nil {
			return err
		}
	}
	Localizer = i18n.NewLocalizer(Bundle, configs.Configs.Settings.Language)
	return nil
}

func GetLocaleMessage(messageID string) string {
	msg, err := Localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
	if err != nil {
		panic(err)
	}
	return msg
}

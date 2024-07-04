//go:build unix

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

package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list         list.Model
	choice       string
	err          error
	needFileName string
	currentDir   string
	files        map[string]os.FileInfo
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			if i, ok := m.list.SelectedItem().(item); ok {
				if file, ok := m.files[i.title]; ok && file.IsDir() || i.title == ".." {
					m.currentDir = filepath.Join(m.currentDir, i.title)
				}
				if i.title == m.needFileName {
					m.choice = filepath.Join(m.currentDir, i.title)
					return m, tea.Quit
				}
				items, err := m.UpdateFiles()
				if err != nil {
					return m, func() tea.Msg { return err }
				}
				m.list.SetItems(items)
			}
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case error:
		m.err = msg
		return m, tea.Quit
	}

	m.list.Title = fmt.Sprintf("当前路径: `%s`\n请选择名为 `%s` 的文件", m.currentDir, m.needFileName)

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) UpdateFiles() ([]list.Item, error) {
	var items []list.Item
	files, err := os.ReadDir(m.currentDir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			return nil, err
		}
		m.files[fileInfo.Name()] = fileInfo

		var fileType string
		switch mode := fileInfo.Mode(); {
		case mode.IsRegular():
			fileType = "普通文件"
		case mode.IsDir():
			fileType = "目录"
		case mode&os.ModeSymlink != 0:
			fileType = "符号链接"
			linkPath := filepath.Join(m.currentDir, fileInfo.Name())
			linkTarget, err := os.Readlink(linkPath)
			if err == nil {
				if !filepath.IsAbs(linkTarget) {
					linkTarget = filepath.Join(filepath.Dir(linkPath), linkTarget)
				}
				linkTargetInfo, err := os.Stat(linkTarget)
				if err == nil {
					fileType += fmt.Sprintf(" -> %s", linkTarget)
					m.files[fileInfo.Name()] = linkTargetInfo
				}
			}
		case mode&os.ModeNamedPipe != 0:
			fileType = "命名管道"
		case mode&os.ModeSocket != 0:
			fileType = "套接字"
		case mode&os.ModeDevice != 0:
			fileType = "设备文件"
		default:
			fileType = "未知(无法识别)"
		}
		items = append(items, item{file.Name(), fileType})
	}
	if m.currentDir != "/" {
		items = append(items, item{"..", "返回上一层目录"})
	}
	return items, nil
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func Run(fileName string) (string, error) {
	userHomeDir, _ := os.UserHomeDir()
	listModel := model{
		needFileName: fileName,
		files:        make(map[string]os.FileInfo),
		currentDir:   userHomeDir,
	}

	items, err := listModel.UpdateFiles()
	if err != nil {
		return "", err
	}
	listModel.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m, err := tea.NewProgram(listModel, tea.WithAltScreen()).Run()
	return m.(model).choice, err
}

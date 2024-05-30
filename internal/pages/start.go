package pages

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/urfave/cli/v2"
)

func Start(c *cli.Context) error {
	configs, err := lib.LoadConfigs()
	if err != nil {
		return err
	}
	server, exists := configs.Servers[c.String("server")]
	if !exists {
		return errors.New("服务器不存在")
	}
	cmd := exec.Command(server.Java.Path)
	cmd.Args = append(cmd.Args, fmt.Sprintf("-Xmx%d", server.Java.Xmx), fmt.Sprintf("-Xms%d", server.Java.Xms), "-Dfile.encoding=UTF-8")
	cmd.Args = append(cmd.Args, server.Java.Args...)
	cmd.Args = append(cmd.Args, "-jar", "server.jar")
	cmd.Args = append(cmd.Args, server.ServerArgs...)
	cmd.Dir = filepath.Join(lib.ServersDir, server.Name)
	cmd.Stdout = c.App.Writer
	cmd.Stderr = c.App.ErrWriter
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

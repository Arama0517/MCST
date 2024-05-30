package pages

import (
	"fmt"

	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/urfave/cli/v2"
)

func ListServers(c *cli.Context) error {
	configs, err := lib.LoadConfigs()
	if err != nil {
		return err
	}
	for _, config := range configs.Servers {
		fmt.Println(config.Name)
	}
	return nil
}

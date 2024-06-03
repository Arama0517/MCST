package pages

import (
	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/urfave/cli/v2"
)

var Delete = cli.Command{
	Name: "删除服务器",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "server",
			Aliases:  []string{"s"},
			Usage:    "要删除的服务器名称",
			Required: true,
		},
	},
	Action: func(context *cli.Context) error {
		server := context.String("server")
		configs, err := lib.LoadConfigs()
		if err != nil {
			return err
		}
		switch c, err := confirm("确认删除此服务器?", context); {
		case err != nil:
			return err
		case c:
			delete(configs.Servers, server)
			err = configs.Save()
			if err != nil {
				return err
			}
		case !c:
			break
		}
		return nil
	},
}

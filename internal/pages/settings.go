package pages

import (
	"fmt"

	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/urfave/cli/v2"
)

var Settings = cli.Command{
	Name:  "settings",
	Usage: "查看/修改程序设置",
	Subcommands: []*cli.Command{
		{
			Name:  "aria2c",
			Usage: "Aria2c 下载器的设置",
			Subcommands: []*cli.Command{
				{
					Name:  "path",
					Usage: "Aria2c 路径",
					Flags: []cli.Flag{
						&cli.PathFlag{
							Name:  "path",
							Usage: "Aria2c 路径, 如果为auto则自动寻找",
						},
					},
					Action: func(ctx *cli.Context) error {
						configs, err := lib.LoadConfigs()
						if err != nil {
							return err
						}
						path := ctx.Path("path")
						if path == "" {
							fmt.Println(configs.Aria2c.Path)
						} else {
							configs.Aria2c.Path = path
							err = configs.Save()
							if err != nil {
								return err
							}
						}
						return nil
					},
				},
				{
					Name:  "args",
					Usage: "Aria2c 参数",
					Flags: []cli.Flag{
						&cli.StringSliceFlag{
							Name:  "arg",
							Usage: "Aria2c 参数",
						},
					},
					Action: func(ctx *cli.Context) error {
						configs, err := lib.LoadConfigs()
						if err != nil {
							return err
						}
						args := ctx.StringSlice("arg")
						if len(args) == 0 {
							for _, arg := range configs.Aria2c.Args {
								fmt.Print(arg + " ")
							}
						} else {
							configs.Aria2c.Args = args
							err = configs.Save()
							if err != nil {
								return err
							}
						}
						return nil
					},
				},
			},
		},
		{
			Name:  "eula",
			Usage: "是否自动同意EULA协议",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "accept",
					Usage: "是否同意EULA协议",
				},
			},
			Action: func(ctx *cli.Context) error {
				configs, err := lib.LoadConfigs()
				if err != nil {
					return err
				}
				if ctx.IsSet("accept") {
					configs.AutoAcceptEULA = ctx.Bool("accept")
					err = configs.Save()
					if err != nil {
						return err
					}
				} else {
					fmt.Println(configs.AutoAcceptEULA)
				}
				return nil
			},
		},
	},
}

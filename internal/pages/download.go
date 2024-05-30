package pages

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"

	api "github.com/Arama-Vanarana/MCServerTool/pkg/API"
	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/urfave/cli/v2"
)

func ListCores(c *cli.Context) error {
	configs, err := lib.LoadConfigs()
	if err != nil {
		return err
	}
	for i, Core := range configs.Cores {
		fmt.Printf("%s(%d): %s\n", Core.FileName, i, Core.FilePath)
	}
	return nil
}

func Local(c *cli.Context) error {
	configs, err := lib.LoadConfigs()
	if err != nil {
		return err
	}
	configs.Cores = append(configs.Cores, lib.Core{
		FileName: c.Path("path"),
	})
	if err := configs.Save(); err != nil {
		return err
	}
	return nil
}

func Remote(c *cli.Context) error {
	configs, err := lib.LoadConfigs()
	if err != nil {
		return err
	}
	url, err := url.Parse(c.String("URL"))
	if err != nil {
		return err
	}
	path, err := (&lib.Downloader{
		URL: *url,
	}).Download()
	if err != nil {
		return err
	}
	configs.Cores = append(configs.Cores, lib.Core{
		URL:      *url,
		FileName: filepath.Base(path),
		FilePath: path,
	})
	if err := configs.Save(); err != nil {
		return err
	}
	return nil
}

func FastMirror(c *cli.Context) error {
	Core := c.String("Core")
	MinecraftVersion := c.String("MinecraftVersion")
	BuildVersion := c.String("BuildVersion")
	list := c.Bool("list")
	fastMirror, err := api.GetFastMirrorDatas()
	if err != nil {
		return err
	}
	if list {
		switch {
		case Core != "" && MinecraftVersion != "":
			fastMirrorBuilds, err := api.GetFastMirrorBuildsDatas(Core, MinecraftVersion)
			if err != nil {
				return err
			}
			for _, data := range fastMirrorBuilds {
				fmt.Printf("%s: 更新时间: %s, SHA1: %s\n", data.Core_Version, data.Update_Time, data.Sha1)
			}
		case Core != "":
			for _, data := range fastMirror[Core].MC_Versions {
				fmt.Println(data)
			}
		default:
			return errors.New("没有这个用法")
		}
	} else {
		if Core == "" || MinecraftVersion == "" || BuildVersion == "" {
			return errors.New("缺少必要参数")
		} else {
			path, err := api.DownloadFastMirrorServer(Core, MinecraftVersion, BuildVersion)
			if err != nil {
				return err
			}
			configs, err := lib.LoadConfigs()
			if err != nil {
				return err
			}
			configs.Cores = append(configs.Cores, lib.Core{
				FileName: filepath.Base(path),
				FilePath: path,
				ExtrasData: map[string]string{
					"core":          Core,
					"mc_version":    MinecraftVersion,
					"build_version": BuildVersion,
				},
			})
			if err := configs.Save(); err != nil {
				return err
			}
		}
	}
	return nil
}

func Polars(c *cli.Context) error {
	TypeID := c.Int("TypeID")
	CoreID := c.Int("CoreID")
	list := c.Bool("list")
	polars, err := api.GetPolarsData()
	if err != nil {
		return err
	}
	if list {
		switch {
		case TypeID == 0 && CoreID == 0:
			for _, data := range polars {
				fmt.Printf("%s(%d): %s\n", data.Name, data.ID, data.Description)
			}
		case TypeID != 0 && CoreID == 0:
			data, err := api.GetPolarsCoresDatas(TypeID)
			if err != nil {
				return err
			}
			for _, core := range data {
				fmt.Printf("%s(%d): %s\n", core.Name, core.ID, core.DownloadURL)
			}
		default:
			return errors.New("没有这个用法")
		}
	} else {
		if TypeID == 0 || CoreID == 0 {
			return errors.New("缺少必要参数")
		} else {
			data, err := api.GetPolarsCoresDatas(TypeID)
			if err != nil {
				return err
			}
			path, err := api.DownloadPolarsServer(data[CoreID].DownloadURL)
			if err != nil {
				return err
			}
			configs, err := lib.LoadConfigs()
			if err != nil {
				return err
			}
			url, err := url.Parse(data[CoreID].DownloadURL)
			if err != nil {
				return err
			}
			configs.Cores = append(configs.Cores, lib.Core{
				URL:      *url,
				FileName: filepath.Base(path),
				FilePath: path,
			})
			if err := configs.Save(); err != nil {
				return err
			}
		}
	}
	return nil
}

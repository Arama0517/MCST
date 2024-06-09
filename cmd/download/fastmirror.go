package download

import (
	"sort"
	"time"

	api "github.com/Arama-Vanarana/MCServerTool/internal/API"
	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

type fastMirrorCmdFlags struct {
	core             string
	minecraftVersion string
	buildVersion     string
}

func newFastMirrorCmd() *cobra.Command {
	flags := fastMirrorCmdFlags{}
	cmd := &cobra.Command{
		Use:     "fastmirror",
		Aliases: []string{"fm"},
		Short:   "从无极镜像下载核心",
		Long:    "无极镜像: <https://www.fastmirror.net/>\n请使用 'MCST download fastmirror list' 获取所需的参数",
		RunE: func(*cobra.Command, []string) error {
			downloader := api.GetFastMirrorDownloader(flags.core, flags.minecraftVersion, flags.buildVersion)
			path, err := downloader.Download()
			if err != nil {
				return err
			}
			configs, err := lib.LoadConfigs()
			if err != nil {
				return err
			}
			configs.Cores = append(configs.Cores, lib.Core{
				URL:      downloader.URL.String(),
				FileName: downloader.FileName,
				FilePath: path,
				ID:       len(configs.Cores) + 1,
			})
			return configs.Save()
		},
	}
	cmd.AddCommand(newListFastMirrorCmd())
	cmd.Flags().StringVarP(&flags.core, "core", "c", "", "核心的名称, 例如: \"Mohist\"")
	cmd.Flags().StringVarP(&flags.minecraftVersion, "mc_version", "m", "", "核心的Minecraft版本, 例如: \"1.20.1\"")
	cmd.Flags().StringVarP(&flags.buildVersion, "build_version", "b", "", "核心的构建版本, 例如: \"build738\"")
	_ = cmd.MarkFlagRequired("core")
	_ = cmd.MarkFlagRequired("mc_version")
	_ = cmd.MarkFlagRequired("build_version")
	return cmd
}

func newListFastMirrorCmd() *cobra.Command {
	flags := fastMirrorCmdFlags{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "获取无极镜像的核心信息",
		Long: `使用此命令时拥有三种情况:
1. 不使用任何参数: 输出所有核心的名称
2. 使用 '--core': 输出此核心的所有Minecraft版本
3. 使用 '--core' 和 '--mc_version': 返回此核心的所有构建版本`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fastMirror, err := api.GetFastMirrorDatas()
			if err != nil {
				return err
			}
			switch cmdFlags := cmd.Flags(); {
			case !cmdFlags.Changed("core") && !cmdFlags.Changed("mc_version"):
				keys := make([]string, 0, len(fastMirror))
				for k := range fastMirror {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, key := range keys {
					data := fastMirror[key]
					log.WithFields(log.Fields{
						"推荐": data.Recommanded,
						"标签": data.Tag,
						"主页": data.Homepage,
					}).Info(data.Name)
				}
			case cmdFlags.Changed("core") && !cmdFlags.Changed("mc_version"):
				data := fastMirror[flags.core].MinecraftVersions
				sort.Strings(data)
				for _, version := range data {
					log.Info(version)
				}
			case cmdFlags.Changed("core") && cmdFlags.Changed("mc_version"):
				fastMirror, err := api.GetFastMirrorBuildsDatas(flags.core, flags.minecraftVersion)
				if err != nil {
					return err
				}
				keys := make([]string, 0, len(fastMirror))
				for k := range fastMirror {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, key := range keys {
					data := fastMirror[key]
					updateTime, err := time.Parse("2006-01-02T15:04:05", data.UpdateTime)
					if err != nil {
						return err
					}
					log.WithFields(log.Fields{
						"SHA1": data.Sha1,
						"更新时间": updateTime.Format("2006-01-02 15:04:05"),
					}).Info(data.CoreVersion)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&flags.core, "core", "c", "", "核心的名称, 例如: \"Mohist\"")
	cmd.Flags().StringVarP(&flags.minecraftVersion, "mc_version", "m", "", "核心的Minecraft版本, 例如: \"1.20.1\"")
	return cmd
}

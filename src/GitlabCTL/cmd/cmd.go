package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/longyuan/gitlab.v3/console"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

func Backup() []*cobra.Command {
	var backupCmd = &cobra.Command{
		Use:     "backup",
		Short:   "Backup Gitlab Repository",
		Example: "backup -H gitlab.com -t WXZXXX1",
		Run: func(cmd *cobra.Command, args []string) {
			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			var host, token string
			if configPath == "" {
				token, err = cmd.Flags().GetString("token")
				if err != nil {
					color.Red(fmt.Sprint(err))
					return
				}
				if token == "" {
					color.Red("Not Set token value ?")
					return
				}
				host, err = cmd.Flags().GetString("host")
				if err != nil {
					color.Red(fmt.Sprint(err))
					return
				}
				if host == "" {
					color.Red("Not Set host value ?")
					return
				}
			} else {
				// 读取Yaml 文件
				fileBytes, err := os.ReadFile(configPath)
				if err != nil {
					color.Red(fmt.Sprint(err))
					return
				}
				var gitlabConfig = struct {
					Host  string `yaml:"host"`
					Token string `yaml:"token"`
				}{}
				err = yaml.Unmarshal(fileBytes, &gitlabConfig)
				if err != nil {
					color.Red(fmt.Sprint(err))
					return
				}
				host = gitlabConfig.Host
				token = gitlabConfig.Token
			}
			output, err := cmd.Flags().GetString("output")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			_, err = console.Backup(host, token, output)
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
		},
	}
	backupCmd.Flags().StringP("host", "H", "", "Gitlab URL")
	backupCmd.Flags().StringP("token", "t", "", "Gitlab Token")
	backupCmd.Flags().StringP("config", "c", "", "Config")
	backupCmd.Flags().StringP("output", "o", "", "Output File")

	var cronBackupCmd = &cobra.Command{
		Use:     "cron-backup",
		Short:   "Backup Gitlab Repository",
		Example: "backup -H gitlab.com -t WXZXXX1",
		Run: func(cmd *cobra.Command, args []string) {
			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if configPath == "" {
				color.Red("ConfigPath Required")
				return
			}
			cron, err := cmd.Flags().GetString("cron")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if cron == "" {
				color.Red("Cron Required")
				return
			}
			cloudStorageConfig, err := cmd.Flags().GetString("cloud-storage")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if cloudStorageConfig == "" {
				color.Red("CloudStorageConfig Required")
				return
			}
			err = console.CronBackup(configPath, cron, cloudStorageConfig)
			if err != nil {
				return
			}
		},
	}
	cronBackupCmd.Flags().StringP("config", "c", "", "Config")
	cronBackupCmd.Flags().String("cron", "", "Cron")
	cronBackupCmd.Flags().String("cloud-storage", "", "CloudStorage Config")

	return []*cobra.Command{
		backupCmd,
		cronBackupCmd,
	}
}

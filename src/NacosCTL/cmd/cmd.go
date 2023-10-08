package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/longyuan/nacos.v3/console"
	"github.com/spf13/cobra"
)

func Backup() []*cobra.Command {
	var backupCmd = &cobra.Command{
		Use:     "backup",
		Short:   "Backup Nacos Config",
		Example: "backup -u nacos -p 1234 -h 127.0.0.1",
		Run: func(cmd *cobra.Command, args []string) {
			host, err := cmd.Flags().GetString("host")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if host == "" {
				color.Red("Host Required")
				return
			}
			username, err := cmd.Flags().GetString("username")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if username == "" {
				color.Red("Username Required")
				return
			}
			password, err := cmd.Flags().GetString("password")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if password == "" {
				color.Red("Password Required")
				return
			}
			output, err := cmd.Flags().GetString("output")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			_, err = console.Backup(host, username, password, output)
			if err != nil {
				return
			}
		},
	}
	backupCmd.Flags().StringP("username", "u", "", "Nacos Username")
	backupCmd.Flags().StringP("password", "p", "", "Nacos Password")
	backupCmd.Flags().StringP("host", "H", "", "Nacos Host")
	backupCmd.Flags().StringP("output", "o", "", "Nacos Output")

	var backupAliyunCmd = &cobra.Command{
		Use:     "ali-backup",
		Short:   "Backup Nacos Config",
		Example: "ali-backup -u nacos -p 1234 -h 127.0.0.1",
		Run: func(cmd *cobra.Command, args []string) {
			accessKeyId, err := cmd.Flags().GetString("accessKeyId")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if accessKeyId == "" {
				color.Red("AccessKeyId Required")
				return
			}
			accessKeySecret, err := cmd.Flags().GetString("accessKeySecret")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if accessKeySecret == "" {
				color.Red("AccessKeySecret Required")
				return
			}
			instanceId, err := cmd.Flags().GetString("instanceId")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if instanceId == "" {
				color.Red("InstanceId Required")
				return
			}
			namespace, err := cmd.Flags().GetString("namespace")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if namespace == "" {
				color.Red("Namespace Required")
				return
			}
			output, err := cmd.Flags().GetString("output")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			_, err = console.AliBackup(accessKeyId, accessKeySecret, instanceId, namespace, output)
			if err != nil {
				return
			}
		},
	}
	backupAliyunCmd.Flags().StringP("accessKeyId", "k", "", "Aliyun AccessKeyId")
	backupAliyunCmd.Flags().StringP("accessKeySecret", "s", "", "Aliyun AccessKeySecret")
	backupAliyunCmd.Flags().StringP("instanceId", "i", "", "Aliyun InstanceId")
	backupAliyunCmd.Flags().StringP("namespace", "n", "", "Aliyun InstanceId Namespace")
	backupAliyunCmd.Flags().StringP("output", "o", "", "Nacos Output")

	var cronBackupCmd = &cobra.Command{
		Use:     "cron-backup",
		Short:   "Backup Nacos Config",
		Example: "backup -u nacos -p 1234 -h 127.0.0.1",
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
		backupAliyunCmd,
		cronBackupCmd,
	}
}

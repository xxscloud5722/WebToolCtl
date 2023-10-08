package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/longyuan/kubernetes.v3/console"
	"github.com/spf13/cobra"
)

func Backup() []*cobra.Command {
	var backupCmd = &cobra.Command{
		Use:     "backup",
		Short:   "Backup Kubernetes Config",
		Example: "backup -c ./conf/kubeconfig-example.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if configPath == "" {
				color.Red("Not Set kubeconfig file ?")
				return
			}
			outputFile, err := cmd.Flags().GetString("output")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			color.Green(fmt.Sprintf("[Kubernetes] kubeconfig: %s , output: %s", configPath, outputFile))
			_, err = console.Backup(configPath, outputFile, nil)
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
		},
	}
	backupCmd.Flags().StringP("config", "c", "", "Config Path")
	backupCmd.Flags().StringP("output", "o", "", "Output Path")

	var cronBackupCmd = &cobra.Command{
		Use:     "cron-backup",
		Short:   "Backup Kubernetes Config",
		Example: "backup -c ./conf/kubeconfig-example.yaml -cron '0 0 * * ?' -cs URL,SecretId,SecretKey",
		Run: func(cmd *cobra.Command, args []string) {
			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if configPath == "" {
				color.Red("Not Set kubeconfig file ?")
				return
			}
			cron, err := cmd.Flags().GetString("cron")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if cron == "" {
				color.Red("Not Set Cron ?")
				return
			}
			cloudStorage, err := cmd.Flags().GetString("cloud-storage")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if cloudStorage == "" {
				color.Red("Not Set CloudStorage Config ?")
				return
			}
			err = console.CronBackup(configPath, cron, cloudStorage)
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
		},
	}
	cronBackupCmd.Flags().StringP("config", "c", "", "Config Path")
	cronBackupCmd.Flags().String("cron", "", "Cron")
	cronBackupCmd.Flags().String("cloud-storage", "", "CloudStorage Config")

	return []*cobra.Command{
		backupCmd,
		cronBackupCmd,
	}
}

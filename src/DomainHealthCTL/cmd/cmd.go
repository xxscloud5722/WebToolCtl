package cmd

import (
	"fmt"
	"github.com/fatih/color"
	console "github.com/longyuan/domain.v3/console"
	"github.com/spf13/cobra"
)

func Cmd() []*cobra.Command {
	var sslCmd = &cobra.Command{
		Use:     "ssl",
		Short:   "Host SSL Info",
		Example: "ssl www.baidu.com",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) <= 0 {
				return
			}
			err := console.SSL(args[0])
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
		},
	}

	var whoisCmd = &cobra.Command{
		Use:     "whois",
		Short:   "Host Whois Info",
		Example: "whois gitlab.com",
		Run: func(cmd *cobra.Command, args []string) {
			original, err := cmd.Flags().GetString("original")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if len(args) <= 0 {
				return
			}
			err = console.Whois(args[0], original != "false")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
		},
	}
	whoisCmd.Flags().StringP("original", "o", "false", "original information")

	var scanCmd = &cobra.Command{
		Use:     "scan",
		Short:   "Scan Config",
		Example: "scan ./domain.txt",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) <= 0 {
				return
			}
			err := console.Scan(args[0])
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
		},
	}

	var cronCmd = &cobra.Command{
		Use:     "cron",
		Short:   "Cron Job",
		Example: "cron -c domain.txt -n ",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := cmd.Flags().GetString("config")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if config == "" {
				color.Red("Not Set config value ?")
				return
			}
			cron, err := cmd.Flags().GetString("cron")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if cron == "" {
				color.Red("Not Set cron value ?")
				return
			}
			notice, err := cmd.Flags().GetString("notice")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if notice == "" {
				color.Red("Not Set notice value ?")
				return
			}
			err = console.CronJob(config, cron, notice)
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
		},
	}
	cronCmd.Flags().StringP("config", "c", "", "Domain Config")
	cronCmd.Flags().String("cron", "", "Cron")
	cronCmd.Flags().StringP("notice", "n", "", "Notice Config")

	return []*cobra.Command{
		sslCmd,
		whoisCmd,
		scanCmd,
		cronCmd,
	}
}

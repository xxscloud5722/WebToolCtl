package main

import (
	"github.com/longyuan/gitlab.v3/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "gctl"}
	for _, it := range cmd.Backup() {
		rootCmd.AddCommand(it)
	}
	err := rootCmd.Execute()
	if err != nil {
		return
	}
}

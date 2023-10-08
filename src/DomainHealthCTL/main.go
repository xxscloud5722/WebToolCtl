package main

import (
	"github.com/longyuan/domain.v3/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "domain"}
	for _, it := range cmd.Cmd() {
		rootCmd.AddCommand(it)
	}
	err := rootCmd.Execute()
	if err != nil {
		return
	}
}

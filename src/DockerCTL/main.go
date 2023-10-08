package main

import (
	"github.com/longyuan/docker.v3/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "dctl"}
	for _, it := range cmd.Pipeline() {
		rootCmd.AddCommand(it)
	}
	for _, it := range cmd.Console() {
		rootCmd.AddCommand(it)
	}
	err := rootCmd.Execute()
	if err != nil {
		return
	}
}

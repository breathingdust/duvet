package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "duvet",
	Short: "CLI tool to generate resource statistics for the Terraform AWS Provider",
	Long:  `CLI tool to generate resource statistics for the Terraform AWS Provider`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}

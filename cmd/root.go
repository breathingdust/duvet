package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	outputFormat string
	outputDir    string
	providerDir  string

	rootCmd = &cobra.Command{
		Use:   "duvet",
		Short: "CLI tool to generate resource statistics for the Terraform AWS Provider",
		Long:  `CLI tool to generate resource statistics for the Terraform AWS Provider`,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	serviceCmd.Flags().StringVarP(&outputFormat, "format", "f", ".", "Format of outputted files (markdown|html)")
	serviceCmd.Flags().StringVarP(&outputDir, "directory", "d", ".", "Directory to place outputted files")
	serviceCmd.Flags().StringVarP(&providerDir, "providerDir", "p", ".", "Path of local AWS provider source")
}

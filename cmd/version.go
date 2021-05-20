package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of go faas",
	Long:  `version of faas`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version 0.0.1.0 for faas")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

package cmd

import (
	"github.com/spf13/cobra"
	"go_faas/util"
)

var showPluginsCmd = &cobra.Command{
	Use:   "showPlugins",
	Short: "show all Plugins",
	Long:  `show all Plugins`,
	Run:   showPlugins,
}

func showPlugins(cmd *cobra.Command, args []string) {
	util.ShowAllPlugins()
}

func init() {
	rootCmd.AddCommand(showPluginsCmd)
}

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"go_faas/util"
)

var callPluginsCmd = &cobra.Command{
	Use:   "callPlugin",
	Short: "call Plugin func",
	Long:  `call Plugin func`,
	Run:   callPlugin,
}

var (
	pluginMethod string
	params       string
)

func callPlugin(cmd *cobra.Command, args []string) {
	var pluginParams []interface{}
	if params != "" {
		if err := json.Unmarshal([]byte(params), &pluginParams); err != nil {
			fmt.Println(err)
			panic("params为json类型")
		}
	}
	fmt.Println("Call PluginMethod: "+pluginMethod+" params: ", pluginParams)
	ret := util.DoInvokePlugin(pluginMethod, pluginParams...)
	if ret.GetCode() != nil {
		fmt.Println("Call PluginMethod: " + pluginMethod + " failed")
	} else {
		ret.ShowData()
	}

}

func init() {
	rootCmd.AddCommand(callPluginsCmd)
	callPluginsCmd.PersistentFlags().StringVar(&pluginMethod, "method", "", "plugin method")
	callPluginsCmd.PersistentFlags().StringVar(&params, "params", "", "plugin method params")
}

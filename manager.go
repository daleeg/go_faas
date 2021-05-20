package main

import (
	"fmt"
	"go_faas/cmd"

	// _ "go_restframework/util"
	"go_faas/util"
	// "plugin"
)

func main() {
	cmd.Execute()
	util.DoInvokePlugin(util.PluginItems, []interface{}{"world"})
	fmt.Println("Process Stop ========")
}

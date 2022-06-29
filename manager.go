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
	util.DoInvokePlugin("helloPlugin.hello.PrintNowTime", []interface{}{"world"})
	fmt.Println("Process Stop ========")
}

package main

import (
	"go_faas/cmd"
	// "plugin"
)

func main() {
	cmd.Start()
	//util.DoInvokePlugin("helloPlugin.hello.PluginPrintNowTime", "world")
	//util.DoInvokePlugin("rsaPlugin.rsa.PluginRSAGenKey", 4096, "privateKey.pem", "publicKey.pem")
	//fmt.Println("Process Stop ========")
}

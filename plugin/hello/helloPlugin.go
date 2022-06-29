package main

import (
	"fmt"
	"go_faas/plugin"
	"time"
)

var PluginBaseInfo = plugin.PluginBaseInfoNode{
	Name: "hello",
	Desc: "hello plugin",
	Function: plugin.PluginFunction{
		Name: "PrintNowTime", //可用函数名
		Params: []plugin.FuncParam{
			{
				Type: "string",
				Key:  "world",
			},
		},
	},
}

// func PluginHandle() {
// 	PrintNowTime()
// }

// PrintNowTime 打印当前时间
func PrintNowTime(world string) {
	nowSecond := time.Now().Second()
	if nowSecond%2 == 0 {
		fmt.Println("Hello,", world)
	} else {
		fmt.Println("Get out,", world)
	}
	fmt.Println(time.Now())
}

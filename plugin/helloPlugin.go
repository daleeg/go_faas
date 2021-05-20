package main

import (
	"fmt"
	"time"
)

var PluginBaseInfo = PluginBaseInfoNode{
	Name: "hello plugin",
	Desc: "hello plugin",
	Function: PluginFunction{
		Name: "PrintNowTime", //可用函数名
		Params: []FuncParam{
			{
				Type: "string",
				key:  "world",
			},
		},
	},
}

// func PluginHandle() {
// 	PrintNowTime()
// }

// 打印当前时间
func PrintNowTime(world string) {
	nowSecond := time.Now().Second()
	if nowSecond%2 == 0 {
		fmt.Println("Hello,", world)
	} else {
		fmt.Println("Get out,", world)
	}
	fmt.Println(time.Now())
}

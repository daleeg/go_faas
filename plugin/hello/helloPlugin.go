package main

import (
	"fmt"
	"time"
)

var (
	PackageName = "hello"
)
// PluginPrintNowTime 打印当前时间
func PluginPrintNowTime(world string) {
	nowSecond := time.Now().Second()
	if nowSecond%2 == 0 {
		fmt.Println("Hello,", world)
	} else {
		fmt.Println("Get out,", world)
	}
	fmt.Println(time.Now())
}

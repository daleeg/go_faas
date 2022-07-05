package util

import "reflect"


type PluginBaseInfoNode struct {
	Name     string         // 插件名称
	Desc     string         // 插件描述
	Function reflect.Value // 插件可用函数
	Params  []reflect.Type // 函数参数
	Returns []reflect.Type // 函数参数
}

type PluginItem struct {
	Name           string
	PluginBaseInfo *PluginBaseInfoNode
}



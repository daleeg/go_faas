package main

type PluginBaseInfoNode struct {
	Name     string         // 插件名称
	Desc     string         // 插件描述
	Function PluginFunction // 插件可用函数
}

type PluginFunction struct {
	Name   string
	Params []FuncParam // 函数参数
}

type FuncParam struct {
	Type string
	key  string
}

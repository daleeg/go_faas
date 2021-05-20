package main

const (
	GetTimeActive   = "get_time_active"   //获取时间的流程
	DoPrintActive   = "do_print_active"   //执行打印的流程
	PrintItemActive = "print_item_active" //执行分别打印
)

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

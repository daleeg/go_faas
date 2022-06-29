package util

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"
	"unsafe"
)

type PluginItem struct {
	Name           string
	PluginBaseInfo *PluginBaseInfoNode
	PluginItem     *plugin.Plugin
}

type Result interface {
	GetCode() error
	GetData() interface{}
}

type ErrorResult struct {
	err error
}

func (e ErrorResult) GetCode() error {
	return e.err
}
func (e ErrorResult) GetData() interface{} {
	return nil
}

type SuccessResult struct {
	data interface{}
}

func (s SuccessResult) GetCode() error {
	return nil
}
func (s SuccessResult) GetData() interface{} {
	return s.data
}

// 所有插件必须实现该方法
const BaseInfo = "PluginBaseInfo"

type eface struct {
	typ, val unsafe.Pointer
}

// LoadAllPlugin 将会过滤一次传入的targetFile,同时将so后缀的文件装载，并返回一个插件信息集合
func LoadAllPlugin(targetFile []string, collection map[string]PluginItem) []PluginItem {
	var res []PluginItem
	index := 1
	for _, fileItem := range targetFile {
		// 过滤插件文件
		if path.Ext(fileItem) == ".so" {
			fmt.Println("load plugin", index, ": ", fileItem)
			index += 1
			pluginFile, err := plugin.Open(fileItem)
			if err != nil {
				fmt.Println("An error occurred while load plugin : [" + fileItem + "]")
				fmt.Println(err)
			}
			//查找指定函数或符号
			targetFunc, err := pluginFile.Lookup(BaseInfo)
			if err != nil {
				fmt.Println("An error occurred while search target info : [" + fileItem + "]")
				fmt.Println(err)
				continue
			}

			baseInfo := (*PluginBaseInfoNode)((*eface)(unsafe.Pointer(&targetFunc)).val)
			// baseInfo, ok := targetFunc.(*PluginBaseInfoNode)
			// if !ok {
			// 	fmt.Println("Can find base info.")
			// 	continue
			// }
			// fmt.Println(baseInfo)
			// fmt.Println(4444)

			//采集插件信息
			pluginInfo := PluginItem{
				Name:           fileItem,
				PluginBaseInfo: baseInfo,
				PluginItem:     pluginFile,
			}
			filename := filepath.Base(fileItem)
			key := fmt.Sprintf("%s.%s.%s", strings.TrimSuffix(filename, filepath.Ext(filename)),
				baseInfo.Name,
				pluginInfo.PluginBaseInfo.Function.Name)
			collection[key] = pluginInfo
		}
	}
	return res
}

// DoInvokePlugin 会根据当前状态执行插件调用
func DoInvokePlugin(pluginFuncName string, args []interface{}) Result {
	fmt.Println(pluginFuncName)
	fmt.Println(pluginCollection)
	if pluginItem, ok := pluginCollection[pluginFuncName]; ok {
		// 判断流程
		funcName := pluginItem.PluginBaseInfo.Function.Name
		funcItem, err := pluginItem.PluginItem.Lookup(funcName)
		fmt.Println(funcName)
		if err != nil {
			fmt.Println("Can't find target func in [" + pluginItem.Name + "].")
			return ErrorResult{err}
		}
		fun := reflect.ValueOf(funcItem)
		params := &pluginItem.PluginBaseInfo.Function.Params
		in := make([]reflect.Value, len(*params))
		for k, param := range *params {
			switch param.Type {
			case "string":
				in[k] = reflect.ValueOf(args[k].(string))
				break
			case "int":
				in[k] = reflect.ValueOf(args[k].(int))
				break
			}
		}
		ret := fun.Call(in)
		return SuccessResult{ret}
	}
	print("Can't find [" + pluginFuncName + "]")
	return ErrorResult{errors.New("Can't find [" + pluginFuncName + "]")}
}

var pluginCollection map[string]PluginItem /*创建集合 */

func init() {
	// 读取plugin文件夹
	pluginsFiles := FindFile("plugin")
	pluginCollection = make(map[string]PluginItem)
	LoadAllPlugin(pluginsFiles, pluginCollection)
	fmt.Println(pluginCollection)
	fmt.Println("Process On ==========")
}

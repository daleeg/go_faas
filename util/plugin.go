package util

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"
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

// PluginPackageName 所有插件必须实现该方法
const PluginPackageName = "PackageName"


func listPluginMethod(p *plugin.Plugin, pluginNamePre string) (names []string) {
	pluginObj := reflect.ValueOf(*p)
	syms := pluginObj.FieldByName("syms")
	symsNames := syms.MapKeys()
	for _, symsName := range symsNames {
		name := symsName.String()
		if strings.HasPrefix(name, pluginNamePre) {
			names = append(names, name)
		}
	}
	return
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

			packageNameType, err := pluginFile.Lookup(PluginPackageName)
			if err != nil {
				fmt.Println("An error occurred while search target info : [" + PluginPackageName + "]")
				fmt.Println(err)
				continue
			}
			packageName, ok := packageNameType.(*string)
			if !ok {
				fmt.Println("Can find packageName ", packageNameType)
				continue
			}

			pluginMethods := listPluginMethod(pluginFile, "Plugin")
			filename := filepath.Base(fileItem)
			for _, pluginMethodName := range pluginMethods {
				fmt.Println("Plugin Method ", pluginMethodName)
				pluginMethod, err := pluginFile.Lookup(pluginMethodName)
				if err != nil {
					fmt.Println("An error occurred while search target info : [" + pluginMethodName + "]")
					fmt.Println(err)
					continue
				}

				fmt.Println("Plugin Method ", pluginMethodName)
				method := reflect.ValueOf(pluginMethod)

				if method.Kind() != reflect.Func {
					fmt.Println(pluginMethod, " is not function")
					continue
				}

				inParam := method.Type()
				fmt.Println("Plugin Method inParam", inParam)
				parameters := make([]reflect.Type, 0, inParam.NumIn())
				for i := 0; i < method.Type().NumIn(); i++ {
					arg := inParam.In(i)
					fmt.Printf("argument %d is %s[%s] type \n", i, arg.Kind(), arg.Name())
					parameters = append(parameters, arg)
				}

				baseInfo := PluginBaseInfoNode{
					Name:     pluginMethodName,
					Desc:     pluginMethodName,
					Function: method,
					Params:   parameters,
				}

				fmt.Println("baseInfo ", baseInfo)
				pluginInfo := PluginItem{
					Name:           fileItem,
					PluginBaseInfo: &baseInfo,
					PluginItem:     pluginFile,
				}
				fmt.Println("pluginInfo ", pluginInfo)
				key := fmt.Sprintf("%s.%s.%s", strings.TrimSuffix(filename, filepath.Ext(filename)),
					*packageName,
					baseInfo.Name)
				fmt.Println("key ", key)
				pluginCollection[key] = pluginInfo

			}
			
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
		fun := pluginItem.PluginBaseInfo.Function

		params := &pluginItem.PluginBaseInfo.Params
		in := make([]reflect.Value, len(*params))
		for k, param := range *params {
			switch param {
			case reflect.TypeOf(string("")):
				in[k] = reflect.ValueOf(args[k].(string))
				break
			case reflect.TypeOf(int(0)):
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

func ShowAllPlugins() {
	for name, pluginItem := range pluginCollection {
		fmt.Println(name, pluginItem)
	}
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

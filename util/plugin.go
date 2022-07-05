package util

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"path/filepath"
	"plugin"
	"reflect"
	"strconv"
	"strings"
)


// PluginPackageName 所有插件必须实现该方法
const PluginPackageName = "PackageName"

var (
	pluginCollection map[string]PluginItem /*创建集合 */
	logger  *logrus.Logger
)



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
			logger.Infoln("load plugin", index, ": ", fileItem)
			index += 1
			pluginFile, err := plugin.Open(fileItem)
			if err != nil {
				logger.Errorln("An error occurred while load plugin : [" + fileItem + "]")
				logger.Errorln(err)
			}

			packageNameType, err := pluginFile.Lookup(PluginPackageName)
			if err != nil {
				logger.Errorln("An error occurred while search target info : [" + PluginPackageName + "]")
				logger.Errorln(err)
				continue
			}
			packageName, ok := packageNameType.(*string)
			if !ok {
				logger.Errorln("Can find packageName ", packageNameType)
				continue
			}

			pluginMethods := listPluginMethod(pluginFile, "Plugin")
			filename := filepath.Base(fileItem)
			for _, pluginMethodName := range pluginMethods {
				pluginMethod, err := pluginFile.Lookup(pluginMethodName)
				if err != nil {
					logger.Errorln("An error occurred while search target info : [" + pluginMethodName + "]")
					logger.Errorln(err)
					continue
				}

				logger.Infoln("Plugin Method ", pluginMethodName)
				method := reflect.ValueOf(pluginMethod)

				if method.Kind() != reflect.Func {
					logger.Errorln(pluginMethod, " is not function")
					continue
				}

				methodParam := method.Type()
				numIn :=  methodParam.NumIn()
				inParameters := make([]reflect.Type, 0, numIn)
				for i := 0; i < numIn; i++ {
					arg := methodParam.In(i)
					logger.Infoln("argument %d is %s[%s] type \n", i, arg.Kind(), arg.Name())
					inParameters = append(inParameters, arg)
				}

				outReturns := make([]reflect.Type, 0, methodParam.NumOut())
				numOut :=  methodParam.NumOut()
				if numOut < 1 {
					logger.Errorln("outs length must greater than 0")
					continue
				}

				for i := 0; i < numOut; i++ {
					arg := methodParam.Out(i)
					logger.Infoln("out %d is %s[%s] type \n", i, arg.Kind(), arg.Name())
					outReturns = append(outReturns, arg)
				}
				if !outReturns[len(outReturns)-1].AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
					logger.Errorln("last output must be error")
					continue
				}
				if !outReturns[len(outReturns)-1].Implements(reflect.TypeOf((*error)(nil)).Elem()) {
					logger.Errorln("last output must be error")
					continue
				}

				baseInfo := PluginBaseInfoNode{
					Name:     pluginMethodName,
					Desc:     pluginMethodName,
					Function: method,
					Params:   inParameters,
					Returns: outReturns,
				}

				logger.Infoln("baseInfo ", baseInfo)
				pluginInfo := PluginItem{
					Name:           fileItem,
					PluginBaseInfo: &baseInfo,
				}
				logger.Infoln("pluginInfo ", pluginInfo)
				key := fmt.Sprintf("%s.%s.%s", strings.TrimSuffix(filename, filepath.Ext(filename)),
					*packageName,
					baseInfo.Name)
				logger.Infoln("key ", key)
				pluginCollection[key] = pluginInfo

			}
			
		}
	}
	return res
}

// DoInvokePlugin 会根据当前状态执行插件调用
func DoInvokePlugin(pluginFuncName string, args ...interface{}) Result {
	logger.Infoln(pluginFuncName)

	if pluginItem, ok := pluginCollection[pluginFuncName]; ok {
		// 判断流程
		fun := pluginItem.PluginBaseInfo.Function

		params := &pluginItem.PluginBaseInfo.Params
		in := make([]reflect.Value, len(*params))
		for k, param := range *params {
			switch param.Kind() {
			case reflect.String:
				in[k] = reflect.ValueOf(args[k].(string))
				break
			case reflect.Int:
				index := args[k]
				switch reflect.TypeOf(index).Kind() {
				case reflect.String:
					if v, e := strconv.Atoi(index.(string)); e != nil {
						return ErrorResult{errors.New("param type error")}
					} else {
						in[k] = reflect.ValueOf(v)
					}
					break
				case reflect.Float64:
					in[k] = reflect.ValueOf(int(index.(float64)))
					break
				}
				break
			}
		}

		ret := fun.Call(in)
		return SuccessResult{&ret, &pluginItem.PluginBaseInfo.Returns}
	}
	print("Can't find [" + pluginFuncName + "]")
	return ErrorResult{errors.New("Can't find [" + pluginFuncName + "]")}
}

func ShowAllPlugins() {
	for name, pluginItem := range pluginCollection {
		fmt.Println(name, pluginItem)
	}
}



func init() {
	// 读取plugin文件夹
	logger = GetLogger("util")
	pluginsFiles := FindFile("plugin")
	pluginCollection = make(map[string]PluginItem)
	LoadAllPlugin(pluginsFiles, pluginCollection)
	logger.Infoln(pluginCollection)
}

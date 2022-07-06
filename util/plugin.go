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
	"sync"
)

// PluginPackageName 所有插件必须实现该方法
const PluginPackageName = "PackageName"

type PluginServer struct {
	pluginCollection sync.Map // map[string]*service
	logger           *logrus.Logger
}

// NewPluginServer returns a new Server.
func NewPluginServer() *PluginServer {
	return &PluginServer{
		logger: GetLogger("pluginServer"),
	}
}

// DefaultServer is the default instance of *Server.
var DefaultServer = NewPluginServer()

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
func (s *PluginServer) LoadAllPlugin(targetFile []string) {
	index := 1
	for _, fileItem := range targetFile {
		// 过滤插件文件
		if path.Ext(fileItem) == ".so" {
			s.logger.Infoln("load plugin", index, ": ", fileItem)
			index += 1
			pluginFile, err := plugin.Open(fileItem)
			if err != nil {
				s.logger.Errorln("An error occurred while load plugin : [" + fileItem + "]")
				s.logger.Errorln(err)
			}

			packageNameType, err := pluginFile.Lookup(PluginPackageName)
			if err != nil {
				s.logger.Errorln("An error occurred while search target info : [" + PluginPackageName + "]")
				s.logger.Errorln(err)
				continue
			}
			packageName, ok := packageNameType.(*string)
			if !ok {
				s.logger.Errorln("Can find packageName ", packageNameType)
				continue
			}

			pluginMethods := listPluginMethod(pluginFile, "Plugin")
			filename := filepath.Base(fileItem)
			for _, pluginMethodName := range pluginMethods {
				pluginMethod, err := pluginFile.Lookup(pluginMethodName)
				if err != nil {
					s.logger.Errorln("An error occurred while search target info : [" + pluginMethodName + "]")
					s.logger.Errorln(err)
					continue
				}

				s.logger.Infoln("Plugin Method ", pluginMethodName)
				method := reflect.ValueOf(pluginMethod)

				if method.Kind() != reflect.Func {
					s.logger.Errorln(pluginMethod, " is not function")
					continue
				}

				methodParam := method.Type()
				numIn := methodParam.NumIn()
				inParameters := make([]reflect.Type, 0, numIn)
				for i := 0; i < numIn; i++ {
					arg := methodParam.In(i)
					s.logger.Infoln("argument %d is %s[%s] type \n", i, arg.Kind(), arg.Name())
					inParameters = append(inParameters, arg)
				}

				outReturns := make([]reflect.Type, 0, methodParam.NumOut())
				numOut := methodParam.NumOut()
				if numOut < 1 {
					s.logger.Errorln("outs length must greater than 0")
					continue
				}

				for i := 0; i < numOut; i++ {
					arg := methodParam.Out(i)
					s.logger.Infoln("out %d is %s[%s] type \n", i, arg.Kind(), arg.Name())
					outReturns = append(outReturns, arg)
				}
				if !outReturns[len(outReturns)-1].AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
					s.logger.Errorln("last output must be error")
					continue
				}
				if !outReturns[len(outReturns)-1].Implements(reflect.TypeOf((*error)(nil)).Elem()) {
					s.logger.Errorln("last output must be error")
					continue
				}

				baseInfo := PluginBaseInfoNode{
					Name:     pluginMethodName,
					Desc:     pluginMethodName,
					Function: method,
					Params:   inParameters,
					Returns:  outReturns,
				}

				s.logger.Infoln("baseInfo ", baseInfo)
				pluginInfo := PluginItem{
					Name:           fileItem,
					PluginBaseInfo: &baseInfo,
				}
				s.logger.Infoln("pluginInfo ", pluginInfo)
				key := fmt.Sprintf("%s.%s.%s", strings.TrimSuffix(filename, filepath.Ext(filename)),
					*packageName,
					baseInfo.Name)
				s.logger.Infoln("key ", key)

				if _, dup := s.pluginCollection.LoadOrStore(key, pluginInfo); dup {
					s.logger.Errorln("ready defined: " + key)
					s.logger.Errorln(err)
					continue
				}
			}

		}
	}
}

// DoInvokePlugin 会根据当前状态执行插件调用
func (s *PluginServer) DoInvokePlugin(pluginFuncName string, args ...interface{}) Result {
	s.logger.Infoln(pluginFuncName)

	if pluginNode, ok := s.pluginCollection.Load(pluginFuncName); ok {
		// 判断流程
		pluginItem := pluginNode.(PluginItem)
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

func (s *PluginServer) ShowAllPlugins() {
	s.pluginCollection.Range(
		func(name, pluginNode interface{}) bool {
			fmt.Println(name.(string), pluginNode.(PluginItem))
			return true
		})
}

func DoInvokePlugin(pluginFuncName string, args ...interface{}) Result {
	return DefaultServer.DoInvokePlugin(pluginFuncName, args...)
}

func ShowAllPlugins() {
	DefaultServer.ShowAllPlugins()
}

func init() {
	// 读取plugin文件夹
	pluginsFiles := FindFile("plugin")
	DefaultServer.LoadAllPlugin(pluginsFiles)
}

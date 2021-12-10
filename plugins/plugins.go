package plugins

import (
	"fmt"
	myPlugin "github.com/basebytes/plugins/plugin"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

const (
	pluginRootDirKey="GO_PLUGIN_DIR"
	defaultRootDir = "plugin"

	initFuncName = "NewPlugin"
	pluginFileSuffix=".so"
)
var myPlugins=make(map[string]myPlugin.Plugin)

func GetPlugin(name string) myPlugin.Plugin{
	return myPlugins[name]
}

func init(){
	var (
		_plugin *plugin.Plugin
		initFunc plugin.Symbol
		rootDir=os.Getenv(pluginRootDirKey)
	)
	if rootDir==""{
		rootDir= defaultRootDir
	}
	log.Printf("Plugin root dir: %s",rootDir)
	err:=filepath.Walk(rootDir, func(path string, info fs.FileInfo, e error)(err error) {
		if !info.IsDir()&&strings.HasSuffix(info.Name(), pluginFileSuffix) {
			log.Printf("Find plugin file %s",path)
			if _plugin,err=plugin.Open(path);err==nil{
				if initFunc,err=_plugin.Lookup(initFuncName);err==nil{
					var _plugins []myPlugin.Plugin
					if targetFunc, OK := initFunc.(func(string)([]myPlugin.Plugin,error));OK{
						if _plugins,err=targetFunc(filepath.Dir(path));err==nil{
							for _,_p:=range _plugins{
								name:=_p.Name()
								if _,OK:= myPlugins[name];OK{
									err=fmt.Errorf("Duplicate plugin name : %s ",name)
									return
								}else{
									myPlugins[name]= _p
									log.Printf("Load plugin [%s] from %s",name,path)
								}
							}
						}
					}
				}
			}
		}
		return
	})
	if err!=nil{
		panic(err)
	}
}
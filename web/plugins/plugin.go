package plugins

import "plugin"

var pluginStorage map[string]*Plugin

type Plugin struct {
	Name      string
	Version   string
	OnEnable  func()
	OnDisable func()
	plugin    *plugin.Plugin
}

func GetPlugin(name string) *Plugin {
	return pluginStorage[name]
}

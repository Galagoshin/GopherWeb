package events

import "github.com/Galagoshin/GoUtils/events"

const (
	StopApplicationEvent = events.EventName("StopApplicationEvent")
	RegisterRouteEvent   = events.EventName("RegisterRouteEvent")
	HotReloadEvent       = events.EventName("HotReloadEvent")
	EnablePluginEvent    = events.EventName("EnablePluginEvent")
	DisablePluginEvent   = events.EventName("DisablePluginEvent")
	StopWebServerEvent   = events.EventName("StopWebServerEvent")
	StartWebServerEvent  = events.EventName("StartWebServerEvent")
)

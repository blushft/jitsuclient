package events

import "github.com/blushft/jitsuclient/event"

func init() {
	event.RegisterType(EventTypeAction)
	event.RegisterType(EventTypeAlias)
	event.RegisterType(EventTypeGroup)
	event.RegisterType(EventTypeIdentify)
	event.RegisterType(EventTypePageview)
	event.RegisterType(EventTypeScreen)
	event.RegisterType(EventTypeSession)
	event.RegisterType(EventTypeTiming)
	event.RegisterType(EventTypeTrack)
}

func mergeOptions(opts ...[]event.Option) []event.Option {
	var mopts []event.Option

	for _, os := range opts {
		mopts = append(mopts, os...)
	}

	return mopts
}

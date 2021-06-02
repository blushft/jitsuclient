package events

import (
	"github.com/blushft/jitsuclient/event"
	"github.com/blushft/jitsuclient/event/contexts"
)

const EventTypeScreen event.Type = "screen"

func Screen(screen *contexts.Screen, opts ...event.Option) *event.Event {
	eopts := mergeOptions(
		[]event.Option{
			event.EventType(EventTypeScreen),
			event.WithContext(screen),
			event.WithValidator(screenValidator()),
		},
		opts,
	)

	return event.New("screen", eopts...)
}

func screenValidator() event.Validator {
	return event.NewValidator(
		event.WithRule("valid_screen", event.ContextValid(contexts.ContextScreen)),
	)
}

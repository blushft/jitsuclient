package events

import (
	"github.com/blushft/jitsuclient/event"
	"github.com/blushft/jitsuclient/event/contexts"
)

const EventTypeAction event.Type = "action"

func Action(action *contexts.Action, opts ...event.Option) *event.Event {
	eopts := mergeOptions(
		[]event.Option{
			event.EventType(EventTypeAction),
			event.WithContext(action),
			event.WithValidator(actionValidator()),
		},
		opts,
	)

	return event.New("action", eopts...)
}

func actionValidator() event.Validator {
	return event.NewValidator(
		event.WithRule("valid_action", event.ContextValid(contexts.ContextAction)),
	)
}

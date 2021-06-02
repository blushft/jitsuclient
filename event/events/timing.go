package events

import (
	"github.com/blushft/jitsuclient/event"
	"github.com/blushft/jitsuclient/event/contexts"
)

const EventTypeTiming event.Type = "timing"

func Timing(t *contexts.Timing, opts ...event.Option) *event.Event {
	eopts := mergeOptions(
		[]event.Option{
			event.EventType(EventTypeTiming),
			event.WithContext(t),
			event.WithValidator(timingValidator()),
		},
		opts,
	)

	return event.New("timing", eopts...)
}

func timingValidator() event.Validator {
	return event.NewValidator(
		event.WithRule("valid_timing", event.ContextValid(contexts.ContextTiming)),
	)
}

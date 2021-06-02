package events

import (
	"github.com/blushft/jitsuclient/event"
	"github.com/blushft/jitsuclient/event/contexts"
)

const EventTypeSession event.Type = "session"

func Session(sess *contexts.Session, opts ...event.Option) *event.Event {
	eopts := mergeOptions(
		[]event.Option{
			event.EventType(EventTypeSession),
			event.WithContext(sess),
			event.WithValidator(sessionValidator()),
		},
		opts,
	)

	return event.New("session", eopts...)
}

func sessionValidator() event.Validator {
	return event.NewValidator(
		event.WithRule("valid_session", event.ContextValid(contexts.ContextSession)),
	)
}

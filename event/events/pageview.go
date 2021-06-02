package events

import (
	"github.com/blushft/jitsuclient/event"
	"github.com/blushft/jitsuclient/event/contexts"
)

const EventTypePageview event.Type = "pageview"

func Pageview(page *contexts.Page, opts ...event.Option) *event.Event {
	eopts := mergeOptions(
		[]event.Option{
			event.EventType(EventTypePageview),
			event.WithContext(page),
			event.WithValidator(pageviewValidator()),
		},
		opts,
	)

	return event.New("pageview", eopts...)
}

func pageviewValidator() event.Validator {
	return event.NewValidator(
		event.WithRule("valid_page", event.ContextValid(contexts.ContextPage)),
	)
}

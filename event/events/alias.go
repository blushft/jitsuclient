package events

import (
	"github.com/blushft/jitsuclient/event"
	"github.com/blushft/jitsuclient/event/contexts"
)

const EventTypeAlias event.Type = "alias"

func Alias(alias *contexts.Alias, user *contexts.User, opts ...event.Option) *event.Event {
	eopts := mergeOptions(
		[]event.Option{
			event.EventType(EventTypeAlias),
			event.WithContext(user),
			event.WithContext(alias),
			event.WithValidator(aliasValidator()),
		},
		opts,
	)

	return event.New("alias", eopts...)
}

func aliasValidator() event.Validator {
	return event.NewValidator(
		event.WithRule("valid_group", event.ContextValid(contexts.ContextAlias)),
		event.WithRule("valid_user", event.ContextValid(contexts.ContextUser)),
	)
}

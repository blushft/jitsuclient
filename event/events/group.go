package events

import (
	"github.com/blushft/jitsuclient/event"
	"github.com/blushft/jitsuclient/event/contexts"
)

const EventTypeGroup event.Type = "group"

func Group(grp *contexts.Group, usr *contexts.User, opts ...event.Option) *event.Event {
	if grp != nil {
		opts = append(opts, event.WithContext(grp))
	}

	if usr != nil {
		opts = append(opts, event.WithContext(usr))
	}

	opts = append(opts, event.EventType(EventTypeGroup), event.WithValidator(groupValidator()))

	return event.New("group", opts...)
}

func groupValidator() event.Validator {
	return event.NewValidator(
		event.WithRule("group_context", event.HasContext(contexts.ContextGroup)),
		event.WithRule("user_context", event.HasContext(contexts.ContextUser)),
	)
}

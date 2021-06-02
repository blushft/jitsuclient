package events

import (
	"github.com/blushft/jitsuclient/event"
	"github.com/blushft/jitsuclient/event/contexts"
)

const EventTypeIdentify event.Type = "identify"

func Identify(user *contexts.User, opts ...event.Option) *event.Event {
	if user != nil {
		opts = append(opts, event.WithContext(user))
	}

	opts = append(opts, event.EventType(EventTypeIdentify), event.WithValidator(identifyValidator()))

	return event.New("identify", opts...)
}

func identifyValidator() event.Validator {
	return event.NewValidator(
		event.WithRule("valid_user", event.ContextValid(contexts.ContextUser)),
	)
}

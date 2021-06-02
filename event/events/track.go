package events

import (
	"github.com/blushft/jitsuclient/event"
)

const EventTypeTrack event.Type = "track"

func Track(evt string, opts ...event.Option) *event.Event {
	opts = append(opts, event.EventType(EventTypeTrack))

	return event.New(evt, opts...)
}

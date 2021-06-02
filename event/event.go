package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Type string

const EventTypeEvent Type = "event"

type Event struct {
	ID         string `json:"id" structs:"id" mapstructure:"id"`
	TrackingID string `json:"trackingId" structs:"trackingID" mapstructure:"trackingID"`
	UserID     string `json:"userId,omitempty" structs:"userID,omitempty" mapstructure:"userID,omitempty"`
	Anonymous  bool   `json:"anonymous" structs:"anonymous" mapstructure:"anonymous"`
	GroupID    string `json:"groupId,omitempty" structs:"groupID,omitempty" mapstructure:"groupID,omitempty"`
	SessionID  string `json:"sessionId,omitempty" structs:"sessionID,omitempty" mapstructure:"sessionID,omitempty"`
	DeviceID   string `json:"deviceId,omitempty" structs:"deviceID,omitempty" mapstructure:"deviceID,omitempty"`

	Event          string    `json:"event" structs:"event" mapstructure:"event"`
	Type           Type      `json:"event_type" structs:"event_type" mapstructure:"event_type"`
	NonInteractive bool      `json:"nonInteractive,omitempty" structs:"nonInteractive,omitempty" mapstructure:"nonInteractive,omitempty"`
	Channel        string    `json:"channel,omitempty" structs:"channel,omitempty" mapstructure:"channel,omitempty"`
	Platform       string    `json:"platform,omitempty" structs:"platform,omitempty" mapstructure:"platform,omitempty"`
	Timestamp      time.Time `json:"timestamp" structs:"timestamp" mapstructure:"timestamp"`
	Context        Contexts  `json:"context,omitempty" structs:"context,omitempty" mapstructure:"context,omitempty"`

	validator Validator
}

func New(evt string, opts ...Option) *Event {
	e := &Event{
		ID:        uuid.New().String(),
		Event:     evt,
		Type:      EventTypeEvent,
		Timestamp: time.Now(),
		Context:   make(map[string]Context),
	}

	for _, o := range opts {
		o(e)
	}

	return e
}

func Empty() *Event {
	return &Event{
		ID:      uuid.New().String(),
		Context: make(map[string]Context),
	}
}

func (e *Event) SetContext(ctx Context) {
	e.Context[string(ctx.Type())] = ctx
}

func (e *Event) SetContexts(ctx ...Context) {
	for _, cc := range ctx {
		e.SetContext(cc)
	}
}

func (e *Event) Apply(opts ...Option) {
	for _, o := range opts {
		o(e)
	}
}

func (e *Event) Render(opts ...Option) ([]byte, error) {
	e.Apply(opts...)

	if len(e.Event) == 0 {
		e.Event = string(e.Type)
	}

	return json.Marshal(e)
}

func (e *Event) Validate() bool {
	if ok := evtreg[e.Type]; !ok {
		return false
	}

	if e.validator != nil {
		return e.validator.Validate(e)
	}

	return true
}

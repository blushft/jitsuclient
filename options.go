package client

import (
	"time"

	"github.com/blushft/jitsuclient/event"
	"github.com/blushft/jitsuclient/event/contexts"
)

type Options struct {
	CollectorURL  string
	ServerKey     string
	S2S           bool
	ClientHeaders map[string]string
	Debug         bool

	AppInfo    *contexts.App
	Platform   string
	TrackingID string

	QueueBuffer   int
	FlushInterval time.Duration
	FlushCount    int
}

func defaultOptions(opts ...Option) Options {
	options := Options{
		CollectorURL:  "http://localhost:8001",
		ClientHeaders: make(map[string]string),
		QueueBuffer:   25,
		FlushInterval: time.Second * 10,
		FlushCount:    100,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

type Option func(*Options)

func (o Options) EventOptions() []event.Option {
	evtOpts := []event.Option{
		event.WithContext(&contexts.Library{Name: "jitsuclient", Version: "v0.1.0"}),
	}

	if o.AppInfo != nil {
		evtOpts = append(evtOpts, event.WithContext(o.AppInfo))
	}

	if len(o.TrackingID) > 0 {
		evtOpts = append(evtOpts, event.TrackingID(o.TrackingID))
	}

	if len(o.Platform) > 0 {
		evtOpts = append(evtOpts, event.Platform(o.Platform))
	} else {
		evtOpts = append(evtOpts, event.Platform("srv"))
	}

	return evtOpts
}

func (o Options) apiPath() string {
	if o.S2S {
		return "/api/v1/s2s/event"
	}

	return "/api/v1/event"
}

func (o Options) apiQueryParams() map[string]string {
	m := make(map[string]string, 0)

	if !o.S2S {
		m["token"] = o.ServerKey
	}

	return m
}

func (o Options) clientHeaders() map[string]string {
	m := map[string]string{
		"Content-Type": "application/json",
	}

	if o.S2S && len(o.ServerKey) > 0 {
		m["X-Auth-Token"] = o.ServerKey
	}

	for k, v := range o.ClientHeaders {
		m[k] = v
	}

	return m
}

func CollectorURL(u string) Option {
	return func(o *Options) {
		o.CollectorURL = u
	}
}

func Debug() Option {
	return func(o *Options) {
		o.Debug = true
	}
}

func WithClientHeader(k, v string) Option {
	return func(o *Options) {
		o.ClientHeaders[k] = v
	}
}

func SetAppInfo(app *contexts.App) Option {
	return func(o *Options) {
		o.AppInfo = app
	}
}

func ServerKey(key string) Option {
	return func(o *Options) {
		o.ServerKey = key
	}
}

func TrackingID(id string) Option {
	return func(o *Options) {
		o.TrackingID = id
	}
}

package client

import (
	"time"

	"github.com/blushft/jitsuclient/event"
	"github.com/blushft/jitsuclient/event/contexts"
)

const (
	minInterval = time.Millisecond * 100
)

type Options struct {
	CollectorURL  string
	ServerKey     string
	S2S           bool
	Bulk          bool
	ClientHeaders map[string]string
	Debug         bool
	Strict        bool

	AppInfo    *contexts.App
	Platform   string
	TrackingID string

	QueueBuffer   int
	FlushInterval time.Duration
	FlushCount    int
	MaxRetries    int
}

func defaultOptions(opts ...Option) Options {
	options := Options{
		CollectorURL:  "http://localhost:8000",
		ClientHeaders: make(map[string]string),
		QueueBuffer:   100,
		FlushInterval: time.Second * 1,
		FlushCount:    250,
		MaxRetries:    10,
	}

	for _, o := range opts {
		o(&options)
	}

	if options.FlushInterval < minInterval {
		options.FlushInterval = minInterval
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
	if o.Bulk {
		return "/api/v1/events/bulk"
	}

	if o.S2S {
		return "/api/v1/s2s/event"
	}

	return "/api/v1/event"
}

func (o Options) apiQueryParams() map[string]string {
	m := make(map[string]string)

	if !o.S2S {
		m["token"] = o.ServerKey
	}

	return m
}

func (o Options) clientHeaders() map[string]string {
	m := map[string]string{
		"Content-Type": o.apiContentType(),
	}

	if o.S2S && len(o.ServerKey) > 0 {
		m["X-Auth-Token"] = o.ServerKey
	}

	for k, v := range o.ClientHeaders {
		m[k] = v
	}

	return m
}

func (o Options) apiContentType() string {
	if o.Bulk {
		return "multipart/form-data"
	}

	return "application/json"
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

func Strict(b bool) Option {
	return func(o *Options) {
		o.Strict = b
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

func QueueBuffer(buf int) Option {
	return func(o *Options) {
		o.QueueBuffer = buf
	}
}

func FlushInterval(d time.Duration) Option {
	return func(o *Options) {
		o.FlushInterval = d
	}
}

func FlushCount(i int) Option {
	return func(o *Options) {
		o.FlushCount = i
	}
}

func S2S() Option {
	return func(o *Options) {
		o.S2S = true
	}
}

func Bulk() Option {
	return func(o *Options) {
		o.Bulk = true
		o.S2S = true
		o.FlushInterval = time.Second * 2
	}
}

func MaxRetries(i int) Option {
	return func(o *Options) {
		o.MaxRetries = i
	}
}

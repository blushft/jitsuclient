package client

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/blushft/jitsuclient/event"
	"github.com/blushft/jitsuclient/event/contexts"
	"github.com/blushft/jitsuclient/event/events"
	"github.com/hashicorp/go-multierror"

	"gopkg.in/resty.v1"
)

type Client struct {
	options Options
	httpc   *resty.Client
	store   Store

	q  chan *event.Event
	cl chan struct{}

	copts []event.Option
}

func New(opts ...Option) (*Client, error) {
	options := defaultOptions(opts...)

	httpc := resty.New().SetHostURL(options.CollectorURL)
	store, err := NewMemStore()
	if err != nil {
		return nil, err
	}

	t := &Client{
		options: options,
		httpc:   httpc,
		store:   store,
		q:       make(chan *event.Event, options.QueueBuffer),
		cl:      make(chan struct{}),
		copts:   options.EventOptions(),
	}

	go t.run()

	return t, nil
}

func (t *Client) SetOption(opt event.Option) *Client {
	t.copts = append(t.copts, opt)
	return t
}

func (t *Client) SetUser(usr *contexts.User) *Client {
	return t.SetOption(event.WithContext(usr))
}

func (t *Client) SetGroup(grp *contexts.Group) *Client {
	return t.SetOption(event.WithContext(grp))
}

func (t *Client) Queue(evt *event.Event) {
	evt.Apply(t.eventOptions()...)

	t.q <- evt
}

func (t *Client) Action(a *contexts.Action, opts ...event.Option) {
	t.Queue(events.Action(a, opts...))
}

func (t *Client) Track(evt string, opts ...event.Option) {
	t.Queue(events.Track(evt, opts...))
}

func (t *Client) Identify(u *contexts.User, opts ...event.Option) {
	t.Queue(events.Identify(u, opts...))
}

func (t *Client) Alias(alias *contexts.Alias, usr *contexts.User, opts ...event.Option) {
	t.Queue(events.Alias(alias, usr, opts...))
}

func (t *Client) Page(page *contexts.Page, opts ...event.Option) {
	t.Queue(events.Pageview(page, opts...))
}

func (t *Client) Screen(screen *contexts.Screen, opts ...event.Option) {
	t.Queue(events.Screen(screen, opts...))
}

func (t *Client) Session(sess *contexts.Session, opts ...event.Option) {
	t.Queue(events.Session(sess, opts...))
}

func (t *Client) Group(g *contexts.Group, u *contexts.User, opts ...event.Option) {
	t.Queue(events.Group(g, u, opts...))
}

func (t *Client) Timing(te *contexts.Timing, opts ...event.Option) {
	t.Queue(events.Timing(te, opts...))
}

func (t *Client) Flush() {
	start := time.Now()
	c, err := t.emit()
	if err != nil {
		log.Printf("error emitting events: %v\n", err)
	}

	if c > 0 && t.options.Debug {
		log.Printf("flush complete. count=%d dur=%dms", c, time.Since(start).Milliseconds())
	}
}

func (t *Client) Close() {
	close(t.cl)

	t.Flush()
}

func (t *Client) eventOptions(opts ...event.Option) []event.Option {
	var eopts []event.Option

	eopts = append(eopts, t.copts...)
	eopts = append(eopts, opts...)

	return eopts
}

func (t *Client) run() {
	tick := time.NewTicker(t.options.FlushInterval)
	defer tick.Stop()

	for {
		select {
		case e := <-t.q:
			if t.options.Strict {
				if !e.Validate() {
					log.Printf("event failed validation")
					continue
				}
			}

			if err := t.store.Set(e); err != nil {
				log.Printf("error storing event: %s\n", err.Error())
				continue
			}

			if t.store.Count() >= t.options.FlushCount {
				t.logDebug("flushing events on flush count")
				t.Flush()
			}
		case <-tick.C:
			t.logDebug("flushing events on interval")
			t.Flush()
		case <-t.cl:
			return
		}
	}
}

type errorHandler func(*StoreEvent, error) error
type eventSender func([]*StoreEvent, errorHandler) (int, error)

func (t *Client) emit() (int, error) {
	if t.store.Count() == 0 {
		return 0, nil
	}

	start := time.Now()
	send := t.sendEach
	if t.options.Bulk {
		send = t.sendAll
	}

	evts, err := t.store.GetAll()
	if err != nil {
		return 0, err
	}

	i, err := send(evts, t.handleSendFailed)
	if err != nil {
		return i, err
	}

	t.logDebug("emitted %d events in %dms", i, time.Since(start).Milliseconds())

	return i, nil
}

func (t *Client) sendEach(evts []*StoreEvent, h errorHandler) (int, error) {
	start := time.Now()
	i := 0

	for _, e := range evts {
		if serr := t.send(e.Event); serr != nil {
			if err := h(e, serr); err != nil {
				log.Println(err.Error())
			}

			continue
		}

		i++

		if err := t.store.Remove(e); err != nil {
			log.Printf("error removing event: %v", err)
		}
	}

	t.logDebug("send complete: events=%d dur=%dms", i, time.Since(start).Milliseconds())

	return i, nil
}

func (t *Client) sendAll(evts []*StoreEvent, h errorHandler) (int, error) {
	start := time.Now()

	bulk := make([][]byte, len(evts))
	for c, e := range evts {
		bulk[c] = e.Event
	}

	onFail := func(evts []*StoreEvent, serr error) error {
		var res error
		for _, e := range evts {
			if err := h(e, serr); err != nil {
				res = multierror.Append(res, err)
			}
		}

		return res
	}

	req := bytes.Join(bulk, []byte("\n"))
	req = append(req, '\n') 

	if err := t.sendBulk(req); err != nil {
		log.Printf("bulk send failed with error: %v", err)
		if ferr := onFail(evts, err); ferr != nil {
			return 0, ferr
		}

		return 0, err
	}

	i := 0
	for _, e := range evts {
		if err := t.store.Remove(e); err != nil {
			log.Printf("error removing event: %v", err)
		}
		i++
	}

	t.logDebug("bulk send complete: count=%d dur=%dms", i, time.Since(start).Milliseconds())

	return i, nil
}

func (t *Client) send(e []byte) error {
	headers := t.options.clientHeaders()

	resp, err := t.httpc.R().
		SetBody(e).
		SetHeaders(headers).
		SetQueryParams(t.options.apiQueryParams()).
		Post(t.options.apiPath())
	if err != nil {
		return err
	}

	if resp != nil && resp.StatusCode() > 299 {
		return fmt.Errorf("send failed with status code %d", resp.StatusCode())
	}

	return err
}

func (t *Client) sendBulk(e []byte) error {
	headers := t.options.clientHeaders()

	resp, err := t.httpc.R().
		SetFileReader("file", "file", bytes.NewReader(e)).
		SetHeaders(headers).
		SetQueryParams(t.options.apiQueryParams()).
		Post(t.options.apiPath())
	if err != nil {
		return err
	}

	if resp != nil && resp.StatusCode() > 299 {
		return fmt.Errorf("bulk send failed with status code %d", resp.StatusCode())
	}

	return nil
}

func (t *Client) handleSendFailed(e *StoreEvent, sendErr error) error {
	e.Attempted = true
	e.Attempts++
	e.LastAttempt = time.Now()

	t.logDebug("event failed to send: %v", sendErr)

	if t.options.MaxRetries > 0 && e.Attempts > t.options.MaxRetries {
		t.logDebug("max retries reached for event, removing")
		if err := t.store.Remove(e); err != nil {
			log.Printf("error removing event: %v", err)
		}
	} else {
		if err := t.store.Update(e); err != nil {
			return fmt.Errorf("error updating event: %w", err)
		}
	}

	return nil
}

func (t *Client) logDebug(msg string, args ...interface{}) {
	if t.options.Debug {
		log.Printf(msg, args...)
	}

}

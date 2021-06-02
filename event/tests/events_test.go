package events_test

import (
	"testing"

	"github.com/blushft/jitsuclient/event"
	"github.com/blushft/jitsuclient/event/contexts"
	"github.com/blushft/jitsuclient/event/events"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type EventSuite struct {
	suite.Suite
}

func TestRunEventSuite(t *testing.T) {
	suite.Run(t, new(EventSuite))
}

func (s *EventSuite) TestNewEvent() {

	evt := events.Action(&contexts.Action{
		Action:   "do_something",
		Category: "stuff",
	},
		event.TrackingID("abc123"),
		event.WithContext(&contexts.User{UserID: "my@guy.com"}),
	)

	s.Require().True(evt.Validate())

	s.Equal(2, len(evt.Context.List()))
}

func (s *EventSuite) TestContextIter() {
	grp := &contexts.Group{
		ID:   uuid.NewString(),
		Name: "testers",
	}

	usr := &contexts.User{
		UserID: uuid.NewString(),
		Name:   "my_guy",
		Traits: contexts.Traits{
			Email: "my@guy.com",
		},
	}
	evt := events.Group(grp, usr,
		event.TrackingID("abc123"),
		event.DeviceID("some_hostname"),
		event.WithContext(contexts.NewNetwork("192.168.0.1")),
	)

	s.Require().True(evt.Validate())

	it := evt.Context.Iter()
	count := 0
	for ctx := it.First(); ctx != nil; ctx = it.Next() {
		s.T().Logf("visiting %v", ctx.Type())
		count++
	}

	s.Equal(3, count)
}

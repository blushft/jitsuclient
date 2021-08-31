package client

import (
	"testing"
	"time"

	"github.com/blushft/jitsuclient/event"
	"github.com/blushft/jitsuclient/event/contexts"
	"github.com/stretchr/testify/suite"
)

type ClientSuite struct {
	suite.Suite
	c *Client
}

func TestRunClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

func (s *ClientSuite) SetupSuite() {
	tr, err := New(
		TrackingID("tracker_test"),
		ServerKey("s2s.f94blb9erxbasgzmdnnlwq.42ncq1nh3h4t1kz40amjx"),
		SetAppInfo(&contexts.App{
			Name:    "tracker_test",
			Version: "v0.1.0",
		}),
		Bulk(),
		Debug())

	s.Require().NoError(err)

	s.c = tr
}

func (s *ClientSuite) TearDownSuite() {
	s.c.Close()
}

func (s *ClientSuite) TestTrackAction() {
	a := &contexts.Action{
		Category: "tests",
		Action:   "Test Action",
		Label:    "test_track",
		Property: "test_val",
		Value:    99,
	}

	s.c.Action(a)

	time.Sleep(time.Second * 5)
}

func (s *ClientSuite) TestBulkEvents() {
	for i := 0; i < 250; i++ {
		s.c.Queue(bulkEvent(i))
	}

	time.Sleep(time.Second * 5)
}

func bulkEvent(i int) *event.Event {
	return event.New("bulk_event",
		event.WithContext(&contexts.Action{
			Category: "current_test",
			Action:   "Bulk Load",
			Label:    "test_bulk",
			Property: "value",
			Value:    i,
		}))
}

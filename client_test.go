package client

import (
	"testing"
	"time"

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
		ServerKey("s2s.3qaiwf2e92fmm3ido942ph.v7878v1hpyiw58j1rttmhs"),
		SetAppInfo(&contexts.App{
			Name:    "tracker_test",
			Version: "v0.1.0",
		}))

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

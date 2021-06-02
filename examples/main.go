package main

import (
	"log"
	"time"

	client "github.com/blushft/jitsuclient"
	"github.com/blushft/jitsuclient/event/contexts"
	"github.com/google/uuid"
)

func main() {
	c, err := client.New(
		client.Debug(),
		client.TrackingID("tracker_test"),
		//client.ServerKey("s2s.3qaiwf2e92fmm3ido942ph.v7878v1hpyiw58j1rttmhs"),
		client.ServerKey("js.3qaiwf2e92fmm3ido942ph.ssv0egh3s7aebssf4297zq"),
		client.SetAppInfo(&contexts.App{
			Name:    "tracker_test",
			Version: "v0.1.0",
		}))
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	usr := &contexts.User{
		AnonID: uuid.New().String(),
	}

	grp := &contexts.Group{
		ID:   uuid.NewString(),
		Name: "anon_group",
	}

	c.SetUser(usr).Identify(nil)
	c.SetGroup(grp).Group(nil, nil)

	c.Action(&contexts.Action{Action: "click", Category: "user_action", Label: "submit_button"})
	t := c.TimingStart("checkout", "viewing_cart", "user_attention")
	time.Sleep(time.Second * 5)

	c.Timing(t.End())
	c.Flush()
}

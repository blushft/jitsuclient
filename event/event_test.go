package event_test

import (
	"fmt"
	"testing"

	"github.com/blushft/jitsuclient/event"
	"github.com/blushft/jitsuclient/event/contexts"
)

func TestCreateEvent(t *testing.T) {
	evt := event.New("testing", event.WithContext(contexts.NewAction("testing", "test")))

	b, err := evt.Render()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
}

package function

import (
	"context"
	"embed"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
)

type StaticFileRequest struct {
	FilePath string `json:"filePath"`
}

// content holds the static web server content
//
//go:embed public/*
var content embed.FS

// Handle an event.
func Handle(ctx context.Context, e event.Event) (*event.Event, error) {
	/*
	 * YOUR CODE HERE
	 *
	 * Try running `go test`.  Add more test as you code in `handle_test.go`.
	 */

	req := &StaticFileRequest{}
	err := e.DataAs(req)
	if err != nil {
		return nil, err
	}

	file, err := content.ReadFile(req.FilePath)
	if err != nil {
		return nil, err
	}

	event := cloudevents.NewEvent()
	event.SetSource("knext/staticserver")
	event.SetType("knext.static.response")
	event.SetData("text/html", string(file))

	fmt.Println(event) // echo to local output
	return &event, nil // echo to caller
}

/*
Other supported function signatures:

	Handle()
	Handle() error
	Handle(context.Context)
	Handle(context.Context) error
	Handle(event.Event)
	Handle(event.Event) error
	Handle(context.Context, event.Event)
	Handle(context.Context, event.Event) error
	Handle(event.Event) *event.Event
	Handle(event.Event) (*event.Event, error)
	Handle(context.Context, event.Event) *event.Event
	Handle(context.Context, event.Event) (*event.Event, error)

*/

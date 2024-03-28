package function

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
)

// TestHandle ensures that Handle accepts a valid CloudEvent without error.
func TestHandle(t *testing.T) {
	// Assemble
	e := event.New()
	e.SetID("id")
	e.SetType("type")
	e.SetSource("source")
	data, err := json.Marshal(map[string]string{"filePath": "index.html"})
	if err != nil {
		t.Fatal(err)
	}

	e.SetData("application/json", data)

	// Act
	resp, err := Handle(context.Background(), e)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if resp == nil {
		t.Errorf("received nil event") // fail on nil
	}
	if string(resp.Data()) != "<h1>Hello world</h1>\n" {
		t.Errorf("the received event expected data to be '<h1>Hello world</h1>\n', got '%s'", resp.Data())
	}
}

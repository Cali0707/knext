package function

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
)

type StaticFileRequest struct {
	FilePath string `json:"filePath"`
}

// content holds the static web server content
//
//go:embed static/*
var content embed.FS

// Handle an event.
func Handle(_ context.Context, e event.Event) (*event.Event, error) {

	req := &StaticFileRequest{}
	err := e.DataAs(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		fmt.Printf("%v\n", req)
		return nil, err
	}

	file, err := content.Open(fmt.Sprintf("static/%s", req.FilePath))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}

	content, contentType, err := readFile(file)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}

	event := cloudevents.NewEvent()
	event.SetSource("knext/staticserver")
	event.SetType("knext.static.success")
	err = event.SetData(contentType, content)

	fmt.Println(event) // echo to local output
	return &event, err // echo to caller
}

func readFile(file fs.File) ([]byte, string, error) {
	var size int
	if info, err := file.Stat(); err == nil {
		size64 := info.Size()
		if int64(int(size64)) == size64 {
			size = int(size64)
		}
	}

	var contentType string
	data := make([]byte, 0, size+1)
	for {
		if len(data) >= cap(data) {
			d := append(data[:cap(data)], 0)
			data = d[:len(data)]
		}
		n, err := file.Read(data[len(data):cap(data)])
		data = data[:len(data)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
				contentType = http.DetectContentType(data)
			}
			return data, contentType, err
		}
	}

}

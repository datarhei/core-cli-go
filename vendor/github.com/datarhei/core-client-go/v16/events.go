package coreclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/datarhei/core-client-go/v16/api"

	"encoding/json"
)

func (r *restclient) Events(ctx context.Context, filters api.LogEventFilters) (<-chan api.LogEvent, error) {
	var buf bytes.Buffer

	e := json.NewEncoder(&buf)
	e.Encode(filters)

	header := make(http.Header)
	header.Set("Accept", "application/x-json-stream")
	header.Set("Connection", "close")

	stream, err := r.stream(ctx, "POST", "/v3/events", nil, header, "application/json", &buf)
	if err != nil {
		return nil, err
	}

	channel := make(chan api.LogEvent, 128)

	go func(stream io.ReadCloser, ch chan<- api.LogEvent) {
		defer stream.Close()
		defer close(channel)

		decoder := json.NewDecoder(stream)

		for decoder.More() {
			var event api.LogEvent
			if err := decoder.Decode(&event); err == io.EOF {
				event.Component = "eof"
			} else if err != nil {
				event.Component = "error"
				event.Message = err.Error()
			}

			// Don't emit keepalives
			if event.Component == "keepalive" {
				continue
			}

			if event.Component == "" {
				event.Component = "eof"
			}

			ch <- event

			if event.Component == "eof" || event.Component == "error" {
				break
			}
		}

		ch <- api.LogEvent{
			Component: "eof",
		}
	}(stream, channel)

	return channel, nil
}

func (r *restclient) MediaEvents(ctx context.Context, storage, pattern string) (<-chan api.MediaEvent, error) {
	header := make(http.Header)
	header.Set("Accept", "application/x-json-stream")
	header.Set("Connection", "close")

	query := &url.Values{}
	query.Set("glob", pattern)

	stream, err := r.stream(ctx, "POST", "/v3/events/media/"+url.PathEscape(storage), query, header, "", nil)
	if err != nil {
		return nil, err
	}

	channel := make(chan api.MediaEvent, 128)

	go func(stream io.ReadCloser, ch chan<- api.MediaEvent) {
		defer stream.Close()
		defer close(channel)

		decoder := json.NewDecoder(stream)

		for decoder.More() {
			event := api.MediaEvent{}

			if err := decoder.Decode(&event); err == io.EOF {
				event.Action = "eof"
			} else if err != nil {
				event.Action = "error"
				event.Name = err.Error()
			}

			// Don't emit keepalives
			if event.Action == "keepalive" {
				continue
			}

			if event.Action == "" {
				event.Action = "eof"
			}

			ch <- event

			if event.Action == "eof" || event.Action == "error" {
				break
			}
		}

		ch <- api.MediaEvent{
			Action: "eof",
		}
	}(stream, channel)

	return channel, nil
}

func (r *restclient) ProcessEvents(ctx context.Context, filters api.ProcessEventFilters) (<-chan api.ProcessEvent, error) {
	var buf bytes.Buffer

	e := json.NewEncoder(&buf)
	e.Encode(filters)

	header := make(http.Header)
	header.Set("Accept", "application/x-json-stream")
	header.Set("Connection", "close")

	stream, err := r.stream(ctx, "POST", "/v3/events/process", nil, header, "application/json", &buf)
	if err != nil {
		return nil, err
	}

	channel := make(chan api.ProcessEvent, 128)

	go func(stream io.ReadCloser, ch chan<- api.ProcessEvent) {
		defer stream.Close()
		defer close(channel)

		decoder := json.NewDecoder(stream)

		for decoder.More() {
			event := api.ProcessEvent{}

			if err := decoder.Decode(&event); err == io.EOF {
				event.Type = "error"
				event.Line = "EOF"
			} else if err != nil {
				event.Type = "error"
				event.Line = err.Error()
			}

			// Don't emit keepalives
			if event.Type == "keepalive" {
				continue
			}

			if event.Type == "" {
				event.Type = "eof"
			}

			ch <- event

			if event.Type == "eof" || event.Type == "error" {
				break
			}
		}

		ch <- api.ProcessEvent{
			Type: "eof",
		}
	}(stream, channel)

	return channel, nil
}

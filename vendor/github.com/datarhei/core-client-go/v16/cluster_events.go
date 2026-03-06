package coreclient

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/datarhei/core-client-go/v16/api"

	"encoding/json"
)

func (r *restclient) ClusterLogEvents(ctx context.Context, filters api.LogEventFilters) (<-chan api.LogEvent, error) {
	var buf bytes.Buffer

	e := json.NewEncoder(&buf)
	e.Encode(filters)

	header := make(http.Header)
	header.Set("Accept", "application/x-json-stream")

	stream, err := r.stream(ctx, "POST", "/v3/cluster/events/log", nil, header, "application/json", &buf)
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
				return
			} else if err != nil {
				event.Component = "error"
				event.Message = err.Error()
			}

			// Don't emit keepalives
			if event.Component == "keepalive" {
				continue
			}

			ch <- event

			if event.Component == "" || event.Component == "error" {
				return
			}
		}
	}(stream, channel)

	return channel, nil
}

func (r *restclient) ClusterProcessEvents(ctx context.Context, filters api.ProcessEventFilters) (<-chan api.ProcessEvent, error) {
	var buf bytes.Buffer

	e := json.NewEncoder(&buf)
	e.Encode(filters)

	header := make(http.Header)
	header.Set("Accept", "application/x-json-stream")

	stream, err := r.stream(ctx, "POST", "/v3/cluster/events/process", nil, header, "application/json", &buf)
	if err != nil {
		return nil, err
	}

	channel := make(chan api.ProcessEvent, 128)

	go func(stream io.ReadCloser, ch chan<- api.ProcessEvent) {
		defer stream.Close()
		defer close(channel)

		decoder := json.NewDecoder(stream)

		for decoder.More() {
			var event api.ProcessEvent
			if err := decoder.Decode(&event); err == io.EOF {
				return
			} else if err != nil {
				event.Type = "error"
				event.Line = err.Error()
			}

			// Don't emit keepalives
			if event.Type == "keepalive" {
				continue
			}

			ch <- event

			if event.Type == "" || event.Type == "error" {
				return
			}
		}
	}(stream, channel)

	return channel, nil
}

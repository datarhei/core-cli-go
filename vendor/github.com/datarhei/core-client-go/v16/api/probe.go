package api

import (
	"encoding/json"
)

// ProbeIO represents a stream of a probed file
type ProbeIO struct {
	// common
	Address  string      `json:"url"`
	Format   string      `json:"format"`
	Index    uint64      `json:"index"`
	Stream   uint64      `json:"stream"`
	Language string      `json:"language"`
	Type     string      `json:"type"`
	Codec    string      `json:"codec"`
	Coder    string      `json:"coder"`
	Bitrate  json.Number `json:"bitrate_kbps" swaggertype:"number" jsonschema:"type=number"`
	Duration json.Number `json:"duration_sec"  swaggertype:"number" jsonschema:"type=number"`

	// video
	FPS    json.Number `json:"fps" swaggertype:"number" jsonschema:"type=number"`
	Pixfmt string      `json:"pix_fmt"`
	Width  uint64      `json:"width"`
	Height uint64      `json:"height"`

	// audio
	Sampling uint64 `json:"sampling_hz"`
	Layout   string `json:"layout"`
	Channels uint64 `json:"channels"`
}

// Probe represents the result of probing a file. It has a list of detected streams
// and a list of log lone from the probe process.
type Probe struct {
	Streams []ProbeIO `json:"streams"`
	Log     []string  `json:"log"`
}

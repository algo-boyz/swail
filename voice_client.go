package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// voiceQueryBase returns the voice-query-service HTTP root.
// Override with VOICE_QUERY_URL (default http://localhost:8080).
func voiceQueryBase() string {
	if u := os.Getenv("VOICE_QUERY_URL"); u != "" {
		return u
	}
	return "http://localhost:8080"
}

var httpClient = &http.Client{Timeout: 15 * time.Second}

func postJSON(ctx context.Context, path string, body any) (json.RawMessage, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	} else {
		buf.WriteString("{}")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, voiceQueryBase()+path, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("voice-query %s: %w", path, err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("voice-query %s: HTTP %d: %s", path, res.StatusCode, string(data))
	}
	return json.RawMessage(data), nil
}

func getJSON(ctx context.Context, path string) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, voiceQueryBase()+path, nil)
	if err != nil {
		return nil, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("voice-query %s: %w", path, err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("voice-query %s: HTTP %d: %s", path, res.StatusCode, string(data))
	}
	return json.RawMessage(data), nil
}

// ---- request shapes (subset of voice-query models) ----

type queryArgs struct {
	MakeModel     string `json:"make_model,omitempty"`
	VehicleClass  string `json:"vehicle_class,omitempty"`
	Color         string `json:"color,omitempty"`
	SceneL1       string `json:"scene_l1,omitempty"`
	Class         string `json:"class,omitempty"`
	TrackID       string `json:"track_id,omitempty"`
	CameraID      string `json:"camera_id,omitempty"`
	Level1        string `json:"level1,omitempty"`
	Type          string `json:"type,omitempty"`
	Contains      string `json:"contains,omitempty"`
	Source        string `json:"source,omitempty"`
	TimeWindowSec int    `json:"time_window_sec,omitempty"`
	Limit         int    `json:"limit,omitempty"`
}

// Unified row for the UI.
type QueryResultRow struct {
	Title      string  `json:"title"`
	Meta       string  `json:"meta"`
	Confidence float64 `json:"confidence"`
	ObjectID   string  `json:"object_id,omitempty"`
	Kind       string  `json:"kind"`
	Raw        string  `json:"raw,omitempty"` // original JSON for debug
}

type voiceHealth struct {
	Status     string `json:"status"`
	DemoMode   bool   `json:"demo_mode"`
	ClickHouse string `json:"clickhouse"`
	Version    string `json:"version"`
}

func pingVoiceQuery(ctx context.Context) (voiceHealth, error) {
	raw, err := getJSON(ctx, "/health")
	if err != nil {
		return voiceHealth{}, err
	}
	var h voiceHealth
	if err := json.Unmarshal(raw, &h); err != nil {
		return voiceHealth{}, err
	}
	return h, nil
}

func queryVoice(ctx context.Context, tool string, args queryArgs) ([]QueryResultRow, error) {
	path := "/v1/query/" + tool
	raw, err := postJSON(ctx, path, args)
	if err != nil {
		return nil, err
	}
	return normalizeHits(tool, raw)
}

func normalizeHits(tool string, raw json.RawMessage) ([]QueryResultRow, error) {
	switch tool {
	case "plates":
		var hits []struct {
			Plate        string  `json:"plate"`
			Confidence   float64 `json:"confidence"`
			TrackID      string  `json:"track_id"`
			CameraID     string  `json:"camera_id"`
			VehicleClass string  `json:"vehicle_class"`
			MakeModel    string  `json:"make_model"`
			Color        string  `json:"color"`
			SceneL1      string  `json:"scene_l1"`
			SceneL2      string  `json:"scene_l2"`
		}
		if err := json.Unmarshal(raw, &hits); err != nil {
			return nil, err
		}
		out := make([]QueryResultRow, 0, len(hits))
		for _, h := range hits {
			out = append(out, QueryResultRow{
				Title:      h.Plate,
				Meta:       fmt.Sprintf("%s · %s · track %s · %s/%s", nz(h.MakeModel, h.VehicleClass), h.CameraID, h.TrackID, h.SceneL1, h.SceneL2),
				Confidence: h.Confidence,
				Kind:       "plate",
			})
		}
		return out, nil

	case "objects":
		var hits []struct {
			TrackID    string  `json:"track_id"`
			Class      string  `json:"class"`
			CameraID   string  `json:"camera_id"`
			Confidence float64 `json:"confidence"`
			SceneL1    string  `json:"scene_l1"`
			SceneL2    string  `json:"scene_l2"`
		}
		if err := json.Unmarshal(raw, &hits); err != nil {
			return nil, err
		}
		out := make([]QueryResultRow, 0, len(hits))
		for _, h := range hits {
			out = append(out, QueryResultRow{
				Title:      h.Class,
				Meta:       fmt.Sprintf("track %s · %s · %s/%s", h.TrackID, h.CameraID, h.SceneL1, h.SceneL2),
				Confidence: h.Confidence,
				Kind:       "object",
			})
		}
		return out, nil

	case "scenes":
		var hits []struct {
			CameraID   string  `json:"camera_id"`
			Level1     string  `json:"level1"`
			Level2     string  `json:"level2"`
			Confidence float64 `json:"confidence"`
		}
		if err := json.Unmarshal(raw, &hits); err != nil {
			return nil, err
		}
		out := make([]QueryResultRow, 0, len(hits))
		for _, h := range hits {
			out = append(out, QueryResultRow{
				Title:      h.Level1 + " / " + h.Level2,
				Meta:       h.CameraID,
				Confidence: h.Confidence,
				Kind:       "scene",
			})
		}
		return out, nil

	case "detections":
		var hits []struct {
			ID         string  `json:"id"`
			Class      string  `json:"class"`
			CameraID   string  `json:"camera_id"`
			TrackID    string  `json:"track_id"`
			Confidence float64 `json:"confidence"`
			Metadata   string  `json:"metadata"`
		}
		if err := json.Unmarshal(raw, &hits); err != nil {
			return nil, err
		}
		out := make([]QueryResultRow, 0, len(hits))
		for _, h := range hits {
			meta := h.CameraID
			if h.TrackID != "" {
				meta += " · track " + h.TrackID
			}
			if h.Metadata != "" {
				meta += " · " + h.Metadata
			}
			out = append(out, QueryResultRow{
				Title:      h.Class,
				Meta:       meta,
				Confidence: h.Confidence,
				Kind:       "detection",
			})
		}
		return out, nil

	case "transcripts":
		var hits []struct {
			Text       string  `json:"text"`
			Speaker    string  `json:"speaker"`
			Source     string  `json:"source"`
			DeviceID   string  `json:"device_id"`
			Confidence float64 `json:"confidence"`
		}
		if err := json.Unmarshal(raw, &hits); err != nil {
			return nil, err
		}
		out := make([]QueryResultRow, 0, len(hits))
		for _, h := range hits {
			src := h.Source
			if src == "" {
				src = h.DeviceID
			}
			meta := src
			if h.Speaker != "" {
				meta = h.Speaker + " · " + src
			}
			out = append(out, QueryResultRow{
				Title:      truncate(h.Text, 80),
				Meta:       meta,
				Confidence: h.Confidence,
				Kind:       "transcript",
			})
		}
		return out, nil
	}
	return nil, fmt.Errorf("unknown tool %q", tool)
}

func nz(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

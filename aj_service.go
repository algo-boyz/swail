package main

import (
	"context"
	"time"
)

// AJService is the Go ↔ frontend bridge for AgentJetson.
type AJService struct{}

// Health reports local + voice-query + object-storage readiness.
func (s *AJService) Health() map[string]any {
	out := map[string]any{
		"ok":      true,
		"service": "swail",
		"mode":    "pr5",
		"targets": map[string]string{
			"voice_query":    voiceQueryBase(),
			"object_storage": objectStorageBase(),
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	h, err := pingVoiceQuery(ctx)
	if err != nil {
		out["voice_query"] = map[string]any{"ok": false, "error": err.Error()}
		out["ok"] = false
	} else {
		out["voice_query"] = map[string]any{
			"ok":         true,
			"status":     h.Status,
			"demo_mode":  h.DemoMode,
			"clickhouse": h.ClickHouse,
			"version":    h.Version,
		}
	}

	ost, err := pingObjectStorage(ctx)
	if err != nil {
		out["object_storage"] = map[string]any{"ok": false, "error": err.Error()}
		out["ok"] = false
	} else {
		ost["ok"] = true
		out["object_storage"] = ost
	}
	return out
}

// QueryArgs is the frontend-facing query form.
type QueryArgs struct {
	Tool          string `json:"tool"` // plates | objects | scenes | detections | transcripts
	Filter        string `json:"filter"`
	TimeWindowSec int    `json:"time_window_sec"`
	Limit         int    `json:"limit"`
	CameraID      string `json:"camera_id"`
}

// Query proxies to voice-query-service /v1/query/{tool}.
func (s *AJService) Query(args QueryArgs) ([]QueryResultRow, error) {
	tool := args.Tool
	if tool == "" {
		tool = "plates"
	}
	q := queryArgs{
		TimeWindowSec: args.TimeWindowSec,
		Limit:         args.Limit,
		CameraID:      args.CameraID,
	}
	switch tool {
	case "plates":
		q.MakeModel = args.Filter
	case "objects":
		q.Class = args.Filter
	case "scenes":
		q.Level1 = args.Filter
	case "detections":
		q.Type = args.Filter
	case "transcripts":
		q.Contains = args.Filter
	default:
		tool = "plates"
		q.MakeModel = args.Filter
	}
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	return queryVoice(ctx, tool, q)
}

// GetEvidence returns metadata + play URL for an object_id.
func (s *AJService) GetEvidence(objectID string) (EvidenceMeta, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return fetchEvidence(ctx, objectID, 3600)
}

// MockAlert is a single alert row for the UI scaffold.
type MockAlert struct {
	ID         string  `json:"id"`
	Kind       string  `json:"kind"`
	CameraID   string  `json:"camera_id"`
	Message    string  `json:"message"`
	Confidence float64 `json:"confidence"`
	Timestamp  string  `json:"timestamp"`
	ObjectID   string  `json:"object_id,omitempty"`
}

// GetRecentAlerts tries detections from voice-query; falls back to static mocks.
func (s *AJService) GetRecentAlerts() []MockAlert {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	rows, err := queryVoice(ctx, "detections", queryArgs{TimeWindowSec: 300, Limit: 20})
	if err != nil || len(rows) == 0 {
		return fallbackAlerts()
	}
	out := make([]MockAlert, 0, len(rows))
	for i, r := range rows {
		out = append(out, MockAlert{
			ID:         r.Title + "-" + string(rune('a'+i%26)),
			Kind:       r.Kind,
			CameraID:   r.Meta,
			Message:    r.Title,
			Confidence: r.Confidence,
			Timestamp:  time.Now().Format("15:04:05"),
			ObjectID:   r.ObjectID,
		})
	}
	return out
}

func fallbackAlerts() []MockAlert {
	return []MockAlert{
		{
			ID: "a-001", Kind: "plate", CameraID: "cam-front",
			Message: "7ABC123 · Mercedes C-Class", Confidence: 0.94,
			Timestamp: "14:22:01", ObjectID: "obj-rec-001",
		},
		{
			ID: "a-002", Kind: "scene", CameraID: "cam-rear",
			Message: "roadway / driving", Confidence: 0.88,
			Timestamp: "14:21:44",
		},
		{
			ID: "a-003", Kind: "officer-down", CameraID: "body-3",
			Message: "Possible officer-down (audio intent)", Confidence: 0.91,
			Timestamp: "14:18:12", ObjectID: "obj-aud-003",
		},
	}
}

// MockCamera is a camera / source tile.
type MockCamera struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Status string `json:"status"` // online | offline | degraded
	Scene  string `json:"scene"`
}

// GetCameras returns static camera list for the scaffold.
func (s *AJService) GetCameras() []MockCamera {
	return []MockCamera{
		{ID: "cam-front", Label: "Dash front", Status: "online", Scene: "roadway"},
		{ID: "cam-rear", Label: "Dash rear", Status: "online", Scene: "roadway"},
		{ID: "body-3", Label: "Body cam 3", Status: "online", Scene: "patrol"},
		{ID: "fixed-north", Label: "Fixed north gate", Status: "degraded", Scene: "perimeter"},
	}
}

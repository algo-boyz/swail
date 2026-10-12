package main

// AJService is the Go ↔ frontend bridge for AgentJetson.
// PR-1: health + mock data only. PR-2 will wire voice-query-service HTTP client.
type AJService struct{}

// Health reports local service readiness (placeholder).
func (s *AJService) Health() map[string]any {
	return map[string]any{
		"ok":      true,
		"service": "swail",
		"mode":    "scaffold",
		"targets": map[string]string{
			"voice_query":   "http://localhost:8080",
			"object_storage": "http://localhost:8081",
		},
	}
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

// GetRecentAlerts returns static mock alerts for the scaffold UI.
func (s *AJService) GetRecentAlerts() []MockAlert {
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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// objectStorageBase returns the object-storage HTTP root.
// Override with OBJECT_STORAGE_URL (default http://localhost:8081).
func objectStorageBase() string {
	if u := os.Getenv("OBJECT_STORAGE_URL"); u != "" {
		return u
	}
	return "http://localhost:8081"
}

func storageGET(ctx context.Context, path string) (json.RawMessage, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, objectStorageBase()+path, nil)
	if err != nil {
		return nil, 0, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("object-storage %s: %w", path, err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}
	if res.StatusCode >= 300 {
		return nil, res.StatusCode, fmt.Errorf("object-storage %s: HTTP %d: %s", path, res.StatusCode, string(data))
	}
	return json.RawMessage(data), res.StatusCode, nil
}

// EvidenceMeta is metadata + playback URL for the UI.
type EvidenceMeta struct {
	ObjectID    string `json:"object_id"`
	Kind        string `json:"kind"`
	Source      string `json:"source"`
	ContentType string `json:"content_type"`
	SizeBytes   int64  `json:"size_bytes"`
	Timestamp   string `json:"timestamp"`
	StorageKey  string `json:"storage_key,omitempty"`
	// PlayURL is the preferred media src (presigned, or proxy fallback).
	PlayURL   string `json:"play_url"`
	ExpiresAt string `json:"expires_at,omitempty"`
	// Mode: "presign" | "proxy"
	Mode string `json:"mode"`
}

type objectMetaEnvelope struct {
	Meta *struct {
		ObjectID       string `json:"ObjectID"`
		Kind           int    `json:"Kind"`
		Source         string `json:"Source"`
		ContentType    string `json:"ContentType"`
		SizeBytes      int64  `json:"SizeBytes"`
		StorageKey     string `json:"StorageKey"`
		Timestamp      string `json:"Timestamp"`
		CreatedAt      string `json:"CreatedAt"`
	} `json:"Meta"`
	// Also accept snake / lowercase from future JSON encodings.
	Error *struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
	} `json:"Error"`
}

type presignResponse struct {
	ObjectID  string `json:"object_id"`
	URL       string `json:"url"`
	ExpiresAt string `json:"expires_at"`
	ExpiresIn int    `json:"expires_in"`
}

var kindNames = map[int]string{
	1: "audio_transcript",
	2: "camera_recording",
	3: "frame_snapshot",
	4: "crop",
	5: "annotated_clip",
	99: "other",
}

func fetchEvidence(ctx context.Context, objectID string, expiresSec int) (EvidenceMeta, error) {
	if objectID == "" {
		return EvidenceMeta{}, fmt.Errorf("object_id required")
	}
	if expiresSec <= 0 {
		expiresSec = 3600
	}

	metaRaw, _, err := storageGET(ctx, "/v1/objects/"+objectID+"?meta=1")
	if err != nil {
		return EvidenceMeta{}, err
	}

	var env objectMetaEnvelope
	_ = json.Unmarshal(metaRaw, &env)

	out := EvidenceMeta{ObjectID: objectID, Mode: "proxy"}
	if env.Meta != nil {
		m := env.Meta
		out.ObjectID = m.ObjectID
		if out.ObjectID == "" {
			out.ObjectID = objectID
		}
		out.Kind = kindNames[m.Kind]
		if out.Kind == "" {
			out.Kind = strconv.Itoa(m.Kind)
		}
		out.Source = m.Source
		out.ContentType = m.ContentType
		out.SizeBytes = m.SizeBytes
		out.StorageKey = m.StorageKey
		out.Timestamp = m.Timestamp
		if out.Timestamp == "" {
			out.Timestamp = m.CreatedAt
		}
	}

	// Prefer presigned URL for browser playback.
	presignRaw, status, err := storageGET(ctx, "/v1/objects/"+objectID+"/presign?expires="+strconv.Itoa(expiresSec))
	if err == nil && status < 300 {
		var pr presignResponse
		if json.Unmarshal(presignRaw, &pr) == nil && pr.URL != "" && !isUnusablePresign(pr.URL) {
			out.PlayURL = pr.URL
			out.ExpiresAt = pr.ExpiresAt
			out.Mode = "presign"
			return out, nil
		}
	}

	// Fallback: stream through object-storage proxy (supports Range).
	out.PlayURL = objectStorageBase() + "/v1/objects/" + objectID
	out.Mode = "proxy"
	out.ExpiresAt = time.Now().UTC().Add(time.Duration(expiresSec) * time.Second).Format(time.RFC3339)
	return out, nil
}

func isUnusablePresign(u string) bool {
	// Filesystem backend may return empty or file:// which browsers cannot use.
	return u == "" || len(u) > 7 && u[:7] == "file://"
}

func pingObjectStorage(ctx context.Context) (map[string]any, error) {
	raw, _, err := storageGET(ctx, "/health")
	if err != nil {
		return nil, err
	}
	var h map[string]any
	if err := json.Unmarshal(raw, &h); err != nil {
		return map[string]any{"raw": string(raw)}, nil
	}
	return h, nil
}

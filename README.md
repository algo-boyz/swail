# Swail — AgentJetson Dashboard

Cross-platform front-end for [AgentJetson](https://www.agentjetson.ai/) core services
(voice-query-service, object-storage, alerts, evidence).

Built with [Wails v3](https://v3.wails.io/) + Svelte 5 + TypeScript.

## Status

**PR-5** — Evidence player streams media via object-storage presign (proxy fallback).

| Route | Purpose |
|-------|---------|
| `/` | Live overview |
| `/query` | Query console → voice-query-service |
| `/cameras` | Camera / source list (mock) |
| `/alerts` | Alert feed from detections |
| `/evidence/:id` | Media player (`presign` or `proxy`) |

## Dev

```bash
# Terminal 1 — voice-query-service
cd agentjetson/core/voice-query-service
DEMO_MODE=true go run ./cmd/server   # :8080

# Terminal 2 — object-storage
cd agentjetson/core/object-storage
DEMO_MODE=true go run ./cmd/server   # :8081

# Put a sample clip (optional)
curl -s -X POST 'http://localhost:8081/v1/objects?kind=camera_recording&source=cam-front' \
  -H 'Content-Type: video/mp4' --data-binary @sample.mp4 | jq
# → note object_id / Meta.ObjectID

# Terminal 3 — dashboard
cd swail
# optional overrides:
# export VOICE_QUERY_URL=http://localhost:8080
# export OBJECT_STORAGE_URL=http://localhost:8081
task dev
```

Open `/evidence/<object_id>` — player uses presigned URL when available, otherwise the Range-capable proxy path.

### Go bridge

| File | Role |
|------|------|
| `voice_client.go` | voice-query HTTP (`/health`, `/v1/query/*`) |
| `storage_client.go` | object-storage HTTP (`/health`, meta, `/presign`) |
| `aj_service.go` | `Health`, `Query`, `GetEvidence`, `GetRecentAlerts`, `GetCameras` |

Env: `VOICE_QUERY_URL`, `OBJECT_STORAGE_URL`.

## Roadmap

1. Live alerts / NATS / packaging
2. Terminal type console that slides in as transparant overlay (green text) displaying log messages in real-time
3. Subtitle overlays of dialogue captured from a mic associated with a cam  


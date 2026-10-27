# Swail — AgentJetson Dashboard

Cross-platform front-end for [AgentJetson](https://www.agentjetson.ai/) core services
(voice-query-service, object-storage, alerts, evidence).

Built with [Wails v3](https://v3.wails.io/) + Svelte 5 + TypeScript.

## Status

**PR-5** — Evidence player streams media via object-storage presign (proxy fallback).

**feat/embed-gods-eye-view** — `/map` embeds [God's Eye View](https://github.com/bilawalsidhu/gods-eye-view) via iframe (`?embed=1`).

| Route | Purpose |
|-------|---------|
| `/` | Live overview |
| `/query` | Query console → voice-query-service |
| `/cameras` | Camera / source list (mock) |
| `/alerts` | Alert feed from detections |
| `/evidence/:id` | Media player (`presign` or `proxy`) |
| `/map` | God's Eye View globe (iframe embed) |

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

### God's Eye View embed (`/map`)

1. Clone and run [gods-eye-view](https://github.com/bilawalsidhu/gods-eye-view) (Node 24+/26):

   ```bash
   git clone https://github.com/bilawalsidhu/gods-eye-view.git
   cd gods-eye-view && npm ci && npm run dev
   # → http://localhost:4173
   ```

2. Allow framing from Swail (GEV blocks arbitrary parents by default). For local dev:

   ```bash
   export GEV_EMBED_FRAME_ANCESTORS="*"
   # or a specific origin if you know the Wails webview origin
   npm run dev
   ```

3. In Swail open **Globe** (`/#/map`). Default base URL is `http://localhost:4173` with `?embed=1`.

   - Change URL under **Settings** on the map page (persisted in `localStorage`).
   - Or pass `#/map?gev=http://host:port`.
   - **Open in browser** launches the same embed URL outside the iframe if framing fails.

Cesium / map keys stay in the GEV process; Swail only hosts the frame.

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
4. Deeper GEV integration (share-link deep links from alerts, private detection layer, MCP bridge)

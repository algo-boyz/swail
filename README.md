# Swail — AgentJetson Dashboard

Cross-platform front-end for [AgentJetson](https://www.agentjetson.ai/) core services
(voice-query-service, object-storage, alerts, evidence).

Built with [Wails v3](https://v3.wails.io/) + Svelte 5 + TypeScript.

## Status

**PR-2** — Query console + alerts wired to `voice-query-service` over HTTP.

| Route | Purpose |
|-------|---------|
| `/` | Live overview (cameras, alerts, health) |
| `/query` | Query console → `POST /v1/query/{plates,objects,scenes,detections,transcripts}` |
| `/cameras` | Camera / source list (mock) |
| `/alerts` | Alert feed from detections (mock fallback) |
| `/evidence/:id` | Evidence player placeholder (presigned media in PR-5) |

## Dev

```bash
# Terminal 1 — core voice-query-service (demo data, no ClickHouse)
cd path/to/agentjetson/core/voice-query-service
export DEMO_MODE=true
go run ./cmd/server   # :8080

# Terminal 2 — dashboard
cd path/to/swail
export VOICE_QUERY_URL=http://localhost:8080   # optional; this is the default
task dev   # or: wails3 dev
```

### Go bridge

- `voice_client.go` — thin HTTP client (`GET /health`, `POST /v1/query/*`)
- `aj_service.go` — `Health`, `Query`, `GetRecentAlerts`, `GetCameras`
- Override base URL with `VOICE_QUERY_URL`

Frontend bindings use `Call.ByName('main.AJService.*')` (see `frontend/bindings/swail/ajservice.ts`).
Optional: `wails3 generate bindings -ts -names` to refresh.

## Roadmap

1. ~~PR-1 — scaffold + placeholder UI~~
2. ~~PR-2 — voice-query-service HTTP client + Query console~~
3. **PR-3/4** — object-storage presign + CORS (core)
4. **PR-5** — Evidence player with presigned URLs
5. Live alerts / NATS / packaging

# Swail — AgentJetson Dashboard

Cross-platform front-end for [AgentJetson](https://www.agentjetson.ai/) core services
(voice-query-service, object-storage, alerts, evidence).

Built with [Wails v3](https://v3.wails.io/) + Svelte 5 + TypeScript.

## Status

**PR-1 scaffold** — SniGo / marketplace UI removed. Placeholder shell aligned with
`agentjetson/core` voice-query-service + evidence plane.

| Route | Purpose |
|-------|---------|
| `/` | Live overview (cameras, alerts, health) |
| `/query` | Query console (plates, objects, scenes, transcripts) |
| `/cameras` | Camera / source list |
| `/alerts` | Alert feed |
| `/evidence/:id` | Evidence player (presigned media) |

## Dev

```bash
task dev          # or: wails3 dev
cd frontend && npm run dev   # frontend-only
```

Point at local core services later:

- voice-query-service → `http://localhost:8080`
- object-storage → `http://localhost:8081`

## Roadmap

1. **PR-1** (this) — scaffold + placeholder UI
2. **PR-2** — voice-query-service HTTP client + Query console live data
3. **PR-3/4** — object-storage presign + CORS (core)
4. **PR-5** — Evidence player with presigned URLs
5. Live alerts / NATS / packaging

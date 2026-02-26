# session-manager-grpc-plugin-server-go

An Extend Override app for the **session manager** written in Go. AGS calls this gRPC server with lifecycle hooks whenever game sessions or party sessions are created, updated, or deleted.

This is a template project — clone it, replace the sample logic in the service implementation, and deploy.

## Build & Test

```bash
make build                           # Build the project
go test ./...                        # Run unit tests
docker compose up --build            # Run locally with Docker
make proto                           # Regenerate proto code
```

Linting: `golangci-lint run` (config in `.golangci.yml`).

## Architecture

AGS invokes this app's gRPC methods instead of its default logic:

```
Game Client → AGS → [gRPC] → This App → Response → AGS
```

The sample implementation handles six lifecycle hooks (OnSessionCreated, OnSessionUpdated, OnSessionDeleted, OnPartyCreated, OnPartyUpdated, OnPartyDeleted) and demonstrates injecting custom attributes into session data on creation events.

### Key Files

| Path | Purpose |
|---|---|
| `main.go` | Entry point — starts gRPC server, wires interceptors and observability |
| `pkg/server/grpcserver.go` | **Service implementation** — your custom logic goes here |
| `pkg/proto/session-manager.proto` | gRPC service definition (AccelByte-provided, do not modify) |
| `pkg/pb/` | Generated code from proto (do not hand-edit) |
| `pkg/common/` | Auth interceptor, tracing, logging utilities |
| `docker-compose.yaml` | Local development setup |
| `.env.template` | Environment variable template |

## Rules

See `.agents/rules/` for coding conventions, commit standards, and proto file policies.

## Environment

Copy `.env.template` to `.env` and fill in your credentials.

| Variable | Description |
|---|---|
| `AB_BASE_URL` | AccelByte base URL (e.g. `https://test.accelbyte.io`) |
| `AB_NAMESPACE` | Target namespace |
| `AB_CLIENT_ID` | OAuth client ID |
| `AB_CLIENT_SECRET` | OAuth client secret |
| `PLUGIN_GRPC_SERVER_AUTH_ENABLED` | Enable gRPC auth (`true` by default) |

## Dependencies

- [AccelByte Go SDK](https://github.com/AccelByte/accelbyte-go-sdk) — AGS platform SDK and gRPC plugin utilities

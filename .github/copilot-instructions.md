# Copilot Instructions

See also: [`.github/gopls.instructions.md`](./gopls.instructions.md) for Go-specific workspace workflows using the gopls MCP server.

## Project Overview

`slacker-bot` is a CLI tool that creates backups of Discord guilds (servers). It fetches guild metadata, channels, members, roles, webhooks, and other entities via the Discord REST API, serializes them as Protocol Buffers, and writes them to disk—including downloading all asset URLs (avatars, icons, banners).

**Environment variables required at runtime:**
- `DISCORD_BOT_TOKEN` — Discord bot authentication token
- `DISCORD_GUILD_ID` — Snowflake ID of the guild to back up

## Commands

```bash
make run          # go run .
make test         # ginkgo run -r (all packages)
make cover        # ginkgo run --cover -r
make build        # Nix build → bin/slacker-bot
make fmt          # nix fmt
make generate     # go generate ./gen/mocks/...
make docker       # Build OCI image → bin/image.tar.gz
make up           # Load image and podman compose up --force-recreate
make import       # gomod2nix import (sync after go.mod changes)
make update       # nix flake update && go get -u ./...
```

**Run a single package's tests:**
```bash
ginkgo run ./pkg/backup
```

**Filter to a specific test by name:**
```bash
ginkgo run ./pkg/backup --label-filter "mapUsers"
ginkgo run ./pkg/backup -grep "mapUsers"
```

## Architecture

```
main.go
  └─ backup.Create(ctx, rest.Rest, snowflake.ID) → *pb.ServerBackup
       Fetches all guild data from the Discord REST API and maps
       Discord types to Protobuf message types.
  └─ backup.Write(ctx, *pb.ServerBackup, ihfs.FS)
       Writes the backup as a binary .binpb file and downloads
       all referenced asset URLs into a local directory structure.
```

- `pkg/backup/create.go` — API fetching and type mapping
- `pkg/backup/writer.go` — Protobuf serialization and asset downloading
- `gen/pb/` — Generated Protobuf Go code (do not edit manually)
- `gen/mocks/` — Generated mock for `rest.Rest` (do not edit manually)

Proto definitions live in the external `buf.build/unmango/apis` module. To regenerate Go code after proto changes, run `buf generate` (configured in `buf.gen.yaml`).

## Key Conventions

### Protobuf Builder API
Generated types use the opaque builder pattern. Always use `_builder` structs and `.Build()`:

```go
b := &pb.Guild_builder{
    Id:   ptr(g.ID.String()),
    Name: ptr(g.Name),
}
return b.Build()
```

### Discord Enum Offset
Discord enumerations start at 0, but Protobuf convention reserves 0 for `UNSPECIFIED`. Apply `+1` when converting:

```go
VerificationLevel: ptr(pb.VerificationLevel(int32(g.VerificationLevel) + 1))
```

### Optional Fields
Use `new(value)` or a `ptr()` helper to produce `*T` for optional proto fields. Return `nil` for absent values:

```go
func optID(id *snowflake.ID) *string {
    if id == nil { return nil }
    s := id.String()
    return &s
}
```

### Logging
Extract the logger from context; attach fields with `.With()`:

```go
logger := log.FromContext(ctx).With("guildId", id)
logger.Info("creating backup")
```

### Testing (Ginkgo + Gomega + uber-go/mock)
- All tests use Ginkgo v2 BDD style (`Describe` / `Context` / `It`)
- Assertions use Gomega (`Expect(...).To(...)`)
- Mocks are generated with `go.uber.org/mock/mockgen`; the REST mock is in `gen/mocks/`
- Use `DescribeTable` for parameterized cases
- Test helpers (e.g., `testRestGuild()`, `testMember()`) live alongside the tests they support

### Generated Code
Never edit files under `gen/` by hand:
- `gen/pb/` — regenerate with `buf generate`
- `gen/mocks/` — regenerate with `make generate`

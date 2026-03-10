# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build -o pomodoro-cli

# Run directly
go run main.go

# Run with options
go run main.go -timer 25 -break 5 -interval 4 -silent -notick

# Install to $GOPATH/bin
go install github.com/maiku1008/pomodoro-cli@latest
```

No tests exist in this codebase yet. No linter is configured.

## Architecture

This is a macOS-only CLI Pomodoro timer. No external dependencies — standard library only.

**Package layout:**

- `main.go` — parses CLI flags, sets up context/signal handling, builds `pomodoro.Config`, calls `pomodoro.Run`
- `internal/pomodoro/pomodoro.go` — core timer logic: runs work/break cycles, manages site blocking/unblocking, plays sounds, renders the progress bar countdown
- `internal/sound/sound.go` — embeds `.wav` files via `//go:embed`, extracts them to temp files on first use (`sync.Once`), plays via macOS `afplay`; ticking loops with overlap to avoid gaps
- `internal/hosts/hosts.go` — adds/removes a `### Pomodoro CLI ###` block from `/etc/hosts` to redirect domains to `127.0.0.1`

**Key behaviors:**
- The 3x longer break only triggers on the final interval when `-interval` is 2 or higher
- Website blocking requires `sudo`; cleanup is via `defer` so sites are always unblocked on exit or Ctrl+C
- Sound is macOS-specific (`afplay`, `afinfo`); `sound.Cleanup()` removes temp files on exit
- Context cancellation propagates from `SIGINT`/`SIGTERM` through all timers

##  Coding standards and guidelines
1. Simpler is better - do not overcomplicate - code should be readable and understandable
2. Comments only when necessary - code should be self explanatory as much as possible
3. Be concise; short README, no emojis
4. When changing something, ALWAYS ensure documentation (README.md, CLAUDE.md, etc.) reflect those changes

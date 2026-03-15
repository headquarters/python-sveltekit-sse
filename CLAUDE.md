# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a monorepo demonstrating Server-Sent Events (SSE) integration between a FastAPI backend and SvelteKit frontend.

**Structure:**
- `/python` - FastAPI SSE server
- `/go` - Go SSE server (drop-in replacement for Python, same port/endpoint)
- `/sveltekit` - SvelteKit frontend client

## Architecture

**Backend (FastAPI — `/python`):**
- Uses `sse-starlette` library with `EventSourceResponse` for proper SSE implementation
- SSE messages include: `event` type, `id` (UUID), `retry` timeout, and JSON `data` payload
- Implements client disconnect detection via `request.is_disconnected()`
- Main endpoint: `/poll` - streams random numbers with timestamps every second

**Backend (Go — `/go`):**
- Standard library only (`net/http`, `crypto/rand`) — no external dependencies
- Uses `http.Flusher` to flush SSE chunks after each write
- Client disconnect detection via `r.Context().Done()` channel
- UUID v4 generated with `crypto/rand` (no external package needed)
- Main endpoint: `/poll` - identical SSE format and behavior to the Python backend
- Drop-in replacement: same port (8000), same endpoint, same event format

**Frontend (SvelteKit):**
- Standard SvelteKit 2.x setup with TypeScript
- Svelte 5.x components
- Vitest for unit testing

## Development Commands

### Python Backend (from `/python` directory)
```sh
# Run development server with auto-reload
./run.sh
# or
poetry run uvicorn main:app --reload

# Install dependencies
poetry install

# Update dependencies after changing pyproject.toml
poetry lock && poetry install
```

### Go Backend (from `/go` directory)
```sh
# Run server
./run.sh
# or
go run .

# Build binary
go build .
```

### SvelteKit Frontend (from `/sveltekit` directory)
```sh
# Development server
npm run dev

# Type checking
npm run check

# Linting
npm run lint

# Format code
npm run format

# Run tests
npm test              # single run
npm run test:unit     # watch mode

# Build for production
npm run build
npm run preview       # preview production build
```

## Key Implementation Details

**SSE Message Format:**
The Python backend sends SSE events with this structure:
```python
{
    "event": "random_number",
    "retry": 15000,
    "data": json.dumps({"value": int, "timestamp": str}),
    "id": str(uuid.uuid4())
}
```

**SSE Best Practices:**
- Always use `EventSourceResponse` from `sse-starlette` (not raw `StreamingResponse`)
- Check `await request.is_disconnected()` in event loops to handle client disconnections
- Include `retry` timeout to control client reconnection behavior
- Use unique `id` for each event to enable event resumption

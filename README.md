# Full-Stack Calculator

A calculator with a **React + TypeScript** frontend and a **Go** REST microservice.
The frontend sends every calculation to the backend, which validates the input and returns the result as JSON.

<img src="docs/screenshot.png" alt="Calculator UI on mobile" width="280" />

**Operations:** addition, subtraction, multiplication, division, exponentiation, square root, and percentage.

---

## Project structure

```
.
├── backend/                     Go REST API (standard library only)
│   ├── cmd/server/main.go       Entry point: config, HTTP server, graceful shutdown
│   └── internal/
│       ├── calculator/          Pure domain logic, no HTTP (100% test coverage)
│       └── api/                 HTTP handlers, JSON contract, error mapping, logging
├── frontend/                    React 19 + TypeScript + Vite
│   └── src/
│       ├── api/client.ts        Typed API client and ApiError
│       ├── hooks/useCalculator  State, validation and request flow
│       ├── components/          Presentational UI
│       └── lib/                 Operation metadata, input parsing, number formatting
├── Dockerfile                   Multi-stage build that serves the full stack from one container
├── docker-compose.yml
├── Makefile                     Shortcuts for running, testing and coverage
└── .github/workflows/ci.yml     Lint, test, coverage and Docker build on every push
```

---

## Setup

### Prerequisites

- Go **1.24+**
- Node.js **22+** and npm
- Docker (optional)

### Option A: Docker (full stack, one command)

```bash
docker compose up --build
# open http://localhost:8080
```

The image builds the frontend and the backend, then runs a single Go binary on a distroless image.
That binary serves both the API and the static React build.

### Option B: Run each layer locally

**Backend** (http://localhost:8080):

```bash
cd backend
go run ./cmd/server
```

**Frontend** (http://localhost:5173), in a second terminal:

```bash
cd frontend
npm install
npm run dev
```
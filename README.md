# TrueGul

A writing training service that helps users improve their writing skills without AI assistance.

## Overview

TrueGul allows users to write content independently, detects whether the text was AI-generated, and provides quality feedback. The goal is to help users develop genuine writing skills and break dependency on AI writing tools.

### Key Features

- **AI Detection**: Analyzes submitted text and returns an AI-generation probability score (0-100%)
- **Quality Feedback**: Provides tailored feedback based on writing type (essay, cover letter)
- **Writing Management**: Create, save drafts, and track submission history
- **Rate Limiting**: Daily submission limits to encourage thoughtful writing

## Tech Stack

| Component | Technology |
|-----------|------------|
| Frontend | Next.js 15, React 19, TypeScript, Tailwind CSS, shadcn/ui |
| API Server | Go |
| ML Server | Python, FastAPI |
| Database | PostgreSQL (AWS RDS) |
| Infrastructure | AWS (ECS/Fargate, ECR, RDS) |

## Project Structure

```
TrueGul/
├── README.md           # This file
├── SPEC.md             # Project specification (Korean)
├── sprints/            # Sprint planning documents
│   ├── README.md
│   ├── sprint-0.md
│   ├── sprint-1.md
│   ├── sprint-2.md
│   ├── sprint-3.md
│   ├── sprint-4.md
│   └── sprint-5.md
├── frontend/           # Next.js frontend (TBD)
├── api-server/         # Go API server (TBD)
└── ml-server/          # Python ML server (TBD)
```

## Getting Started

### Prerequisites

- Node.js 22+
- pnpm 10+
- Go 1.25+
- Python 3.11+
- Docker & Docker Compose

### Local Development (without ML Server)

**1. Start backend services (PostgreSQL, Redis, API Server)**

```bash
# Clone the repository
git clone https://github.com/your-username/TrueGul.git
cd TrueGul

# Start services with Docker Compose (without ML server)
docker compose -f docker-compose.dev.yml up -d --build

# Check services are running
docker compose -f docker-compose.dev.yml ps
```

**2. Start frontend (in a separate terminal)**

```bash
cd frontend
pnpm install
pnpm dev
```

**3. Access the application**

| Service | URL |
|---------|-----|
| Frontend | http://localhost:3000 |
| API Server | http://localhost:8080 |
| PostgreSQL | localhost:5432 |
| Redis | localhost:6379 |

**4. Stop services**

```bash
# Stop all services
docker compose -f docker-compose.dev.yml down

# Stop and remove volumes (reset database)
docker compose -f docker-compose.dev.yml down -v
```

### Default Development Credentials

| Service | Credential |
|---------|------------|
| PostgreSQL User | `truegul` |
| PostgreSQL Password | `truegul123` |
| PostgreSQL Database | `truegul` |

### Full Stack (with ML Server)

```bash
# Requires .env file with OPENAI_API_KEY
docker compose up -d --build
```

## Roadmap

| Version | Goal | Description |
|---------|------|-------------|
| v0 | MVP | Single model, basic features, minimal logging |
| v1 | Model Separation | Separate AI detection and feedback models |
| v2 | Monitoring | ML monitoring + server monitoring |
| v3 | Typing Analysis | Typing pattern collection and analysis |

## Documentation

- [Project Specification](./SPEC.md) (Korean)
- [Sprint Plans](./sprints/README.md)

## License

TBD

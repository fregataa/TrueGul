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

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Web UI    │────▶│  API Server │────▶│  Database   │
│  (Next.js)  │     │ (Go/Fargate)│     │ (RDS)       │
└─────────────┘     └──────┬──────┘     └─────────────┘
                           │
                           ▼ (Async/Polling)
                    ┌─────────────┐     ┌─────────────┐
                    │  ML Server  │────▶│ OpenAI API  │
                    │  (Fargate)  │     │ (Feedback)  │
                    └─────────────┘     └─────────────┘
```

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

- Node.js 20+
- Go 1.21+
- Python 3.11+
- Docker & Docker Compose
- PostgreSQL 15+

### Local Development

```bash
# Clone the repository
git clone https://github.com/your-username/TrueGul.git
cd TrueGul

# Start services with Docker Compose
docker-compose up -d

# Run database migrations
# (Instructions will be added after Sprint 0)
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

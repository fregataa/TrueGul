# TrueGul

AI-powered TOPIK (Test of Proficiency in Korean) writing assessment service.

## Overview

TrueGul helps TOPIK II test takers practice and improve their essay writing skills through AI-based scoring and feedback. The service evaluates Korean essays based on official TOPIK scoring criteria and provides detailed feedback for self-directed learning.

### Key Features

- **AI Scoring**: Evaluates essays based on TOPIK official rubric (Content 20pts, Structure 15pts, Language 15pts = 50pts total)
- **Detailed Feedback**: Provides specific improvement suggestions for each scoring category
- **AI Detection**: Detects AI-generated content to encourage authentic writing practice
- **Push Notifications**: Async processing with FCM/APNs notifications when scoring completes
- **Submission History**: Track progress and review past submissions

### Target Users

- TOPIK II (Level 3-6) test preparation students
- International students seeking Korean university admission
- Professionals needing TOPIK certification for employment

## Tech Stack

| Component | Technology |
|-----------|------------|
| Mobile App | Flutter or React Native |
| API Server | Go/Gin |
| ML Server | Python/FastAPI |
| Database | PostgreSQL |
| Message Queue | Redis |
| AI/ML | LLM API (Claude/GPT) + RoBERTa (AI detection) |
| Push Notification | FCM (Android) / APNs (iOS) |
| Infrastructure | AWS (ECS/Fargate, ECR, RDS, S3) |

## Project Structure

```
TrueGul/
├── README.md              # This file
├── SPEC.md                # Service specification (Korean)
├── sprints/               # Sprint implementation plans
│   ├── README.md          # Sprint planning index
│   ├── sprint-0.md        # Planning & Setup
│   ├── sprint-1.md        # API Server
│   ├── sprint-2.md        # ML Server
│   ├── sprint-3.md        # Mobile App
│   ├── sprint-4.md        # Integration & QA
│   └── sprint-5.md        # Deployment
├── versions/              # Version roadmap documents
│   ├── v1-ocr.md          # OCR feature
│   ├── v2-data-collection.md
│   ├── v3-model-optimization.md
│   └── v4-expansion.md
├── api-server/            # Go API server
├── ml-server/             # Python ML server
└── infra/                 # Terraform infrastructure
```

## Getting Started

### Prerequisites

- Go 1.25+
- Python 3.11+
- Docker & Docker Compose
- Flutter 3.x or React Native (for mobile development)

### Local Development

**1. Start backend services**

```bash
# Clone the repository
git clone https://github.com/your-username/TrueGul.git
cd TrueGul

# Start services with Docker Compose
docker compose -f docker-compose.dev.yml up -d --build

# Check services are running
docker compose -f docker-compose.dev.yml ps
```

**2. Access services**

| Service | URL |
|---------|-----|
| API Server | http://localhost:8080 |
| ML Server | http://localhost:8000 |
| PostgreSQL | localhost:5432 |
| Redis | localhost:6379 |

**3. Stop services**

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

## Roadmap

| Version | Goal | Key Features |
|---------|------|--------------|
| **v0** | MVP | Text input, LLM scoring, AI detection, Push notifications |
| **v1** | OCR | Handwritten answer image upload with OCR |
| **v2** | Data Collection | User feedback, actual score linking, analytics |
| **v3** | Optimization | Score calibration, fine-tuning, RAG |
| **v4+** | Expansion | Additional question types (51-53), personalization |

## Documentation

- [Service Specification](./SPEC.md) (Korean)
- [Sprint Plans](./sprints/)

## License

TBD

# TSETMC Options Platform

A dockerized platform for collecting, storing, and analyzing Tehran Stock Exchange (TSETMC) option-market tick data. The system provides a modular architecture with an automated data ingestion pipeline, a high-performance Go backend with PostgreSQL storage, and a TypeScript/React frontend dashboard for monitoring and data exploration.

## Quick Start

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```
2. Build and start all services via Docker Compose:
   ```bash
   docker compose up -d --build
   ```
3. Open the frontend in your browser:
   [http://localhost:3000](http://localhost:3000)

4. Verify backend health endpoints:
   - Application health: [http://localhost:8080/api/v1/health](http://localhost:8080/api/v1/health)
   - Database health: [http://localhost:8080/api/v1/health/db](http://localhost:8080/api/v1/health/db)

## Folder Layout

```text
├── backend/       # Go REST API service, migrations, and CLI commands
├── frontend/      # React + TypeScript + Vite SPA dashboard
├── crawler/       # Python tick data crawler (added in later cards)
├── docs/          # Platform documentation and guides (cookbook in E1-1)
├── infra/         # Infrastructure configurations and deployment assets
├── .github/       # GitHub Actions CI workflows and PR templates
└── docker-compose.yml # Orchestration for local development and services
```

## Workflow

To maintain code quality and delivery standards, adhere to the following workflow rules:
- **No Direct Push to `main`**: All changes must go through feature branches and pull requests.
- **Branch Naming**: Use the format `<card-id>-<desc>` (e.g., `E6-1-repo-skeleton`).
- **Commit Messages**: Follow conventional commits: `type(scope): message` (e.g., `feat(backend): add health endpoints`).
- **Pull Requests**: Open a PR using the provided PR template, conduct self-review, and ensure CI checks are green before merging.

## AI Usage Disclosure

An AI coding agent was used for initial repository scaffolding, boilerplate generation, and configuration setup. The author thoroughly reviewed, tested, and understands all code within this repository. Core domain components including crawler routines, contract registry logic, and quality control (QC) checks are authored and verified by the platform author.

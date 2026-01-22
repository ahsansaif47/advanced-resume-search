# Advanced Resume Search

A high-performance resume search application built with Go, Weaviate (vector database), and Temporal (workflow orchestration). Uses AI-powered semantic search for finding relevant resumes.

## Architecture Overview

This project follows a microservices architecture with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Architecture                   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐     ┌──────────────┐                     │
│  │  Go Backend  │────▶│   Weaviate   │ (Vector Database)   │
│  │  (Fiber API) │     │              │                     │
│  └──────────────┘     └──────────────┘                     │
│         │                     │                             │
│         ▼                     ▼                             │
│  ┌──────────────┐     ┌──────────────┐                     │
│  │   Temporal   │     │  Text2Vec    │ (ML Model)          │
│  │  (Workflows) │     │  Reranker    │                     │
│  └──────────────┘     └──────────────┘                     │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## File Structure

### Key Files

- **`Dockerfile`**: Builds the Go application for Railway deployment (production)
- **`docker-compose.yml`**: Runs complete stack locally with all dependencies (development)
- **`Makefile`**: Build automation with swagger generation
- **`scripts/gen-swagger.sh`**: Generates API documentation
- **`.env.example`**: Template for environment variables

### Why Separate Files?

| File | Purpose | When Used | Environment |
|------|---------|-----------|-------------|
| `Dockerfile` | Builds Go binary only | Railway deployment | Production (cloud) |
| `docker-compose.yml` | Runs full stack | Local development | Development (local) |

**Important:** Do NOT integrate `docker-compose.yml` into `Dockerfile`. Each serves a different purpose:
- The `Dockerfile` creates a small, efficient image for the Go app only
- The `docker-compose.yml` orchestrates multiple services for local development
- Railway requires individual services, not a monolithic container

## Quick Start

### Prerequisites

- Go 1.25 or later
- Docker and Docker Compose
- Git

### Local Development

1. **Clone the repository:**
   ```bash
   git clone <your-repo-url>
   cd advanced-resume-search
   ```

2. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your API keys
   ```

3. **Start the infrastructure:**
   ```bash
   docker-compose up -d
   ```
   
   This starts:
   - Weaviate (vector database) on `localhost:8081`
   - Text2Vec model (embeddings) on `localhost:8080`
   - Reranker model (search reranking) on `localhost:8080`
   - Temporal (workflow engine) on `localhost:7233`
   - Temporal UI on `localhost:8080`

4. **Build and run the Go application:**
   ```bash
   # Generate swagger docs and build
   make build
   
   # Run the application
   ./out
   ```
   
   Or run directly:
   ```bash
   go run cmd/server/main.go
   ```

5. **Access the services:**
   - API: http://localhost:8080
   - Weaviate Console: http://localhost:8081
   - Temporal UI: http://localhost:8080

## Railway Deployment

### Step-by-Step Guide

#### 1. Prepare the Dockerfile

The Dockerfile is already configured for Railway deployment. It:
- Uses Alpine Linux for minimal image size (~50MB)
- Generates Swagger docs during build
- Includes all necessary build tools
- Produces a production-ready binary

#### 2. Create Railway Services

You'll need to deploy these services on Railway:

**Service 1: Go Application**
- Repository: Connect your Git repository
- Dockerfile: Automatically detected
- Environment Variables: Configure as shown below

**Service 2: Weaviate**
- Add new service → Select "Weaviate"
- Railway will provide connection details

**Service 3: Temporal** (optional)
- Add new service → Select "Temporal"
- Or use Railway's managed Temporal

#### 3. Configure Environment Variables

On Railway, set these variables for your Go application:

```bash
# Server

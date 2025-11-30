# Paper Airplane Rating System - Go gRPC Microservice

## Project Overview
A gRPC microservice for rating paper airplane designs. Supports HTTP/2 transcoding for RESTful access alongside native gRPC. MongoDB stores entities, rating categories, and ratings submissions.

**Core APIs:**
- Airplane entity CRUD (name, description, picture)
- Rating category CRUD (name, description, scale, weight)
- Rating submission (ratings per category within defined scales)
- Rating reports and weighted overall scores

## Architecture

### Data Layer
- **MongoDB collections:** `airplane-entity`, `airplane-rating-category`, `airplane-rating` (singular names)
- One document type per collection
- Connection managed in devcontainer docker-compose setup

### Service Separation
- **AirplaneService** - Manages airplane entities (basic info only, no categories)
- **CategoryService** - Manages rating categories (scale, weight, description)
- **RatingService** - Manages ratings and summaries (references categories by ID)

### Logging
- JSON structured logs to stdout/stderr (external collection)
- TODO: Verify logging structure implementation when code exists

### HTTP/2 Transcoding
- gRPC services expose both gRPC and REST endpoints
- Check `.proto` annotations for HTTP route mappings

## Development Environment

### Devcontainer Setup
**Required:** Use devcontainer (VSCode or manual docker-compose in `.devcontainers/`)
- Includes MongoDB and dev dependencies
- Minimizes local OS requirements

### Build System
- **go-task** for builds and tests (not Make)
- Run tasks before creating PRs: `task build`, `task test`
- Check `Taskfile.yml` for available tasks

## Git Workflow (Scaled Trunk-Based)

### Branching Strategy
- **Base all work on `main`** (NOT release branches)
- Feature branches: `feature/description` or forked repo branches (recommended to avoid conflicts)
- Hotfix branches: `hotfix/description`
- Release branches: `release/{major}.{minor}.x`

### Commits
- **Conventional Commits format required** (https://www.conventionalcommits.org)
- Examples: `feat: add airplane entity creation`, `fix: rating calculation overflow`
- Scope commits to specific changes

### PR Process
1. Commit all changes to feature/hotfix branch
2. Create PR to `main` (triggers CI/CD)
3. Features → next release
4. Hotfixes → cherry-picked to active release branches

**⚠️ NEVER base branches off release branches**

## Deployment
- Helm charts in `deployments/` directory
- Semantic versioning on release branches

## Code Patterns (When Implemented)

### Project Structure
```
cmd/server/              # Server entrypoint (main.go)
proto/                   # Proto definitions
  ├── airplane/v1/      # Airplane entity service
  ├── category/v1/      # Rating category service
  └── rating/v1/        # Rating service
internal/                # Private application code
  ├── pb/               # Generated proto code (DO NOT EDIT)
  ├── server/           # gRPC server setup
  ├── service/          # Business logic
  ├── repository/       # MongoDB data access
  └── middleware/       # gRPC interceptors
pkg/                    # Public libraries
scripts/                # Build scripts
deployments/            # Helm charts
```

### Proto Definitions
- Proto files in `proto/{service}/v1/` with HTTP transcoding annotations
- **Three separate services:** airplane, category, rating
- Categories are referenced by ID in ratings (not embedded)
- Generate: `task proto` (downloads googleapis automatically)
- Never edit files in `internal/pb/` - they're auto-generated

### MongoDB Operations
- Collection names: singular form (`airplane-entity`, `airplane-rating-category`, `airplane-rating`)
- **Separate collections for each resource type**
- Implement repository pattern for data access
- Use context for timeouts and cancellation

### Rating Calculations
- Ratings reference category IDs (not category names)
- Category data (scale, weight) stored separately in category collection
- Overall rating computed by joining rating values with category weights
- Validate ratings against category scale from category collection

### Testing
- Use `bufconn` for in-memory gRPC testing
- MongoDB tests use `airplane_rating_test` database
- Run: `task test`

### Code Style
- Format: `task fmt` (runs gofmt + goimports)
- Lint: `task lint` (golangci-lint with config in `.golangci.yml`)
- All checks: `task ci` (run before PRs)

## Quick Reference
- **First time setup:** `task install-tools && task proto-deps && task init`
- **Generate proto:** `task proto`
- **Run tests:** `task test`
- **Build:** `task build`
- **Run server:** `task run`
- **All CI checks:** `task ci`
- **Start devcontainer:** VSCode → Reopen in Container, or `docker-compose -f .devcontainer/docker-compose.yml up`
- **Commit format:** `type(scope): description` (conventional commits)
- **MongoDB URI:** `mongodb://mongodb:27017` (in devcontainer)

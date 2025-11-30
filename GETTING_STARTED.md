# Paper Airplane Rating System - Quick Start

## Initial Setup

### 1. Open in Devcontainer
This project requires the devcontainer to run. Open VSCode and:
- Press `F1` → "Dev Containers: Reopen in Container"
- Or use the popup notification to reopen in container

The devcontainer includes:
- Go 1.23
- Protocol Buffer compiler
- MongoDB (accessible at `mongodb://mongodb:27017`)
- All required gRPC tooling
- go-task for builds

The devcontainer will automatically run:
- task install-tools
- task proto
- task init

## Development Commands

```bash
# List all available tasks
task

# Build the server
task build

# Run tests
task test

# Run the server
task run

# Format code
task fmt

# Run linters
task lint

# Clean build artifacts
task clean

# Run all CI checks
task ci
```

## Project Structure

```
├── cmd/server/              # Server entrypoint
├── proto/                   # Protocol buffer definitions
│   ├── airplane/v1/        # Airplane entity service
│   ├── category/v1/        # Rating category service
│   └── rating/v1/          # Rating service
├── internal/               # Private application code
│   ├── pb/                 # Generated proto code (auto-generated)
│   ├── server/             # gRPC server implementation
│   ├── service/            # Business logic layer
│   └── repository/         # MongoDB data access
├── pkg/                    # Public packages
├── deployments/            # Helm charts
└── scripts/                # Build and utility scripts
```

## MongoDB Access

MongoDB is available in the devcontainer at:
- **URI:** `mongodb://mongodb:27017`
- **Database:** `airplane_rating`

Collections:
- `airplane-entity` - Airplane design entities (name, description, picture)
- `airplane-rating-category` - Rating categories (name, scale, weight)
- `airplane-rating` - Rating submissions (references category IDs)

## Proto Files

The project uses gRPC with HTTP/2 transcoding for REST support:
- **gRPC:** Port 8080
- **HTTP/REST:** Port 9090 (via grpc-gateway)

Proto files in `proto/` define three separate services:

**AirplaneService** - Airplane entity CRUD
- Create, read, update, delete airplane entities
- Airplanes store basic info (name, description, picture)

**CategoryService** - Rating category CRUD
- Create, read, update, delete rating category
- Category define scale (min/max) and weight for ratings

**RatingService** - Rating submission and reports
- Submit ratings (references category IDs)
- Get ratings for an airplane
- Get rating summaries with weighted scores

HTTP route mappings for REST access are defined in each proto file.

## Next Steps

1. Implement service layer in `internal/service/`
2. Add MongoDB repository in `internal/repository/`
3. Wire up services in `cmd/server/main.go`
4. Add tests in `*_test.go` files

## Git Workflow

- Base all work on `main` branch
- Use conventional commits: `feat:`, `fix:`, `docs:`, etc.
- Create PRs to `main` for features and hotfixes
- Run `task ci` before creating PRs

See `.github/copilot-instructions.md` for detailed development guidelines.

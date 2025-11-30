# Paper Airplane Rating System

A Go gRPC microservice for rating paper airplane designs with HTTP/2 transcoding support.

## Quick Start

See [GETTING_STARTED.md](GETTING_STARTED.md) for setup instructions.

**TL;DR:**
1. Open in devcontainer (VSCode → "Reopen in Container")
2. Run `task install-tools && task init && task proto`
3. Run `task build` and `task test`

## Project Overview
This is a sandbox project for experimenting with Go gRPC microservices.
The architecture follows modern microservice patterns with protocol buffers for a service definition per resource type.
HTTP/2 transcoding is used to allow RESTful access alongside gRPC, this is implemented with grpc-gateway (https://github.com/grpc-ecosystem/grpc-gateway).
Data is all be stored in a Mongodb database.

This project is a simple microservice that is a paper airplane rating system.
Each paper airplane is an entity which includes name, description, and picture
Each paper airplane will be rated based on rating categories which include the scale and weight of that category.

There are APIs for:
* CRUD operations on the paper airplane entity, this api allows entities to be managed and store all the data about that entity.
* CRUD operations on the category, this allows categories to be managed and adjusted.
* CRUD operations on ratings about each entity. Each rating is for rating categories, within the rating scale of that category.
* Getting reports on the summary of ratings for an entity along with getting an overall rating based on the weights on each category.

### Data
MongoDB is used to store and lookup all the persistent data.
Each document type is in its own collection.
Each collection is named for the singular verb of the document contents of that collection.

Collections:
* airplane-entity
* airplane-rating
* airplane-rating-category

### Logging
Logging uses JSON structured logs and outputs the logs on stdout and stderr. The logs will be picked up and recorded external to the service.

TODO logging structure

## Development
This project makes use of devcontainer to minimize the requirements needed on the local operating system.
The devcontainer setup uses docker compose to create the devcontainer along with dependent resources needed for development such as MongoDB.

You can use use VSCode devcontainers to start up the development environment, or you can run it manually using the docker compose file in the `.devcontainers` folder.

This project uses a standard scaled trunk based development workflow (https://trunkbaseddevelopment.com/).
You can use either feature branches or forked branches to work on features or hotfixes.
Forked branches are recommended as there is lower risk of causing git conflicts when rebasing or squashing commits.

The local build system uses go-task to build and test the project to ensure features are working as defined before a PR is requested to start the release process.

All feature and hotfix branches are based on the `main` branch of the primary git repository.
Features are merged back to the `main` branch and will become part of the next release.
Hotfixes are merged back to the `main` branch and then cherry-picked into the active release branches to start the CI/CD release process.

[!NOTE]
DO NOT base a feature or hotfix branch off of a release branch.

All commits should be scoped to the changes needed for the feature or hotfix.
Commit messages need to follow the conventional commits message structure (https://www.conventionalcommits.org/en/v1.0.0/)

### Development Workflow

```bash
# List all tasks
task

# Generate proto code
task proto

# Build
task build

# Run tests
task test

# Run locally
task run

# Format and lint
task fmt
task lint

# Run all CI checks before PR
task ci
```

See [Taskfile.yml](Taskfile.yml) for all available tasks.

### Release Workflow
Releases are tagged off of the release branches using semantic versioning. 
The release branches are named `release/{major}.{minor}.x`, where the major and minor are the semantic versions of that release.

Once a feature or hotfix seems to be ready, make sure to commit all of your changes and then create a pull request (PR) from your branch into the `main` branch in the repo.

This will trigger the CICD process to get this feature/hotfix added to the next release and in the case of hotfixes will be cherry-picked to the active release branches.

## Deployment
Deployment is handled by a helm chart in the `deployments` directory

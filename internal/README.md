# Internal Packages

This directory contains private application code that should not be imported by external projects.

## Structure

- `pb/` - Generated protocol buffer code
- `server/` - gRPC server implementation and setup
- `service/` - Business logic layer
- `repository/` - Data access layer for MongoDB
- `middleware/` - gRPC interceptors (logging, auth, etc.)

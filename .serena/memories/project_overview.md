# df-permission Project Overview

## Project Purpose
df-permission is a comprehensive permission management library for Minecraft Bedrock Edition Dragonfly servers, written in Go. It provides a robust, thread-safe, cached permission system with wildcard support, group-based permissions, and flexible storage options.

## Key Features
- **Wildcard Permission System**: Supports `*` (global) and `prefix.*` (prefix) wildcards
- **Group-Based Management**: Hierarchical permission management through groups
- **Thread-Safe Operations**: Full concurrent operation support with RWMutex
- **Performance Optimized**: TTL-based caching with automatic cleanup
- **Flexible Configuration**: Options pattern for customizable setup
- **Clean Architecture**: Layered design following domain-driven principles

## Tech Stack
- **Language**: Go 1.24+
- **Dependencies**: 
  - github.com/google/uuid v1.6.0 (UUID generation and handling)
  - Standard library only for core functionality
- **Build/Release**: GoReleaser for automated releases
- **CI/CD**: GitHub Actions for release workflow

## Project Type
This is a Go library/package (not an executable application) designed to be imported and used by other Go projects, specifically Minecraft Dragonfly servers.
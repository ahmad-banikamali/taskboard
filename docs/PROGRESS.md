# Progress

## Goal
Task board API in Go. Project 1 of the path to senior backend developer.
Hexagonal architecture. Final project of the path: bread delivery store.

## Architecture rules
- Layout: cmd/taskboard, internal/<module>/{domain,port,service,adapter/inbound,adapter/outbound}
- Modules: identity, board. Modules talk only through ports.
- Core (domain, port, service) must not import adapters or frameworks.
- Domain errors live in core. Adapters translate their own errors.
- Timestamps in UTC.
- main calls run() error. run handles graceful shutdown.
- Config comes from environment variables (internal/config).
- Tests run with -count=1.

## Done
- Milestone 1, Step 1: repo skeleton

## Current state
Empty skeleton. main prints one line.

## Next task
Milestone 1, Step 2: config package and graceful shutdown.

## Known issues
None yet.

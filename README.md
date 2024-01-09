# Redis Connection Factory
Gets a redis client in a standard and safe manner.

## Features
- Multiple attempts to connect with a backoff
- Backoff jitter to prevent thundering herd connection scenarios
- No additional external dependencies

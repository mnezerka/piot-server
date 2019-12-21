# PIOT Server

## Architecture

The startup procedure of the server initiates global context which is later
propagated to all handlers. This instance of the context is signleton and holds
configuration parameters as well as instances of all services. The idea is to
encapsulate all generic functionality and be albe to pass it to the chain of
handlers.

```
  +------------------+  +------------------+
  | Logging Service  |  | Things Service   |  .... all services
  +---------+--------+  +--------+---------+
            |                    |
            +--------------------+  Context holds instances (singletons)
            |
  +------------------+
  | Context          +-------------------+
  +------------------+                   |
                                         |
  +----------------------------------+   | Context is passed to handlers
  | CORS handler                     |   |
  | +------------------------------+ |   |
  | | AddContext handler           +-----+
  | | +--------------------------+ | |
  | | | Logging handler          | | |
  | | | +----------------------+ | | |
  | | | | Auth handler         | | | |
  | | | | +------------------+ | | | |
  | | | | | GraphQL handler  | | | | |
  | | | | +------------------+ | | | |
  | | | +----------------------+ | | |
  | | +--------------------------+ | |
  | +------------------------------+ |
  +----------------------------------+

```

## Development Environment

1. Run only mongodb docker container

   ```
    docker-compose up -d mongodb
   ```

2. Run script ``scripts/env.sh`` to get IP address of mongo container
   and set env variable for piot server

3. Run tests:

   ```
   # all tests
   go test ./...

   # tests for selected package (handler)
   go test ./handler

   # tests for selected test case (matched against regexp)
   go test --run ShortNotation

   ```

# PIOT Server

![](https://dev.azure.com/michalnezerka/PIOT/_apis/build/status/mnezerka.piot-server?branchName=master)

## Architecture

The startup procedure of the server initiates global context which is later
propagated to all handlers. This instance of the context is signleton and holds
configuration parameters as well as instances of all services. The idea is to
encapsulate all generic functionality and be albe to pass it to the chain of
handlers

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

## Development environment - minimal

1. Run only mongodb docker container::
   ```bash
   docker-compose up -d mongodb
   ```

2. Run script `scripts/env.sh` to get IP address of mongo container
   and set env variable for piot server

3. Run tests (not in parallel since shared mongodb is used):
   ```bash
   # all tests
   go test -p 1 ./...

   # tests for selected package (handler)
   go test ./handler

   # tests for selected test case (matched against regexp)
   go test --run ShortNotation
   ```
## Development environment - full stack

1. Run all containers::
   ```bash
   docker-compose up -d
   ```

2. Configure mysql - create schema (see doc/deployment.rst)

3. Build and start piot server::
   ```bash
   go build && ./piot-server --mysqldb-host localhost --mysqldb-user piot --mysqldb-password piot -l debug
   ```

Refer to `doc` folder for more information

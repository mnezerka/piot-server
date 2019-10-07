# PIOT Mosquitto Auth Server

## Development Environment

1. Run only mongodb docker container

   ```
    docker-compose up -d mongodb
   ```

2. Run script ``scripts/env.sh`` to get IP address of mongo container
   and set env variable for piot server

3. Run all tests:

   ```
   go test ./...
   ```

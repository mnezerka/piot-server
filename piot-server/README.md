# PIOT Server

## Resources

MongoDB and Golang:
https://vkt.sh/go-mongodb-driver-cookbook/

Testing Tips:
- https://medium.com/@benbjohnson/structuring-tests-in-go-46ddee7a25c

## Development Environment

1. Run only mongodb docker container

   ```
    docker-compose up -d mongodb
   ```

2. Run script ``scripts/env.sh`` to get IP address of mongo container
   and set env variable for piot server

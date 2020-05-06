Development of PIOT
===================

Development environment
-----------------------

Start MongoDB in daemon mode::

    docker-compose up

Configure environment for piot-server::

    cd piot-server
    source scripts/env.sh

Start piot-server::

    go build && ./piot-server -l DEBUG --mqtt-uri tcp://localhost:1883 --mqtt-user piot --mqtt-password piot

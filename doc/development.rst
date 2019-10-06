Development of PIOT
===================

Development environment
-----------------------

Start MongoDB in daemon mode::

    docker-compose up -d mongodb

Temporarily switch networking for mosquitto mqtt broker to host network by
uncommenting following line in ``mqtt`` section of ``docker-compose.yml``::

    #network_mode: host

Temporarily change hostname of mosquitto authentication to localhost by 
modifying ./mosquitto/conf/conf.d/go-auth.conf file::

    auth_opt_http_host           localhost
    #auth_opt_http_host           piot-server

Start mqtt in no-daemon mode to see logs::

    docker-compose up mqtt


Configure environment for piot-server::

    cd piot-server
    source scripts/env.sh

Start piot-server::

    go build && ./piot-server -l DEBUG --mqtt-user piot --mqtt-password piot

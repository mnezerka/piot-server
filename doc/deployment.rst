Deployment of PIOT
==================

NGINX Password
--------------

Password for HTTP basic authentication which is used for all services that are
designed as unprotected (e.g. Prometheus) could be generated as::

    openssl passwd -apr1

The generated string needs to be pasted to ``./nginx/nginx-passwd`` file on
production server

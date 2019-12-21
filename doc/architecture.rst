PIOT Architecture
=================

The central component of the whole system is MQTT broker mosquitto together
with PIOT Server, which provides GraphQL API for system administration.

Schema::

    +---------+              +----------+    +------------+    +-------------+    +-----------+
    |  MQTT   | TLS          |   MQTT   |    | Prometehus |    | Prometheus  |    |  Grafana  |
    | Device  +--------------+  Broker  +----+  Exporter  +----+             +----+           |
    +---------+       +------+----+-----+    +------------+    +-------------+    +---+--+----+
                      |           |                                                   |  |
                      |           |       Grafana provisioning Alert Notifications    |  |
                      |           |   +-----------------------------------------------+  |
                      |           |   |                                                  |
    +----------+      |      +----+---+-+                                          +-----+----+
    |   PIOT   +------+      |   PIOT   | GraphQL API                              |   NGINX  |  HTTPS (TLS)
    |  Adapter +-------------+  Server  +------------------------------------------+          +---------------
    +-----+----+             +----------+                                          +----------+
          |                       |  |      +-----------+                                |
          |Internet               |  |      |  MongoDB  |                                |
          |GPRS                   |  +------+           |                                |
          |                       |         +-----------+                                |
    +-----+----+             +----------+                                                |
    |   PIOT   |             |   PIOT   | Web Application                                |
    |  Device  |             |  Manager +------------------------------------------------+
    +----------+             +----------+


:MQTT Broker:
    The broker is heart and central point for distributing MQTT messages between
    clients (devices, sensors, applications). It can handle up to thousands of
    concurrently connected MQTT clients. It is responsible for receiving all messages,
    filtering, determining who is subscribed to each message, and sending the message
    to these subscribed clients.

:PIOT Server:
    Server responsible for management of all things (devices), users and organizations
    as well as for collecting data from MQTT Broker and persisting those data. It
    provides HTTP API (GraphQL) for almost all actions it provides.

:PIOT Manager:
    Web application connected to PIOT Server API that allows users to interactively
    manage IOT data (e.g. create organization, assign discovered things).

:MQTT Device:
    Any device capable to communicate over MQTT protocol. It can be single purpose
    device (e.g. thermo meter), another MQTT broker or even business application
    that consumes specific MQTT messages and generates new based on logic
    (e.g. https://nodered.org)

:PIOT Adapter:
    Simple server that listens for messages in proprieatry protocol used by PIOT
    Devices, translates all incoming traffic to MQTT and publishes it on MQTT Broker.

:PIOT Device:
    Device, typically based on ESP8266 or AT-family chip, that sends sensor reading
    messages in proprietary format.

:Prometheus Exporter:
    Standalone server that acts as MQTT client with single purpose - read statistics
    from MQTT server (e.g. number of connections) and publish it in a form of
    prometehus metrics

:Prometheus:
    Time series database that collects data from various sources and provides it in
    a time bound format. It is easily consumable by e.g. Grafana platform.

:Grafana:
    Analytics and monitoring solution that provides nice dashboards and charts.
    All content could be made multi tenant by spliting it into Organizations.

:NGINX:
    Main role of this server is to act as reverse proxy for other servers that
    don't support TLS or even any kind of protection (e.g. prometheus). It also
    makes it easier to organize servers running in docker containers into single
    structure of endpoint behind single domain name.

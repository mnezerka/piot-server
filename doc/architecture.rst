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


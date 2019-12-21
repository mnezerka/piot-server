MQTT
====

The MQTT protocol is just transporting and distribution mechanism, but doesn't
tell anything about structure of topics and format of values. It is up to each
application or platform to decide on details. PIOT is trying to put focus on:

* Each topic starts with organization and thing name
* Each published value is scalar (number, string)


Structure of Topics
-------------------

The basic idea of topic hierarchy is following

``Organization/ThingName/[ThingDetail]``

where ``ThingDetail`` refers subtree specific for given thing. Those are
generic attributes common for all things (if possible)

:Availability:
    Topic ``Organization/ThingName/available`` is dedicated to availability
    of given thing. The possible values are ``yes`` and ``no``

:Temperature:

    * Topic ``Organization/ThingName/temperature`` for value
    * Topic ``Organization/ThingName/temperature/unit`` for unit

:Humidity:

    * Topic ``Organization/ThingName/humidity`` for value
    * Topic ``Organization/ThingName/humidity/unit`` for unit

:Pressure:

    * Topic ``Organization/ThingName/pressure`` for value
    * Topic ``Organization/ThingName/pressure/unit`` for unit

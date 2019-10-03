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


PIOT Protocol and MQTT
----------------------

Data are received over HTTP in POST request. Server replies with HTTP
status codes, where 200 OK means the readings were accepted. Format of the raw
http body is json that matches following schema::

    {
        "device": "DEVICE_ID",
        "ip": "device ip address (optional)",
        "wifi-ssid": "device wifi ssid (optional)",
        "wifi-strength": deivice wifi strength (optional),
        "time": unix_timestamp (optional),
        "readings": [
            {
                "address": "unique sensor address",
                "t": temperature (optional),
                "h": humidity (optional);
                "p": pressure (optional);
            }
        ]
    }

The only mandatory field on global level is ``device``. If there
is at least one entry in ``readings``, it has to contain ``address`` attribute.

The minimal http chunk could look like which is kind of hart beat
notification saying that device is alive::

    {
        "device": "Device123",
    }

This is example of minimal http chunk with sensor data (temperature)::

    {
        "device": "Device123",
        "readings": [
            {
                "address": "SensorXYZ",
                "t": 23,
            }
        ]
    }


Encryption
..........

Server accepts both unecrypted and encrypted data. The only supported
algorithm is AES 128bit ECB with PKCS7 padding. Free implementation
in C is available here: https://github.com/mnezerka/blue-aes

Translation to MQTT
...................

Each device have to be assigned to organization. Let's assume we have
device identified as *CHIP23* assigned to organization *TestOrg*. The 
following packet::

    {
        "device": "CHIP23",
        "ip": "192.168.0.55",
        "wifi-ssid": "TestOrgWifi",
        "wifi-strength": 57,
        "readings": [
            {
                "address": "67543465",
                "t": 25,
                "h": 45,
                "p": 567;
            }
        ]
    }

will be translated into following sequence of MQTT publish calls::

    PUBLISH /TestOrg/CHIP23/available -> "yes"
    PUBLISH /TestOrg/CHIP23/net/ip-> "192.168.0.55"
    PUBLISH /TestOrg/CHIP23/net/wifi/ssid -> ""TestOrgWifi"
    PUBLISH /TestOrg/CHIP23/net/wifi/strength": "57"
    PUBLISH /TestOrg/67543465/available -> "yes"
    PUBLISH /TestOrg/67543465/temperature -> "25"
    PUBLISH /TestOrg/67543465/temperature/unit -> "C"
    PUBLISH /TestOrg/67543465/numidity -> "45"
    PUBLISH /TestOrg/67543465/numidity/unit -> "%"
    PUBLISH /TestOrg/67543465/pressure -> "567"
    PUBLISH /TestOrg/67543465/pressure/unit -> "Pa"

package main

import (
	"errors"
	"fmt"
	"piot-server/config"
	"strconv"
	"time"

	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// topic name used for publishing sensor readings
const PIOT_MEASUREMENT_TOPIC = "value"

type PiotDevices struct {
	log    *logging.Logger
	things *Things
	mqtt   IMqtt
	params *config.Parameters
	cache  map[string]time.Time
}

// constructor
func NewPiotDevices(logger *logging.Logger, things *Things, mqtt IMqtt, params *config.Parameters) *PiotDevices {
	p := PiotDevices{log: logger, things: things, mqtt: mqtt, params: params}
	p.cache = make(map[string]time.Time)
	return &p
}

func (p *PiotDevices) ProcessPacket(packet PiotDevicePacket) error {
	p.log.Debugf("Process PIOT device packet: %v", packet)

	// handle short notation of attributes (assign short to long attributes)
	if len(packet.DeviceShort) > 0 {
		packet.Device = packet.DeviceShort
	}
	if len(packet.ReadingsShort) > 0 {
		packet.Readings = packet.ReadingsShort
	}

	// DOS Protection
	// allow to process data from this packet only if it didn't come too close
	// to previous packet from the same device (treshold in seconds is defined)
	if lastSeen, ok := p.cache[packet.Device]; ok {
		delta := time.Now().Sub(lastSeen)
		p.log.Debugf("Time cache holds entry for device: %s, time diff is %f seconds", packet.Device, delta.Seconds())

		if delta <= p.params.DOSInterval {
			return errors.New("Exceeded dos protection treshold")
		}
	}

	// name of the device cannot be empty
	if packet.Device == "" {
		return errors.New("Device name cannot be empty")
	}

	// store device name to cache together with date it was seen
	p.cache[packet.Device] = time.Now()

	// get instance of Things service and look for the device (chip),
	// register it if it doesn't exist
	thing, err := p.things.FindPiot(packet.Device)
	if err != nil {
		// register device
		thing, err = p.things.RegisterPiot(packet.Device, THING_TYPE_DEVICE)
		if err != nil {
			return err
		}

		// configure availability topic
		if err := p.things.SetAvailabilityTopic(thing.Id, "available"); err != nil {
			return err
		}
		if err := p.things.SetAvailabilityYesNo(thing.Id, "yes", "no"); err != nil {
			return err
		}
	}

	// if thing is assigned to org
	if thing.OrgId != primitive.NilObjectID {
		// try to push data to mqtt
		if err = p.processDevice(thing, packet); err != nil {
			return err
		}
	} else {
		p.log.Debugf("Ignoring processing of data for thing <%s> that is not assigned to any organization", thing.Name)
	}

	// look for sensors and register those that doesn't exist
	for _, reading := range packet.Readings {

		// handle short notation of address attribute
		if len(reading.AddressShort) > 0 {
			reading.Address = reading.AddressShort
		}

		if reading.Temperature != nil {
			class := THING_CLASS_TEMPERATURE
			if err = p.processReading(class, thing, reading); err != nil {
				p.log.Debugf("Failed to process reading data for thing <%s>", thing.Name)
			}
		}

		if reading.Humidity != nil {
			class := THING_CLASS_HUMIDITY
			if err = p.processReading(class, thing, reading); err != nil {
				p.log.Debugf("Failed to process humidity reading data for thing <%s>", thing.Name)
			}
		}

		if reading.Pressure != nil {
			class := THING_CLASS_PRESSURE
			if err = p.processReading(class, thing, reading); err != nil {
				p.log.Debugf("Failed to process pressure reading data for thing <%s>", thing.Name)
			}
		}
	}

	return nil
}

func (p *PiotDevices) processDevice(thing *Thing, packet PiotDevicePacket) error {

	p.log.Debugf("Process PIOT device data: %v", packet)

	// dont' push anything if device is disabled
	if !thing.Enabled {
		return nil
	}

	// update avalibility channel
	err := p.mqtt.PushThingData(thing, TOPIC_AVAILABLE, VALUE_YES)
	if err != nil {
		return err
	}

	if packet.Ip != nil {
		err := p.mqtt.PushThingData(thing, TOPIC_IP, *packet.Ip)
		if err != nil {
			return err
		}
	}

	if packet.WifiSSID != nil {
		err := p.mqtt.PushThingData(thing, TOPIC_WIFI_SSID, *packet.WifiSSID)
		if err != nil {
			return err
		}
	}

	if packet.WifiStrength != nil {
		if err := p.mqtt.PushThingData(thing, TOPIC_WIFI_STRENGTH, fmt.Sprintf("%f", *packet.WifiStrength)); err != nil {
			return err
		}
	}

	return nil
}

func (p *PiotDevices) processReading(class string, thing *Thing, reading PiotSensorReading) error {
	p.log.Debugf("Process PIOT device reading data of class \"%s\": %v", class, reading)

	var address string = reading.Address
	var value string
	var unit string

	// determine address from class
	// this is necessary to have separate things for all sensor measurements
	switch class {
	case THING_CLASS_TEMPERATURE:
		address = "T" + address
		unit = "C"
		if reading.Temperature != nil {
			value = strconv.FormatFloat(float64(*reading.Temperature), 'f', -1, 32)
		}
	case THING_CLASS_HUMIDITY:
		address = "H" + address
		unit = "%"
		if reading.Humidity != nil {
			value = strconv.FormatFloat(float64(*reading.Humidity), 'f', -1, 32)
		}
	case THING_CLASS_PRESSURE:
		address = "P" + address
		unit = "mPa"
		if reading.Pressure != nil {
			value = strconv.FormatFloat(float64(*reading.Pressure), 'f', -1, 32)
		}
	}

	// look for thing representing sensor
	sensor_thing, err := p.things.Find(address)

	// if thing not found
	if err != nil {

		// register register device
		sensor_thing, err = p.things.RegisterPiot(address, THING_TYPE_SENSOR)
		if err != nil {
			return err
		}

		// register topics for measurements (if presetn)
		if p.things.SetSensorMeasurementTopic(sensor_thing.Id, PIOT_MEASUREMENT_TOPIC); err != nil {
			return err
		}

		// set proper device class according to received measurement type
		if err := p.things.SetSensorClass(sensor_thing.Id, class); err != nil {
			return err
		}
	}

	// update parent thing (this can happen any time since sensor can be
	// re-connected to another device
	if sensor_thing.ParentId != thing.Id {
		err = p.things.SetParent(sensor_thing.Id, thing.Id)
		if err != nil {
			return err
		}
	}

	// if thing is not assigned to org
	if sensor_thing.OrgId == primitive.NilObjectID {
		p.log.Debugf("Ignoring processing of data for thing <%s> that is not assigned to any organization", sensor_thing.Name)

		// stop processing here
		return nil
	}

	// dont' push anything if device is disabled
	if !thing.Enabled {
		return nil
	}

	// update avalibility channel
	err = p.mqtt.PushThingData(sensor_thing, TOPIC_AVAILABLE, VALUE_YES)
	if err != nil {
		return err
	}

	if value != "" {

		if err := p.mqtt.PushThingData(sensor_thing, PIOT_MEASUREMENT_TOPIC, value); err != nil {
			return err
		}
		if err := p.mqtt.PushThingData(sensor_thing, fmt.Sprintf("%s/%s", PIOT_MEASUREMENT_TOPIC, TOPIC_UNIT), unit); err != nil {
			return err
		}
	} else {
		p.log.Warningf("Processing unkonwn sensor reading data <%v>", reading)
	}

	return nil
}

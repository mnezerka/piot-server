package main

/* Note: checking of res.ModifiedCount could be tricky since it is

 attribute is updated only in the first call

if res.ModifiedCount == 0 {
    return fmt.Errorf("Thing <%s> not found", id.Hex())
}
*/

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Things struct {
	Db  *mongo.Database
	Log *logging.Logger
}

func NewThings(db *mongo.Database, log *logging.Logger) *Things {
	things := &Things{Db: db, Log: log}
	return things
}

func (t *Things) Get(id primitive.ObjectID) (*Thing, error) {
	t.Log.Debugf("Get thing: %s", id.Hex())

	var thing Thing

	collection := t.Db.Collection("things")
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&thing)
	if err != nil {
		t.Log.Warningf("Things.Get failed for id <%s> (%v)", id.Hex(), err)
		return nil, err
	}

	return &thing, nil
}

func (t *Things) GetFiltered(filter interface{}) ([]*Thing, error) {

	collection := t.Db.Collection("things")

	ctx := context.TODO()

	var result []*Thing

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		t.Log.Errorf("GQL: error : %v", err)
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		// To decode into a struct, use cursor.Decode()
		thing := Thing{}
		err := cur.Decode(&thing)
		if err != nil {
			t.Log.Errorf("GQL: error : %v", err)
			return nil, err
		}
		result = append(result, &thing)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (t *Things) Find(name string) (*Thing, error) {
	t.Log.Debugf("Finding thing by name <%s>", name)

	var thing Thing

	// try to find thing in DB by its name
	err := t.Db.Collection("things").FindOne(context.TODO(), bson.M{"name": name}).Decode(&thing)
	if err != nil {
		return nil, errors.New("Thing not found")
	}

	return &thing, nil
}

func (t *Things) FindPiot(id string) (*Thing, error) {
	t.Log.Debugf("Finding piot thing by id <%s>", id)

	var thing Thing

	// try to find thing in DB by its name
	err := t.Db.Collection("things").FindOne(context.TODO(), bson.M{"piot_id": id}).Decode(&thing)
	if err != nil {
		return nil, errors.New("Thing not found")
	}

	return &thing, nil
}

func (t *Things) RegisterPiot(id string, deviceType string) (*Thing, error) {
	t.Log.Debugf("Registering new piot thing: %s of type %s", id, deviceType)
	// check if string of same name already exists
	_, err := t.FindPiot(id)
	if err == nil {
		return nil, fmt.Errorf("piot Thing identified by %s already exists", id)
	}

	// thing does not exist -> create new one
	var thing Thing
	thing.Name = id
	thing.PiotId = id
	thing.Type = deviceType
	thing.Created = int32(time.Now().Unix())
	thing.LastSeen = int32(time.Now().Unix())

	res, err := t.Db.Collection("things").InsertOne(context.TODO(), thing)
	if err != nil {
		t.Log.Errorf("Thing %s cannot be stored (%v)", id, err)
		return nil, errors.New("error while storing new thing")
	}

	thing.Id = res.InsertedID.(primitive.ObjectID)

	return &thing, nil
}

func (t *Things) SetParent(id primitive.ObjectID, id_parent primitive.ObjectID) error {
	t.Log.Debugf("Setting thing <%v>, setting parent to <%s>", id.Hex(), id_parent.Hex())

	_, err := t.Get(id)
	if err != nil {
		t.Log.Errorf("Thing %s not found", id.Hex())
		return errors.New("child thing not found when setting new parent")
	}

	_, err = t.Get(id_parent)
	if err != nil {
		t.Log.Errorf("Thing %s not found", id_parent.Hex())
		return errors.New("Parent thing not found when setting new parent for thing")
	}

	_, err = t.Db.Collection("things").UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"parent_id": id_parent}})
	if err != nil {
		t.Log.Errorf("Thing %s cannot be updated (%v)", id.Hex(), err)
		return errors.New("error while updating thing parent")
	}

	return nil
}

func (t *Things) SetAvailabilityTopic(id primitive.ObjectID, topic string) error {
	t.Log.Debugf("Setting thing <%s>, setting avalibility topic to <%s>", id.Hex(), topic)

	_, err := t.Db.Collection("things").UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"availability_topic": topic}})
	if err != nil {
		t.Log.Errorf("Thing %s cannot be updated (%v)", id.Hex(), err)
		return errors.New("error while updating thing attributes")
	}

	return nil
}

func (t *Things) SetAvailabilityYesNo(id primitive.ObjectID, yes, no string) error {
	t.Log.Debugf("Setting thing <%s>, setting avalibility topic values to <%s> and <%s>", id.Hex(), yes, no)

	_, err := t.Db.Collection("things").UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"availability_yes": yes, "availability_no": no}})
	if err != nil {
		t.Log.Errorf("Thing %s cannot be updated (%v)", id.Hex(), err)
		return errors.New("error while updating thing attributes")
	}

	return nil
}

func (t *Things) SetTelemetry(id primitive.ObjectID, telemetry string) error {
	t.Log.Debugf("Setting thing <%s> telemetry", id.Hex())

	_, err := t.Db.Collection("things").UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"telemetry": telemetry}})
	if err != nil {
		t.Log.Errorf("Thing %s cannot be updated (%v)", id.Hex(), err)
		return errors.New("error while updating thing attributes")
	}

	return nil
}

func (t *Things) SetLocationMqttTopic(id primitive.ObjectID, topic string) error {
	t.Log.Debugf("Setting thing <%s>, setting location topic to <%s>", id.Hex(), topic)

	_, err := t.Db.Collection("things").UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"loc_mqtt_topic": topic}})
	if err != nil {
		t.Log.Errorf("Thing %s cannot be updated (%v)", id.Hex(), err)
		return errors.New("error while updating thing attributes")
	}

	return nil
}

func (t *Things) SetLocationMqttValues(id primitive.ObjectID, lat, lng, sat, ts string) error {
	t.Log.Debugf("Setting thing <%s>, setting location mqtt params topic to <%s>, <%s>, <%s>, <%s>", id.Hex(), lat, lng, sat, ts)

	params := bson.M{
		"loc_mqtt_lat_value": lat,
		"loc_mqtt_lng_value": lng,
		"loc_mqtt_sat_value": sat,
		"loc_mqtt_ts_value":  ts,
	}
	_, err := t.Db.Collection("things").UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": params})
	if err != nil {
		t.Log.Errorf("Thing %s cannot be updated (%v)", id.Hex(), err)
		return errors.New("error while updating thing attributes")
	}

	return nil
}

func (t *Things) SetLocation(id primitive.ObjectID, lat, lng float64, sat, ts int32) error {
	t.Log.Debugf("Setting thing <%s> location", id.Hex())

	_, err := t.Db.Collection("things").UpdateOne(
		context.TODO(),
		bson.M{
			"_id": id,
			"$or": []interface{}{
				bson.M{
					"loc_ts": bson.M{
						"$exists": false,
					},
				},
				bson.M{
					"loc_ts": bson.M{
						"$lte": ts,
					},
				},
			},
		},
		bson.M{
			"$set": bson.M{
				"loc_lat": lat,
				"loc_lng": lng,
				"loc_sat": sat,
				"loc_ts":  ts,
			},
		},
	)

	if err != nil {
		t.Log.Errorf("Thing %s cannot be updated (%v)", id.Hex(), err)
		return errors.New("error while updating thing attributes")
	}

	return nil
}

func (t *Things) SetSensorMeasurementTopic(id primitive.ObjectID, topic string) error {
	t.Log.Debugf("Setting thing <%s> sensor measurement topic to <%s>", id.Hex(), topic)

	_, err := t.Db.Collection("things").UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"sensor.measurement_topic": topic}})
	if err != nil {
		t.Log.Errorf("Thing %s cannot be updated (%v)", id.Hex(), err)
		return errors.New("error while updating thing attributes")
	}

	return nil
}

func (t *Things) SetSensorClass(id primitive.ObjectID, class string) error {
	t.Log.Debugf("Setting thing <%s> sensor class to <%s>", id.Hex(), class)

	_, err := t.Db.Collection("things").UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"sensor.class": class}})
	if err != nil {
		t.Log.Errorf("Thing %s cannot be updated (%v)", id.Hex(), err)
		return errors.New("error while updating thing attributes")
	}

	return nil
}

func (t *Things) SetSensorValue(id primitive.ObjectID, value string) error {
	t.Log.Debugf("Setting thing <%s> sensor value to <%s>", id, value)

	_, err := t.Db.Collection("things").UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"sensor.value": value}})
	if err != nil {
		t.Log.Errorf("Thing %s cannot be updated (%v)", id.Hex(), err)
		return errors.New("error while updating thing attributes")
	}

	return nil
}

func (t *Things) SetSwitchState(id primitive.ObjectID, value bool) error {
	t.Log.Debugf("Setting thing <%s> switch value to <%v>", id, value)

	_, err := t.Db.Collection("things").UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"switch.state": value}})
	if err != nil {
		t.Log.Errorf("Thing %s cannot be updated (%v)", id.Hex(), err)
		return errors.New("error while updating thing attributes")
	}

	return nil
}

func (t *Things) TouchThing(id primitive.ObjectID) error {
	t.Log.Debugf("Touch thing <%s>", id.Hex())

	_, err := t.Db.Collection("things").UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"last_seen": int32(time.Now().Unix())}})
	if err != nil {
		e := fmt.Errorf("thing <%s> cannot be touched (%v)", id.Hex(), err)
		t.Log.Errorf(e.Error())
		return e
	}

	return nil
}

func (t *Things) SetAlarm(id primitive.ObjectID, active bool) error {
	t.Log.Debugf("Setting thing <%s> alarm to %v", id.Hex(), active)

	// try to find thing to be updated
	var thing Thing
	collection := t.Db.Collection("things")
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&thing)
	if err != nil {
		return errors.New("thing does not exist")
	}

	if thing.AlarmActive == active {
		// no need to set same value again
		return nil
	}

	_, err = t.Db.Collection("things").UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"alarm_active":    active,
				"alarm_activated": int32(time.Now().Unix()),
			},
		},
	)

	if err != nil {
		t.Log.Errorf("Thing %s alarm cannot be set (%v)", id.Hex(), err)
		return errors.New("error while setting thing alarm")
	}

	return nil
}

func (t *Things) Delete(id primitive.ObjectID) error {

	t.Log.Debugf("Deleting thing <%s>", id.Hex())

	collection := t.Db.Collection("things")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		t.Log.Errorf("Cannot delete thing %s (%v)", id.Hex(), err)
		return errors.New("error while deleting thing")
	}

	t.Log.Debugf("Thing %s deleted", id.Hex())
	return nil
}

func (t *Things) SetBatteryLevel(id primitive.ObjectID, level int32) error {
	t.Log.Debugf("Setting thing <%s> battery level to <%d>", id.Hex(), level)

	params := bson.M{
		"battery_level": level,
	}
	_, err := t.Db.Collection("things").UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": params})
	if err != nil {
		t.Log.Errorf("Thing %s cannot be updated (%v)", id.Hex(), err)
		return errors.New("error while updating thing attributes")
	}

	return nil
}

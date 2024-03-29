package main

import (
	"errors"
	//"fmt"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type thingUpdateInput struct {
	Id                    graphql.ID
	PiotId                *string
	Name                  *string
	Type                  *string
	Description           *string
	Alias                 *string
	Enabled               *bool
	LastSeenInterval      *int32
	Voltage               *float64
	OrgId                 *graphql.ID
	AvailabilityTopic     *string
	TelemetryTopic        *string
	StoreInfluxDb         *bool
	StoreMysqlDb          *bool
	StoreMysqlDbInterval  *int32
	LocationLat           *float64
	LocationLng           *float64
	LocationTracking      *bool
	LocationMqttTopic     *string
	LocationMqttLatValue  *string
	LocationMqttLngValue  *string
	LocationMqttSatValue  *string
	LocationMqttTsValue   *string
	BatteryLevelTracking  *bool
	BatteryMqttTopic      *string
	BatteryMqttLevelValue *string
}

type thingSensorDataUpdateInput struct {
	Id               graphql.ID
	Class            *string
	MeasurementTopic *string
	MeasurementValue *string
}

type thingSwitchDataUpdateInput struct {
	Id           graphql.ID
	Class        *string
	StateTopic   *string
	StateOn      *string
	StateOff     *string
	CommandTopic *string
	CommandOn    *string
	CommandOff   *string
}

type ThingResolver struct {
	log    *logging.Logger
	orgs   *Orgs
	things *Things
	users  *Users
	db     *mongo.Database
	t      *Thing
}

func (r *ThingResolver) Id() graphql.ID {
	return graphql.ID(r.t.Id.Hex())
}

func (r *ThingResolver) PiotId() string {
	return r.t.PiotId
}

func (r *ThingResolver) Name() string {
	return r.t.Name
}

func (r *ThingResolver) Description() string {
	return r.t.Description
}

func (r *ThingResolver) Alias() string {
	return r.t.Alias
}

func (r *ThingResolver) Type() string {
	return r.t.Type
}

func (r *ThingResolver) Enabled() bool {
	return r.t.Enabled
}

func (r *ThingResolver) Created() int32 {
	return r.t.Created
}

func (r *ThingResolver) LastSeen() int32 {
	return r.t.LastSeen
}

func (r *ThingResolver) Voltage() float64 {
	return r.t.Voltage
}

func (r *ThingResolver) LastSeenInterval() int32 {
	return r.t.LastSeenInterval
}

func (r *ThingResolver) Org() *OrgResolver {

	r.log.Debugf("GQL: Fetching org for thing: %s", r.t.Id.Hex())

	if r.t.OrgId != primitive.NilObjectID {

		org, err := r.orgs.Get(r.t.OrgId)
		if err != nil {
			r.log.Errorf("GQL: Fetching org %v for thing %v failed", r.t.OrgId, r.t.Id)
		} else {
			return &OrgResolver{r.log, r.db, r.users, org}
		}
	}

	return nil
}

func (r *ThingResolver) Parent() *ThingResolver {

	r.log.Debugf("GQL: Fetching parent for thing: %s", r.t.Id.Hex())

	if r.t.ParentId != primitive.NilObjectID {

		parentThing, err := r.things.Get(r.t.ParentId)
		if err != nil {
			r.log.Errorf("GQL: Fetching parent %v for thing %v failed", r.t.ParentId, r.t.Id)
		} else {
			return &ThingResolver{r.log, r.orgs, r.things, r.users, r.db, parentThing}
		}
	}

	return nil
}

func (r *ThingResolver) AvailabilityTopic() string {
	return r.t.AvailabilityTopic
}

func (r *ThingResolver) AvailabilityYes() string {
	return r.t.AvailabilityYes
}

func (r *ThingResolver) AvailabilityNo() string {
	return r.t.AvailabilityNo
}

func (r *ThingResolver) TelemetryTopic() string {
	return r.t.TelemetryTopic
}

func (r *ThingResolver) Telemetry() string {
	return r.t.Telemetry
}

func (r *ThingResolver) StoreInfluxDb() bool {
	return r.t.StoreInfluxDb
}

func (r *ThingResolver) StoreMysqlDb() bool {
	return r.t.StoreMysqlDb
}

func (r *ThingResolver) StoreMysqlDbInterval() int32 {
	return r.t.StoreMysqlDbInterval
}

func (r *ThingResolver) LocationLat() float64 {
	return r.t.LocationLatitude
}

func (r *ThingResolver) LocationLng() float64 {
	return r.t.LocationLongitude
}

func (r *ThingResolver) LocationSat() int32 {
	return r.t.LocationSatelites
}

func (r *ThingResolver) LocationTs() int32 {
	return r.t.LocationTs
}

func (r *ThingResolver) LocationTracking() bool {
	return r.t.LocationTracking
}

func (r *ThingResolver) LocationMqttTopic() string {
	return r.t.LocationMqttTopic
}

func (r *ThingResolver) LocationMqttLatValue() string {
	return r.t.LocationMqttLatValue
}

func (r *ThingResolver) LocationMqttLngValue() string {
	return r.t.LocationMqttLngValue
}

func (r *ThingResolver) LocationMqttSatValue() string {
	return r.t.LocationMqttSatValue
}

func (r *ThingResolver) LocationMqttTsValue() string {
	return r.t.LocationMqttTsValue
}

func (r *ThingResolver) AlarmActive() bool {
	return r.t.AlarmActive
}

func (r *ThingResolver) AlarmActivated() int32 {
	return r.t.AlarmActivated
}

func (r *ThingResolver) BatteryLevel() int32 {
	return r.t.BatteryLevel
}

func (r *ThingResolver) BatteryLevelTracking() bool {
	return r.t.BatteryLevelTracking
}

func (r *ThingResolver) BatteryMqttTopic() string {
	return r.t.BatteryMqttTopic
}

func (r *ThingResolver) BatteryMqttLevelValue() string {
	return r.t.BatteryMqttLevelValue
}

func (r *ThingResolver) Sensor() *SensorResolver {

	if r.t.Type == THING_TYPE_SENSOR {
		return &SensorResolver{r.log, r.t}
	}

	return nil
}

func (r *ThingResolver) Switch() *SwitchResolver {

	if r.t.Type == THING_TYPE_SWITCH {
		return &SwitchResolver{r.log, r.t}
	}

	return nil
}

/////////////// Sensor Data Resolver

type SensorResolver struct {
	log *logging.Logger
	t   *Thing
}

func (r *SensorResolver) MeasurementTopic() string {

	return r.t.Sensor.MeasurementTopic
}

func (r *SensorResolver) MeasurementValue() string {
	return r.t.Sensor.MeasurementValue
}

func (r *SensorResolver) Value() string {
	return r.t.Sensor.Value
}

func (r *SensorResolver) Unit() string {
	return r.t.Sensor.Unit
}

func (r *SensorResolver) Class() string {
	return r.t.Sensor.Class
}

/////////////// Switch Data Resolver

type SwitchResolver struct {
	log *logging.Logger
	t   *Thing
}

func (r *SwitchResolver) State() bool {
	return r.t.Switch.State
}

func (r *SwitchResolver) StateTopic() string {
	return r.t.Switch.StateTopic
}

func (r *SwitchResolver) StateOn() string {
	return r.t.Switch.StateOn
}

func (r *SwitchResolver) StateOff() string {
	return r.t.Switch.StateOff
}

func (r *SwitchResolver) CommandTopic() string {
	return r.t.Switch.CommandTopic
}

func (r *SwitchResolver) CommandOn() string {
	return r.t.Switch.CommandOn
}

func (r *SwitchResolver) CommandOff() string {
	return r.t.Switch.CommandOff
}

/////////////// Resolver

type ThingFilter struct {
	Name         *string
	NameContains *string
}

type ThingSort struct {
	Field string
	Order string
}

func (r *Resolver) Thing(args struct{ Id graphql.ID }) (*ThingResolver, error) {

	r.log.Debugf("GQL: Fetch thing: %v", args.Id)

	// create ObjectID from string
	id, err := primitive.ObjectIDFromHex(string(args.Id))
	if err != nil {
		r.log.Errorf("Graphql error : %v", err)
		return nil, errors.New("cannot decode ID")
	}

	thing := Thing{}

	collection := r.db.Collection("things")
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&thing)
	if err != nil {
		r.log.Errorf("Graphql error : %v", err)
		return nil, err
	}

	r.log.Debugf("GQL: Retrieved thing %v", thing)
	return &ThingResolver{r.log, r.orgs, r.things, r.users, r.db, &thing}, nil
}

func (r *Resolver) Things(ctx context.Context, args struct {
	Sort   *ThingSort
	Filter *ThingFilter
	All    *bool
}) ([]*ThingResolver, error) {

	// authorization checks
	profileValue := ctx.Value("profile")
	if profileValue == nil {
		r.log.Errorf("GQL: Missing user profile")
		return nil, errors.New("missing user profile")
	}
	profile := profileValue.(*UserProfile)
	r.log.Debugf("arg.all: %v", args.All)
	r.log.Debugf("arg.sort: %v", args.Sort)
	r.log.Debugf("arg.filter %v", args.Filter)
	r.log.Debugf("ctx.is admin: %v", profile.IsAdmin)

	all := false

	// if caller provided all atribute
	if args.All != nil {

		// and its value is true
		if *args.All {

			// check if caller is authorized to get all things
			if !profile.IsAdmin {
				r.log.Errorf("GQL: No authorization to request all things")
				return nil, errors.New("no authorization to request all things")

			} else {
				all = true
			}
		}
	}

	filter := bson.M{}

	if !all {
		if profile.OrgId.IsZero() {
			r.log.Errorf("GQL: No organization assigned")
			return nil, errors.New("no organization assigned")
		}

		// prepare filter
		filter["org_id"] = profile.OrgId
	}

	if args.Filter != nil {
		if args.Filter.Name != nil {
			filter["name"] = *args.Filter.Name
		}

		if args.Filter.NameContains != nil {
			filter["name"] = bson.M{"$regex": *args.Filter.NameContains}
		}
	}

	// prepare sorting
	opts := options.Find().SetSort(bson.M{"created": -1})
	if args.Sort != nil {
		order := 1
		if args.Sort.Order == "desc" {
			order = -1
		}
		opts.SetSort(bson.M{args.Sort.Field: order})
	}

	collection := r.db.Collection("things")

	cur, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		r.log.Errorf("GQL: error : %v", err)
		return nil, err
	}
	defer cur.Close(context.TODO())

	var result []*ThingResolver

	for cur.Next(context.TODO()) {
		// To decode into a struct, use cursor.Decode()
		thing := Thing{}
		err := cur.Decode(&thing)
		if err != nil {
			r.log.Errorf("GQL: error : %v", err)
			return nil, err
		}
		result = append(result, &ThingResolver{r.log, r.orgs, r.things, r.users, r.db, &thing})
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Resolver) CreateThing(args *struct {
	Name string
	Type string
}) (*ThingResolver, error) {

	thing, err := NewThing(r.db, r.log, args.Name, args.Type)
	if err != nil {
		return nil, err
	}

	return &ThingResolver{r.log, r.orgs, r.things, r.users, r.db, thing}, nil
}

func (r *Resolver) UpdateThing(args struct{ Thing thingUpdateInput }) (*ThingResolver, error) {

	r.log.Debugf("Updating thing %s", args.Thing.Id)

	// create ObjectID from string
	id, err := primitive.ObjectIDFromHex(string(args.Thing.Id))
	if err != nil {
		return nil, err
	}

	// try to find thing to be updated
	var thing Thing
	collection := r.db.Collection("things")
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&thing)
	if err != nil {
		return nil, errors.New("Thing does not exist")
	}

	// try to find similar thing matching new name
	if args.Thing.Name != nil {
		var similarThing Thing
		err := collection.FindOne(context.TODO(), bson.M{"$and": []bson.M{{"name": args.Thing.Name}, {"_id": bson.M{"$ne": id}}}}).Decode(&similarThing)
		if err == nil {
			return nil, errors.New("Thing of such name already exists")
		}
	}

	// thing exists -> update it
	updateFields := bson.M{}
	if args.Thing.PiotId != nil {
		updateFields["piot_id"] = *args.Thing.PiotId
	}
	if args.Thing.Name != nil {
		updateFields["name"] = *args.Thing.Name
	}
	if args.Thing.Type != nil {
		updateFields["type"] = *args.Thing.Type
	}
	if args.Thing.Description != nil {
		updateFields["description"] = *args.Thing.Description
	}
	if args.Thing.Alias != nil {
		updateFields["alias"] = *args.Thing.Alias
	}
	if args.Thing.Enabled != nil {
		updateFields["enabled"] = *args.Thing.Enabled
	}
	if args.Thing.LastSeenInterval != nil {
		updateFields["last_seen_interval"] = *args.Thing.LastSeenInterval
	}
	if args.Thing.Voltage != nil {
		updateFields["voltage"] = *args.Thing.Voltage
	}
	if args.Thing.AvailabilityTopic != nil {
		updateFields["availability_topic"] = *args.Thing.AvailabilityTopic
	}
	if args.Thing.TelemetryTopic != nil {
		updateFields["telemetry_topic"] = *args.Thing.TelemetryTopic
	}
	if args.Thing.StoreInfluxDb != nil {
		updateFields["store_influxdb"] = *args.Thing.StoreInfluxDb
	}
	if args.Thing.StoreMysqlDb != nil {
		updateFields["store_mysqldb"] = *args.Thing.StoreMysqlDb
	}
	if args.Thing.StoreMysqlDbInterval != nil {
		updateFields["store_mysqldb_interval"] = *args.Thing.StoreMysqlDbInterval
	}
	if args.Thing.LocationMqttTopic != nil {
		updateFields["loc_mqtt_topic"] = *args.Thing.LocationMqttTopic
	}
	if args.Thing.LocationMqttLatValue != nil {
		updateFields["loc_mqtt_lat_value"] = *args.Thing.LocationMqttLatValue
	}
	if args.Thing.LocationMqttLngValue != nil {
		updateFields["loc_mqtt_lng_value"] = *args.Thing.LocationMqttLngValue
	}
	if args.Thing.LocationMqttSatValue != nil {
		updateFields["loc_mqtt_sat_value"] = *args.Thing.LocationMqttSatValue
	}
	if args.Thing.LocationMqttTsValue != nil {
		updateFields["loc_mqtt_ts_value"] = *args.Thing.LocationMqttTsValue
	}
	if args.Thing.LocationLat != nil {
		updateFields["loc_lat"] = *args.Thing.LocationLat
	}
	if args.Thing.LocationLng != nil {
		updateFields["loc_lng"] = *args.Thing.LocationLng
	}
	if args.Thing.LocationTracking != nil {
		updateFields["loc_tracking"] = *args.Thing.LocationTracking
	}
	if args.Thing.BatteryLevelTracking != nil {
		updateFields["battery_level_tracking"] = *args.Thing.BatteryLevelTracking
	}
	if args.Thing.BatteryMqttTopic != nil {
		updateFields["battery_mqtt_topic"] = *args.Thing.BatteryMqttTopic
	}
	if args.Thing.BatteryMqttLevelValue != nil {
		updateFields["battery_mqtt_level_value"] = *args.Thing.BatteryMqttLevelValue
	}

	if args.Thing.OrgId != nil {
		// create ObjectID from string
		orgId, err := primitive.ObjectIDFromHex(string(*args.Thing.OrgId))
		if err != nil {
			return nil, err
		}
		updateFields["org_id"] = orgId
	}
	update := bson.M{"$set": updateFields}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		r.log.Errorf("Updating thing failed %v", err)
		return nil, errors.New("error while updating thing")
	}

	// read thing
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&thing)
	if err != nil {
		return nil, errors.New("cannot fetch thing data")
	}

	r.log.Debugf("Thing updated %v", thing)
	return &ThingResolver{r.log, r.orgs, r.things, r.users, r.db, &thing}, nil
}

func (r *Resolver) UpdateThingSensorData(args struct{ Data thingSensorDataUpdateInput }) (*ThingResolver, error) {

	r.log.Debugf("Updating thing %s sensor data", args.Data.Id)

	// create ObjectID from string
	id, err := primitive.ObjectIDFromHex(string(args.Data.Id))
	if err != nil {
		return nil, err
	}

	// try to find thing to be updated
	var thing Thing
	collection := r.db.Collection("things")
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&thing)
	if err != nil {
		return nil, errors.New("Thing does not exist")
	}

	// thing exists -> update it
	updateFields := bson.M{}
	if args.Data.Class != nil {
		updateFields["sensor.class"] = *args.Data.Class
	}
	if args.Data.MeasurementTopic != nil {
		updateFields["sensor.measurement_topic"] = *args.Data.MeasurementTopic
	}
	if args.Data.MeasurementValue != nil {
		updateFields["sensor.measurement_value"] = *args.Data.MeasurementValue
	}
	update := bson.M{"$set": updateFields}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		r.log.Errorf("Updating thing failed %v", err)
		return nil, errors.New("error while updating thing")
	}

	// read thing
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&thing)
	if err != nil {
		return nil, errors.New("cannot fetch thing data")
	}

	r.log.Debugf("Thing sensor data updated %v", thing)
	return &ThingResolver{r.log, r.orgs, r.things, r.users, r.db, &thing}, nil
}

func (r *Resolver) UpdateThingSwitchData(args struct{ Data thingSwitchDataUpdateInput }) (*ThingResolver, error) {

	r.log.Debugf("Updating thing %s switch data", args.Data.Id)

	// create ObjectID from string
	id, err := primitive.ObjectIDFromHex(string(args.Data.Id))
	if err != nil {
		return nil, err
	}

	// try to find thing to be updated
	var thing Thing
	collection := r.db.Collection("things")
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&thing)
	if err != nil {
		return nil, errors.New("Thing does not exist")
	}

	// thing exists -> update it
	updateFields := bson.M{}
	if args.Data.StateTopic != nil {
		updateFields["switch.state_topic"] = *args.Data.StateTopic
	}
	if args.Data.StateOn != nil {
		updateFields["switch.state_on"] = *args.Data.StateOn
	}
	if args.Data.StateOff != nil {
		updateFields["switch.state_off"] = *args.Data.StateOff
	}
	if args.Data.CommandTopic != nil {
		updateFields["switch.command_topic"] = *args.Data.CommandTopic
	}
	if args.Data.CommandOn != nil {
		updateFields["switch.command_on"] = *args.Data.CommandOn
	}
	if args.Data.CommandOff != nil {
		updateFields["switch.command_off"] = *args.Data.CommandOff
	}
	update := bson.M{"$set": updateFields}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		r.log.Errorf("Updating thing failed %v", err)
		return nil, errors.New("error while updating thing")
	}

	// read thing
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&thing)
	if err != nil {
		return nil, errors.New("cannot fetch thing data")
	}

	r.log.Debugf("Thing switch data updated and refetched %v", thing)
	return &ThingResolver{r.log, r.orgs, r.things, r.users, r.db, &thing}, nil
}

func (r *Resolver) SetThingAlarm(args *struct {
	Id     graphql.ID
	Active bool
}) (*bool, error) {

	r.log.Debugf("Updating thing %s alarm to %v", args.Id, args.Active)

	// create ObjectID from string
	id, err := primitive.ObjectIDFromHex(string(args.Id))
	if err != nil {
		return nil, err
	}

	// set alarm
	err = r.things.SetAlarm(id, args.Active)

	if err != nil {
		r.log.Errorf("Setting thing alarm failed %v", err)
		return nil, err
	}

	r.log.Debugf("Thing alarm updated")
	return &args.Active, nil
}

func (r *Resolver) DeleteThing(args *struct{ Id graphql.ID }) (*bool, error) {

	r.log.Debugf("Delete thing %s", args.Id)

	// create ObjectID from string
	id, err := primitive.ObjectIDFromHex(string(args.Id))
	if err != nil {
		return nil, err
	}

	err = r.things.Delete(id)

	if err != nil {
		r.log.Errorf("Delete thing failed %v", err)
		return nil, err
	}

	r.log.Debugf("Thing deleted")
	return nil, nil
}

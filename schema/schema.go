package schema

func GetRootSchema() string {

    return `
        schema {
            query: Query
            mutation: Mutation
        }

        type Query {
            userProfile(): UserProfile
            user(id: ID!): User
            users(): [User]!
            orgs(): [Org]!
            org(id: ID!): Org
            things(): [Thing]!
            thing(id: ID!): Thing
        }

        type Mutation {
            createOrg(name: String!, description: String!): Org
            updateOrg(org: OrgUpdate!): Org
            removeOrgUser(orgId: ID!, userId: ID!): Boolean

            createUser(user: UserCreate!): User
            updateUser(user: UserUpdate!): User
            addOrgUser(orgId: ID!, userId: ID!): Boolean

            createThing(name: String!, type: String!): Thing
            updateThing(thing: ThingUpdate!): Thing
            updateThingSensorData(data: ThingSensorDataUpdate!): Thing
            updateThingSwitchData(data: ThingSwitchDataUpdate!): Thing
        }

        type User {
            id: ID!
            email: String!
            password: String!
            created: Int!
            orgs: [Org!]!
        }

        type UserProfile {
            email: String!
        }

        type Org {
            id: ID!
            name: String!
            description: String!
            created: Int!
            users: [User!]!
            influxdb: String!
            influxdb_username: String!
            influxdb_password: String!
            mysqldb: String!
            mysqldb_username: String!
            mysqldb_password: String!
            mqtt_username: String!
            mqtt_password: String!
        }

        type SensorData {
            value: String!
            unit: String!
            class: String!
            measurement_topic: String!
            measurement_value: String!
        }

        type SwitchData {
            state: Boolean!
            state_topic: String!
            state_on: String!
            state_off: String!
            command_topic: String!
            command_on: String!
            command_off: String!
        }

        type Thing {
            id: ID!
            piot_id: String!
            name: String!
            description: String!
            alias: String!
            type: String!
            enabled: Boolean!
            created: Int!
            last_seen: Int!
            last_seen_interval: Int!
            org: Org
            parent: Thing
            availability_topic: String!
            availability_yes: String!
            availability_no: String!
            telemetry_topic: String!
            telemetry: String!
            store_influxdb: Boolean!
            store_mysqldb: Boolean!
            store_mysqldb_interval: Int!
            sensor: SensorData
            switch: SwitchData
            location_lat: Float!
            location_lng: Float!
            location_sat: Int!
            location_ts: Int!
            location_mqtt_topic: String!
            location_mqtt_lat_value: String!
            location_mqtt_lng_value: String!
            location_mqtt_sat_value: String!
            location_mqtt_ts_value: String!
        }

        input UserUpdate {
            id: ID!
            email: String
        }

        input UserCreate {
            email: String!
            orgId: ID
        }

        input ThingUpdate {
            id: ID!
            piotId: String
            name: String
            type: String
            description: String
            alias: String
            orgId: ID
            enabled: Boolean
            last_seen_interval: Int
            availability_topic: String
            telemetry_topic: String
            store_influxdb: Boolean
            store_mysqldb: Boolean
            store_mysqldb_interval: Int
            location_lat: Float
            location_lng: Float
            location_mqtt_topic: String
            location_mqtt_lat_value: String
            location_mqtt_lng_value: String
            location_mqtt_sat_value: String
            location_mqtt_ts_value: String
        }

        input ThingSensorDataUpdate {
            id: ID!
            class: String
            measurement_topic: String
            measurement_value: String
        }

        input ThingSwitchDataUpdate {
            id: ID!
            state_topic: String
            state_on: String
            state_off: String
            command_topic: String
            command_on: String
            command_off: String
        }

        input OrgUpdate {
            id: ID!
            name: String
            description: String
            influxdb: String
            influxdb_username: String
            influxdb_password: String
            mysqldb: String
            mysqldb_username: String
            mysqldb_password: String
            mqtt_username: String
            mqtt_password: String
        }
    `
}

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
        }

        type SensorData {
            value: String!
            unit: String!
            class: String!
            measurement_topic: String!
            measurement_value: String!
            store_influxdb: Boolean!
        }

        type SwitchData {
            state_topic: String!
            state_on: String!
            state_off: String!
            command_topic: String!
            command_on: String!
            command_off: String!
            store_influxdb: Boolean!
        }

        type Thing {
            id: ID!
            name: String!
            alias: String!
            type: String!
            enabled: Boolean!
            created: Int!
            last_seen: Int!
            org: Org
            parent: Thing
            availability_topic: String!
            availability_yes: String!
            availability_no: String!
            sensor: SensorData
            switch: SwitchData
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
            name: String
            alias: String
            orgId: ID
            enabled: Boolean
        }

        input ThingSensorDataUpdate {
            id: ID!
            class: String
            store_influxdb: Boolean
            measurement_topic: String
            measurement_value: String
        }

        input OrgUpdate {
            id: ID!
            name: String
            description: String
            influxdb: String
            influxdb_username: String
            influxdb_password: String
        }
    `
}

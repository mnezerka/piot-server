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
            updateOrg(id: ID!, name: String, description: String): Org
            removeOrgUser(orgId: ID!, userId: ID!): Boolean

            createUser(user: UserCreate!): User
            updateUser(user: UserUpdate!): User
            addOrgUser(orgId: ID!, userId: ID!): Boolean

            updateThing(thing: ThingUpdate!): Thing
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
        }

        type SensorData {
            class: String!
            measurement_topic: String!
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
            availability_topic: String!
            availability_yes: String!
            availability_no: String!
            sensor: SensorData
        }

        input UserUpdate {
            id: ID!,
            email: String
        }

        input UserCreate {
            email: String!
            orgId: ID
        }

        input ThingUpdate {
            id: ID!,
            name: String
            alias: String,
            orgId: ID
        }
    `
}

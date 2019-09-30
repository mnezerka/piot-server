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
            createUser(user: UserCreate!): User
            updateUser(user: UserUpdate!): User
            addOrgUser(orgId: ID!, userId: ID!): Boolean
            removeOrgUser(orgId: ID!, userId: ID!): Boolean
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

        type Thing {
            id: ID!
            name: String!
            type: String!
            enabled: Boolean!
            created: Int!
            org: Org
        }

        input UserUpdate {
            id: ID!,
            email: String
            orgId: ID
        }

        input UserCreate {
            email: String!
            orgId: ID
        }
    `
}

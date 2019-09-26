package schema

func GetRootSchema() string {

    return `
        schema {
            query: Query
            mutation: Mutation
        }

        type Query {
            userProfile(): UserProfile
            user(email: String!): User
            users(): [User]!
            orgs(): [Org]!
            org(id: String!): Org
            things(): [Thing]!
            thing(ID: String!): Thing
        }

        type Mutation {
            createOrg(name: String!, description: String!): Org
            updateOrg(id: ID!, name: String, description: String): Org
            createUser(user: UserCreate!): User
            updateUser(user: UserUpdate!): User
        }

        type User {
            id: ID!
            email: String!
            password: String!
            created: Int!
            org: Org
        }

        type UserProfile {
            email: String!
        }

        type Org {
            id: ID!
            name: String!
            description: String!
            created: Int!
        }

        type Thing {
            name: String!
            type: String!
            available: Boolean!
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

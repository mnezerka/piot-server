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
            customers(): [Customer]!
            customer(id: String!): Customer
            things(): [Thing]!
            thing(ID: String!): Thing
        }

        type Mutation {
            createCustomer(name: String!, description: String!): Customer
            updateCustomer(id: ID!, name: String, description: String): Customer
            createUser(email: String!): User
            updateUser(id: ID!, email: String): User
        }

        type User {
            id: ID!
            email: String!
            password: String!
            created: Int!
            customer: Customer
        }

        type UserProfile {
            email: String!
        }

        type Customer {
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
            customer: Customer
        }
    `
}

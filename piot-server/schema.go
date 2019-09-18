package main

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
            devices(): [Device]!
            device(ID: String!): Device
        }

        type Mutation {
            createCustomer(name: String!, description: String!): Customer
        }

        type User {
            email: String!
            password: String!
            created: Int!
        }

        type UserProfile {
            email: String!
        }

        type Customer {
            name: String!
            description: String!
            created: Int!
        }

        type Device {
            name: String!
            type: String!
            available: Boolean!
            created: Int!
            customer: Customer
        }
    `
}

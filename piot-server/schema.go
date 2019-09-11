package main

func GetRootSchema() string {

    return `
        schema {
            query: Query
            mutation: Mutation
        }

        type Query {
            user(email: String!): User
            users(): [User]!
        }

        type User {
            email: String!
            password: String
        }
    `
}

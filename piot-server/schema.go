package main

func GetRootSchema() string {

    return `
        schema {
            query: Query
        }

        type Query {
            user(email: String!): User
            users(): [User]!
        }

        type User {
            email: String!
            password: String!
            created: Int!
        }
    `
}

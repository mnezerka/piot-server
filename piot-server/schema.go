package main

func GetRootSchema() string {

    return `
        schema {
            query: Query
        }

        type Query {
            userProfile(): UserProfile
            user(email: String!): User
            users(): [User]!
        }

        type User {
            email: String!
            password: String!
            created: Int!
        }

        type UserProfile {
            email: String!
        }
    `
}

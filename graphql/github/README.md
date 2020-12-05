# GitHub GraphQL API client

## Goal

Get a list of all of my repositories from GitHub, alongside and/or sorted by creation date.

## From the Explorer

```GraphQL
{
  repositoryOwner(login: "jlucktay") {
    login
    repositories(first: 100, isFork: false, orderBy: {field: CREATED_AT, direction: ASC}) {
      edges {
        node {
          createdAt
          name
        }
      }
    }
  }
}
```

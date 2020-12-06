# GitHub GraphQL API client

## Goal

Get a list of all of repositories owned by me from GitHub, alongside and/or sorted by creation date.

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

## TODO

### Doing

### Done

- get forks and not-forks as two separate queries
  - in Terraform, these would be maintained as two separate resources, one volatile and one less so
- Pagination (starter limit is 100 and we're almost there already)

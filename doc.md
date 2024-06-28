# Base URI: `http://localhost:8080/graphql`

## Create user:

```graphql
mutation {
  createUser(name: "John Doe") {
    id
    name
  }
}
```

## Get user:

```graphql
query {
  user(id: 1) {
    id
    name
  }
}
```


## Update user:

```graphql
mutation {
  updateUser(id: 1, name: "Jane Doe") {
    id
    name
  }
}
```

## Delete user:

```graphql
mutation {
  deleteUser(id: 1)
}
```


## Sample curl request

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"query": "query { user(id: 1) { id name } }"}' \
  http://localhost:8080/graphql
```

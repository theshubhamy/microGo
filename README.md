# microGo

## Project Structure

The project consists of the following main components:

- Account Service
- Catalog Service
- Order Service
- GraphQL API Gateway

Each service has its own database:

- Account and Order services use PostgreSQL
- Catalog service uses Elasticsearch

## Getting Started

1. Clone the repository:

   ```
   git clone <repository-url>
   cd <project-directory>
   ```

2. Start the services using Docker Compose:

   ```
   docker-compose up -d --build
   ```

3. Access the GraphQL playground at `http://localhost:8000/playground`

Steps for grpc file generation -

1. wget https://github.com/protocolbuffers/protobuf/releases/download/v23.0/protoc-23.0-linux-x86_64.zip
2. unzip protoc-23.0-linux-x86_64.zip -d protoc
3. sudo mv protoc/bin/protoc /usr/local/bin/
4. go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
5. go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
6. echo $PATH
7. export PATH="$PATH:$(go env GOPATH)/bin"
8. source ~/.bashrc
9. create the pb folder in your project
10. add this to account.proto - option go_package = "./pb";
11. finally run this command - protoc --go_out=./pb --go-grpc_out=./pb account.proto

## GraphQL API Usage

The GraphQL API provides a unified interface to interact with all the microservices.

### Query Accounts

```graphql
query {
  accounts {
    id
    name
  }
}
```

### Create an Account

```graphql
mutation {
  createAccount(account: { name: "New Account" }) {
    id
    name
  }
}
```

### Query Products

```graphql
query {
  products {
    id
    name
    price
  }
}
```

### Create a Product

```graphql
mutation {
  createProduct(
    product: { name: "New Product", description: "A new product", price: 19.99 }
  ) {
    id
    name
    price
  }
}
```

### Create an Order

```graphql
mutation {
  createOrder(
    order: {
      accountId: "account_id"
      products: [{ id: "product_id", quantity: 2 }]
    }
  ) {
    id
    totalPrice
    products {
      name
      quantity
    }
  }
}
```

### Query Account with Orders

```graphql
query {
  accounts(id: "account_id") {
    name
    orders {
      id
      createdAt
      totalPrice
      products {
        name
        quantity
        price
      }
    }
  }
}
```

## Advanced Queries

### Pagination and Filtering

```graphql
query {
  products(pagination: { skip: 0, take: 5 }, query: "search_term") {
    id
    name
    description
    price
  }
}
```

### Calculate Total Spent by an Account

```graphql
query {
  accounts(id: "account_id") {
    name
    orders {
      totalPrice
    }
  }
}
```

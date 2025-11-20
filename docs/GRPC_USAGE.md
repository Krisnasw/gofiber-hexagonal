# gRPC API Usage

This boilerplate now supports both REST and gRPC APIs with centralized response handling.

## Starting the Server

To start both REST and gRPC servers, simply run:

```bash
go run cmd/main.go
```

The servers will start on:
- REST API: http://localhost:4001
- gRPC API: localhost:4002

## gRPC API Endpoints

### UserService

The UserService provides the following RPC methods:

1. **GetUser(GetUserRequest) returns (GetUserResponse)**
   - Retrieves a user by ID
   - Request: `GetUserRequest` with `id` field
   - Response: `GetUserResponse` with centralized response structure

2. **CreateUser(CreateUserRequest) returns (CreateUserResponse)**
   - Creates a new user
   - Request: `CreateUserRequest` with `name` and `email` fields
   - Response: `CreateUserResponse` with centralized response structure

## Response Structure

Both REST and gRPC APIs use a centralized response structure:

```protobuf
message GetUserResponse {
  bool error = 1;
  int32 code = 2;
  string message = 3;
  User data = 4;
}
```

This matches the REST API response structure:
```json
{
  "error": false,
  "code": 200,
  "message": "User retrieved successfully",
  "data": {
    "id": "1",
    "name": "John Doe",
    "email": "john.doe@example.com"
  }
}
```

## Example Client Usage

Here's how to use the gRPC API in a Go client:

```go
// Connect to the gRPC server
conn, err := grpc.Dial("localhost:4002", grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}
defer conn.Close()

// Create a client
client := v1.NewUserServiceClient(conn)

// Create a user
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

createResp, err := client.CreateUser(ctx, &v1.CreateUserRequest{
    Name:  "John Doe",
    Email: "john.doe@example.com",
})
if err != nil {
    log.Printf("Failed to create user: %v", err)
    return
}

fmt.Printf("Created user: %+v\n", createResp)

// Get the user
getResp, err := client.GetUser(ctx, &v1.GetUserRequest{
    Id: createResp.GetData().GetId(),
})
if err != nil {
    log.Printf("Failed to get user: %v", err)
    return
}

fmt.Printf("Retrieved user: %+v\n", getResp)
```

## Centralized Response Handling

Both APIs use the same response structure for consistency:
- `error`: Boolean indicating if there was an error
- `code`: Status code (HTTP status codes for REST, gRPC codes for gRPC)
- `message`: Human-readable message
- `data`: The actual data payload

This ensures consistent error handling and response formatting across both APIs.
[//]: # (desc: MCP server use http transport)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Gone MCP HTTP Demo

This example demonstrates how to use the Gone MCP component to build HTTP-based client and server applications.

## Project Structure

```
.
├── client/         # Client example code
│   └── main.go     # Client main program
├── config/         # Configuration directory
│   └── default.yaml # Default configuration file
├── go.mod          # Go module definition
└── server/         # Server example code
    ├── functional_add/    # Function definitions
    ├── goner_define/      # Gone component definitions
    ├── import.gone.go     # Gone import file
    └── main.go            # Server main program
```

## Features

This example implements the following features:

1. **Resource Service**: Provides user resource access
   - Supports accessing user information via URI template `users://{id}/profile`
   - Returns user data in JSON format

2. **Code Review Service**: Code review functionality
   - Provides code review assistance
   - Accepts PR number as parameter
   - Returns review suggestions

3. **Calculator Tool**: Basic arithmetic operations
   - Supports addition, subtraction, multiplication, and division
   - Provides parameter validation and error handling
   - Returns calculation results

## Usage

### Server

Run the server program:
   ```bash
   go generator ./...
   go run ./server
   ```

### Client

Run the client program:
   ```bash
   go run ./client
   ```

## Configuration

The configuration file is located at `config/default.yaml`, containing server and client configuration information.

## API Examples

### 1. Accessing User Resources

```http
GET users://123/profile
```

Example response:
```json
{
  "id": 10,
  "name": "Jim"
}
```

### 2. Using the Calculator Tool

Request parameters:
- operation: Operation type (add/subtract/multiply/divide)
- x: First number
- y: Second number

Example:
```json
{
  "operation": "add",
  "x": 10,
  "y": 20
}
```

Response:
```json
{
  "result": 30
}
```

## Notes

1. Ensure the Gone framework and related dependencies are installed
2. Verify the configuration file is correct before running the server
3. Division by zero errors are checked during division operations

## Related Documentation

- [Gone MCP Component Documentation](../../../mcp)
- [Gone Framework Documentation](https://github.com/gone-io/gone)
- [Gone Viper Configuration Documentation](../../../viper)
- [mcp-go](github.com/mark3labs/mcp-go)
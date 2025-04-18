# Gone MCP Quick Start Example

This is a quick start example based on the Gone MCP component, demonstrating how to create a simple MCP tool server and client. Through this example, you can quickly understand and master the basic usage of the Gone MCP component.

## Project Structure

```
.
├── client/         # Client example code
│   └── main.go     # Client main program
├── go.mod          # Go module definition
└── server/         # Server example code
    └── main.go     # Server main program
```

## Feature Description

This example implements a simple greeting service:

1. **Server**:
   - Implements an MCP tool named `hello_world`
   - Accepts a required string parameter `name`
   - Returns formatted greeting

2. **Client**:
   - Communicates with server via stdio
   - Demonstrates complete MCP client usage flow:
     - Initialize connection
     - Get available tools list
     - Call tool and process result

## Usage

### Prerequisites

Ensure Go 1.16 or higher is installed.

### Run Example

1. Clone project and enter example directory:
   ```bash
   cd examples/mcp/quick_start
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run client program:
   ```bash
   go run ./client
   ```

   The client program will automatically start the server and perform:
   - Initialize connection with server
   - Get and display available tools list
   - Call `hello_world` tool and display result

## Example Output

After running client, you will see output similar to:

```
Initialized with server: quick-start 0.0.1

Listing available tools...
- hello_world: Say hello to someone

Calling `hello_world`
Hello, John!
```

## Code Explanation

### Server (server/main.go)

- Uses `gone.Flag` to implement Goner definition
- Defines tool name, description and parameters via `Define()` method
- Implements tool logic in `Handler()` method

### Client (client/main.go)

- Gets MCP client via dependency injection
- Demonstrates complete client initialization and tool calling flow
- Includes error handling and result display best practices

## Extension Suggestions

1. Try adding more parameters to `hello_world` tool
2. Implement new tools like calculator or text processor
3. Try using other communication methods like HTTP or WebSocket

## Related Documents

- [Gone MCP Component Documentation](../../../mcp)
- [Gone Framework Documentation](https://github.com/gone-io/gone)
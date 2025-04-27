[//]: # (desc: MCP server example with standard input/output)

<p>
    English&nbsp ｜&nbsp <a href="README_CN.md">中文</a>
</p>

# Gone MCP Stdio Example

This is a stdio communication example based on the Gone MCP component, demonstrating how to create MCP tool server and client using standard input/output (stdio) method. Through this example, you can learn how to use Gone MCP component in inter-process communication scenarios.

## Project Structure

```
.
├── client/         # Client example code
│   └── main.go     # Client main program
├── config/         # Configuration directory
│   └── default.yaml # Default configuration file
├── go.mod          # Go module definition
└── server/         # Server example code
    ├── functional_add/  # Functional module directory
    ├── goner_define/    # Gone definition files
    ├── import.gone.go   # Gone import file
    └── main.go         # Server main program
```

## Features

This example demonstrates how to implement MCP service using stdio communication:

1. **Server**:
   - Implements stdio-based MCP tool server
   - Supports standard input/output stream communication
   - Contains complete Gone project structure

2. **Client**:
   - Communicates with server via stdio
   - Demonstrates complete MCP client workflow:
     - Initialize connection
     - Get available tool list
     - Call tools and process results

## Usage Scenarios

stdio communication is particularly suitable for:

1. **Inter-process communication**: When communication between different processes on the same machine is needed
2. **Command-line tools**: Developing CLI tools that need to interact with other programs
3. **Plugin systems**: Implementing communication between main program and plugins
4. **Debugging environment**: Convenient for viewing communication data during development and testing

## How to Use

### Prerequisites

Ensure Go 1.16 or later is installed.

### Running the Example

1. Enter example directory:
   ```bash
   cd examples/mcp/stdio
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Generate helper code:
   ```bash
   go generator./...
   ```

4. Run client program:
   ```bash
   go run ./client
   ```

   The client will automatically start the server and communicate via stdio.

## Configuration

### Server Configuration

Server configuration is in `config/default.yaml`, including:

- Basic server information
- Tool configuration
- Log configuration

### Client Configuration

Client is configured in code:

```go
client *client.Client `gone:"*,type=stdio,param=go run ./server"`
```

- `type=stdio`: Specifies stdio communication method
- `param=go run ./server`: Specifies server startup command

## Extension Suggestions

1. Add more tool implementations
2. Implement bidirectional communication
3. Add data compression and encryption
4. Implement more complex inter-process communication scenarios

## Notes

1. stdio communication only works for local inter-process communication
2. Ensure server and client use same protocol version
3. Handle stdio stream buffering and closing properly
4. Recommended for development environment, production may need other communication methods

## Related Documents

- [Gone MCP Component Documentation](../../../mcp)
- [Gone Framework Documentation](https://github.com/gone-io/gone)
- [Gone Viper Configuration Documentation](../../../viper)
- [mcp-go](github.com/mark3labs/mcp-go)
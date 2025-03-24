# The Journey of Gone Framework's Modular Redesign

- [The Journey of Gone Framework's Modular Redesign](#the-journey-of-gone-frameworks-modular-redesign)
  - [I. Redesign Background](#i-redesign-background)
    - [1. The Positioning of Gone Framework](#1-the-positioning-of-gone-framework)
    - [2. The Design of Pluggable Components](#2-the-design-of-pluggable-components)
    - [3. Issues with the Original Architecture](#3-issues-with-the-original-architecture)
  - [II. The Redesign Process](#ii-the-redesign-process)
    - [1. Repository Separation](#1-repository-separation)
    - [2. Single Repository Multiple Modules (Mono Repo) Redesign](#2-single-repository-multiple-modules-mono-repo-redesign)
    - [3. Interface Abstraction and Dependency Decoupling](#3-interface-abstraction-and-dependency-decoupling)
  - [III. Existing Issues and Solutions](#iii-existing-issues-and-solutions)
    - [1. Version Management for Single Repository Multiple Modules in Go Projects](#1-version-management-for-single-repository-multiple-modules-in-go-projects)
      - [Sub-module Version Tags](#sub-module-version-tags)
      - [Version Dependency Management](#version-dependency-management)
    - [2. Unit Testing for Single Repository Multiple Modules in Go Projects](#2-unit-testing-for-single-repository-multiple-modules-in-go-projects)
      - [Module-level Testing](#module-level-testing)
      - [Centralized Test Reporting](#centralized-test-reporting)
    - [3. Component Interdependency Issues](#3-component-interdependency-issues)
      - [Interface Abstraction](#interface-abstraction)
      - [Dependency Injection](#dependency-injection)
      - [Plugin Design](#plugin-design)
  - [IV. Redesign Results](#iv-redesign-results)
  - [V. Future Outlook](#v-future-outlook)
  - [References](#references)


## I. Redesign Background

### 1. The Positioning of Gone Framework

Gone Framework is a lightweight dependency injection framework that helps developers easily build modular, testable applications through a clean API and annotation-based dependency declarations. As a lightweight framework, Gone adheres to the "small is beautiful" design philosophy, providing only the necessary core functionality while avoiding excessive bloat.

### 2. The Design of Pluggable Components

Gone Framework offers a series of pluggable components (called Goners) covering multiple aspects of web development, database access, configuration management, logging, and more:

- gin: Web framework integration
- gorm: ORM framework integration
- viper: Configuration management
- redis: Cache service
- grpc: RPC service
- apollo/nacos: Configuration centers
- zap: Logging component
- And more

This design allows developers to select components based on actual needs, enabling flexible combinations and avoiding forced "all-in-one" dependencies.

### 3. Issues with the Original Architecture

Before the redesign, Gone's core code and all Goner components were managed in the same repository, sharing a single Go module. This architecture had the following problems:

- **Dependency bloat**: Users who needed only a specific function (such as gin integration) still had to import the entire Gone module, including many unnecessary dependencies
- **Difficult version management**: All components shared the same version number, making it impossible to iterate individual components independently
- **High maintenance costs**: As the number of components increased, maintaining a single repository became increasingly difficult
- **Contradiction with lightweight philosophy**: Too many dependencies contradicted Gone Framework's lightweight design philosophy

## II. The Redesign Process

### 1. Repository Separation

First, we separated the Goner component library from the original `github.com/gone-io/gone` repository into an independent new repository: `github.com/gone-io/goner`.

This step included:

- Creating a new Git repository
- Migrating relevant code while preserving commit history
- Adjusting import paths

### 2. Single Repository Multiple Modules (Mono Repo) Redesign

We redesigned `github.com/gone-io/goner` as a Mono Repo, managing each Goner component as an independent sub-module:

```
github.com/gone-io/goner/
├── g/                  # Basic interface definitions
├── gin/                # Gin framework integration
├── gorm/               # GORM framework integration
├── redis/              # Redis client integration
├── viper/              # Viper configuration management
├── apollo/             # Apollo configuration center
├── nacos/              # Nacos configuration center
├── zap/                # Zap logging component
├── grpc/               # gRPC service integration
├── cmux/               # Connection multiplexing
├── schedule/           # Scheduled tasks
├── tracer/             # Distributed tracing
├── urllib/             # HTTP client
├── xorm/               # XORM framework integration
└── ...
```

Each subdirectory is an independent Go module with its own `go.mod` file that can be versioned and published independently.

### 3. Interface Abstraction and Dependency Decoupling

To solve the problem of interdependencies between components, we adopted the "depend on interfaces, not implementations" design principle:

- Created the `github.com/gone-io/goner/g` sub-module for defining all public interfaces
- All components depend on this basic interface module rather than directly depending on other components
- Achieved loose coupling between components through interface abstraction

For example, in `g/tracer.go`, tracing-related interfaces are defined rather than directly depending on the specific implementation of `tracer`:

```go
// /g/tracer.go
package g

// Tracer used for log tracing, assigning a unified traceId to the same call chain for easy log tracking
type Tracer interface {

    //SetTraceId sets a traceId for the calling function. If traceId is an empty string, one will be generated automatically.
    //The calling function can get the traceId through the GetTraceId() method.
    SetTraceId(traceId string, fn func())

    //GetTraceId gets the traceId of the current goroutine
    GetTraceId() string
 
    //Go starts a new goroutine, replacing the native `go func`, and can pass the traceId to the new goroutine
    Go(fn func())
}
```

## III. Existing Issues and Solutions

### 1. Version Management for Single Repository Multiple Modules in Go Projects

Version management for multiple modules in a single repository is a common challenge in the Go module system. We adopted the following strategies:

#### Sub-module Version Tags

We set independent Git tags for each sub-module for version management, ensuring that dependency references can accurately locate the corresponding version of code:

```
<module>/<version>
```

For example:
- `gin/v1.0.0` - Indicates version 1.0.0 of the gin module
- `gorm/v0.2.1` - Indicates version 0.2.1 of the gorm module

This approach allows each sub-module to be versioned independently while maintaining the convenience of management in a single repository.

#### Version Dependency Management

In the `go.mod` files of sub-modules, we use the `replace` directive to handle dependencies between local modules:

```go
// Use local path in development environment
replace github.com/gone-io/goner/g => ../g

// After release, use specific version
require github.com/gone-io/goner/g v0.1.0
```

This ensures that local modules can be used during development, while specific versions of dependencies are used after release.

### 2. Unit Testing for Single Repository Multiple Modules in Go Projects

For unit testing in a single repository with multiple modules, we adopted a strategy of distributed testing with centralized reporting:

#### Module-level Testing

Each sub-module has its own test suite that can be run and validated independently:

```bash
cd gin && go test -v ./...
cd gorm && go test -v ./...
```

#### Centralized Test Reporting

We implemented automated test running and report merging through GitHub Actions workflow (`.github/workflows/go.yml`):

```yaml
- name: Run coverage
  run: find . -name go.mod -not -path "*/example/*" -not -path "*/examples/*" | xargs -n1 dirname | xargs -L1 bash -c 'cd "$0" && pwd && go test -race -coverprofile=coverage.txt -covermode=atomic ./... || exit 255' || echo "Tests failed"
- name: Merge coverage
  run: find . -name "coverage.txt" -exec cat {} \; > total_coverage.txt
- name: Upload coverage reports to Codecov
  uses: codecov/codecov-action@v5
  with:
      token: ${{ secrets.CODECOV_TOKEN }}
      files: ./total_coverage.txt
```

This workflow:
1. Finds all `go.mod` files (excluding example directories)
2. Runs tests in each module directory and generates coverage reports
3. Merges all coverage reports
4. Uploads the merged report to Codecov for analysis

### 3. Component Interdependency Issues

There are interdependencies between Goner components that might lead to circular dependencies or excessive coupling. We solved this through:

#### Interface Abstraction

Abstracting all shared interfaces into the `github.com/gone-io/goner/g` sub-module:

```go
// g/cmux.go
package g

// Listener defines a network listener interface
type Listener interface {
    // ...
}

// g/tracer.go
package g

// Tracer defines a distributed tracing interface
type Tracer interface {
    // ...
}
```

#### Dependency Injection

Using Gone Framework's dependency injection mechanism to dynamically resolve dependencies at runtime:

```go
type MyService struct {
    gone.Flag
    Logger g.Logger `gone:"*"`  // Depending on interface, not specific implementation
}
```

#### Plugin Design

Adopting a plugin design that allows components to be loaded dynamically at runtime:

```go
func main() {
    gone.
        Loads(
            zap.Load,     // Load zap logging component
            gin.Load,     // Load gin component
            gorm.Load,    // Load gorm component
        ).
        Run(func() {
            // Application startup
        })
}
```

## IV. Redesign Results

Through this modular redesign, Gone Framework and its ecosystem have achieved significant improvements:

1. **Lighter dependencies**: Users can import only the components they need, greatly reducing unnecessary dependencies
2. **Flexible version management**: Each component can be versioned independently, facilitating targeted upgrades and maintenance
3. **Better testability**: Loose coupling between components makes unit tests easier to write and maintain
4. **Simplified development process**: Clear interface definitions and dependency relationships make developing new components simpler
5. **Better alignment with Go design philosophy**: Following the "small is beautiful" design principle, each module focuses on solving specific problems

## V. Future Outlook

The modular redesign of Gone Framework is an ongoing process. In the future, we plan to:

1. **Improve documentation**: Provide detailed usage documentation and examples for each component
2. **Expand component ecosystem**: Develop more practical components such as message queues, distributed locks, etc.
3. **Performance optimization**: Optimize core components to improve overall efficiency
4. **Community building**: Encourage community contributions and establish a more active open-source ecosystem

Through these efforts, we hope Gone Framework can become a lightweight, efficient, and easy-to-use dependency injection framework in the Go language ecosystem, providing developers with a better development experience.

## References

- [Gone Framework Official Repository](https://github.com/gone-io/gone)
- [Goner Component Library](https://github.com/gone-io/goner)
- [Go Modules Documentation](https://go.dev/ref/mod)
# util

A collection of Go utility functions for common tasks.

## Features

- **Environment**: Helper functions for looking up environment variables with default values and type conversion (bool,
  int, URL).
- **File**: Secure file operations, path expansion, and structured data loading/saving (JSON, YAML).
- **String**: String template expansion and sensitive string masking.
- **Wait**: Flexible waiting mechanisms for conditions, files, and function returns.

## Installation

```bash
go get github.com/dioad/util
```

## Usage

### Environment Variables

```go
import "github.com/dioad/util"

// Get with default
port := util.LookupEnvWithDefault("PORT", "8080")

// Get as boolean
debug, err := util.LookupEnvBool("DEBUG")

// Get as integer
timeout, err := util.LookupEnvInt("TIMEOUT")
```

### File Operations

```go
import "github.com/dioad/util"

// Safely open a file with path expansion (~, $ENV)
file, err := util.CleanOpen("~/config.json")

// Load a struct from JSON or YAML
type Config struct {
    Name string `json:"name"`
}
config, err := util.LoadStructFromFile[Config]("config.json")
```

### String Masking

```go
import "github.com/dioad/util"

// Create a masked string for sensitive data
apiKey := util.NewMaskedString("secret-api-key")
fmt.Println(apiKey) // Prints masked version
```

### Waiting

```go
import "github.com/dioad/util"

// Wait for a file to exist
err := util.WaitForFile(ctx, time.Second, 10, "ready.txt")
```

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
